// Package events is a small stand-in for what Spring Modulith does for the Kotlin version:
// an event publication registry.
//
// Publish stores one publication per listener of the event (the EVENT_PUBLICATION table in Spring
// Modulith), then hands the event to each listener on its own goroutine. A publication is marked
// COMPLETED when its listener has returned without error, or FAILED otherwise. Failed publications
// stay stored and can be resubmitted.
//
// It lives in the top-level internal/ package: every module of the shop may use it, but it is not a
// module itself. Publications are kept in memory, as the Kotlin version keeps them in an in-memory H2.
package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"sync"
	"time"

	"lweiss.dev/northbean/internal/uuid"
)

// Status of a publication, with the same names as in Spring Modulith.
type Status string

const (
	Published   Status = "PUBLISHED"   // stored, not yet handed to the listener
	Processing  Status = "PROCESSING"  // the listener is working on it
	Completed   Status = "COMPLETED"   // the listener has finished
	Failed      Status = "FAILED"      // the listener returned an error, the publication stays stored
	Resubmitted Status = "RESUBMITTED" // handed to the listener again
)

// Publication is one event for one listener: one row of EVENT_PUBLICATION.
type Publication struct {
	ID                 string
	EventType          string // "lweiss.dev/northbean/order.OrderPlaced"
	ListenerID         string // "lweiss.dev/northbean/inventory/internal.OrderPlacedListener.On"
	SerializedEvent    string // the event as JSON
	Status             Status
	PublicationDate    time.Time
	CompletionDate     *time.Time
	CompletionAttempts int
}

// Publisher is what a module needs to publish events, like Spring's ApplicationEventPublisher.
type Publisher interface {
	Publish(event any) error
}

// Listener handles events of type E, like a method annotated with @ApplicationModuleListener.
// It runs on its own goroutine. If it returns an error (or panics), its publication is marked FAILED.
type Listener[E any] interface {
	On(ctx context.Context, event E) error
}

type subscription struct {
	listenerID string
	eventType  string
	handle     func(ctx context.Context, serialized string) error
}

// Registry stores publications and hands events to listeners.
type Registry struct {
	log *slog.Logger

	mu            sync.Mutex
	subscriptions []subscription
	publications  []*Publication

	running sync.WaitGroup
}

func NewRegistry(log *slog.Logger) *Registry {
	return &Registry{log: log}
}

// Listen registers a listener for events of type E. Its ID is built from its package and type name,
// so the event log can tell which module it belongs to.
func Listen[E any](r *Registry, listener Listener[E]) {
	t := reflect.TypeOf(listener)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	s := subscription{
		listenerID: typeName(t) + ".On",
		eventType:  typeName(reflect.TypeFor[E]()),
		// Listeners get the event from its stored JSON, as Spring Modulith does on a resubmit.
		handle: func(ctx context.Context, serialized string) error {
			var event E
			if err := json.Unmarshal([]byte(serialized), &event); err != nil {
				return fmt.Errorf("read event: %w", err)
			}
			return listener.On(ctx, event)
		},
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.subscriptions = append(r.subscriptions, s)
}

// Publish stores one publication per listener of the event and hands it to them asynchronously.
func (r *Registry) Publish(event any) error {
	serialized, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("serialize event: %w", err)
	}
	eventType := typeName(reflect.TypeOf(event))
	now := time.Now()

	r.mu.Lock()
	var deliveries []delivery
	for _, s := range r.subscriptions {
		if s.eventType != eventType {
			continue
		}
		p := &Publication{
			ID:              uuid.New(),
			EventType:       eventType,
			ListenerID:      s.listenerID,
			SerializedEvent: string(serialized),
			Status:          Published,
			PublicationDate: now,
		}
		r.publications = append(r.publications, p)
		deliveries = append(deliveries, delivery{p, s})
	}
	r.mu.Unlock()

	for _, d := range deliveries {
		r.dispatch(d)
	}
	return nil
}

// ResubmitFailed hands all failed publications to their listeners again.
func (r *Registry) ResubmitFailed() {
	r.mu.Lock()
	var deliveries []delivery
	for _, p := range r.publications {
		if p.Status != Failed {
			continue
		}
		for _, s := range r.subscriptions {
			if s.listenerID == p.ListenerID && s.eventType == p.EventType {
				p.Status = Resubmitted
				deliveries = append(deliveries, delivery{p, s})
			}
		}
	}
	r.mu.Unlock()

	for _, d := range deliveries {
		r.dispatch(d)
	}
}

// Recent returns copies of the newest publications, newest first.
func (r *Registry) Recent(limit int) []Publication {
	r.mu.Lock()
	defer r.mu.Unlock()
	recent := make([]Publication, 0, len(r.publications))
	for _, p := range slices.Backward(r.publications) {
		recent = append(recent, *p)
	}
	slices.SortStableFunc(recent, func(a, b Publication) int { return b.PublicationDate.Compare(a.PublicationDate) })
	return recent[:min(limit, len(recent))]
}

// Wait blocks until all listeners started so far have finished. Tests use it.
func (r *Registry) Wait() {
	r.running.Wait()
}

type delivery struct {
	publication  *Publication
	subscription subscription
}

func (r *Registry) dispatch(d delivery) {
	r.running.Go(func() {
		r.mu.Lock()
		d.publication.Status = Processing
		d.publication.CompletionAttempts++
		r.mu.Unlock()

		err := call(d.subscription, d.publication.SerializedEvent)

		r.mu.Lock()
		defer r.mu.Unlock()
		if err != nil {
			d.publication.Status = Failed
			r.log.Error("Listener failed, the publication stays stored", "listener", d.subscription.listenerID, "error", err)
			return
		}
		now := time.Now()
		d.publication.Status = Completed
		d.publication.CompletionDate = &now
	})
}

// call runs the listener and turns a panic into an error, so a broken listener cannot stop the shop.
func call(s subscription, serialized string) (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("panic: %v", p)
		}
	}()
	return s.handle(context.Background(), serialized)
}

func typeName(t reflect.Type) string {
	return t.PkgPath() + "." + t.Name()
}
