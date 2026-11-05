// Package eventlog is the module that shows what the event registry has stored.
//
// Each published event gets one publication per listener. It is marked COMPLETED when that listener
// has finished, or FAILED if it returned an error. Failed publications stay and can be resubmitted.
// This module only reads the registry and depends on no other module.
package eventlog

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"lweiss.dev/northbean/internal/events"
	"lweiss.dev/northbean/internal/web"
)

// basePackage is the Go counterpart of the Kotlin version's main package dev.lweiss.northbean.
const basePackage = "lweiss.dev/northbean/"

type publication struct {
	ID          string          `json:"id"`
	EventType   string          `json:"eventType"`
	Module      string          `json:"module"`
	Listener    string          `json:"listener"`
	Event       json.RawMessage `json:"event"`
	Status      events.Status   `json:"status"`
	PublishedAt time.Time       `json:"publishedAt"`
	CompletedAt *time.Time      `json:"completedAt"`
	Attempts    int             `json:"attempts"`
}

// Routes adds the HTTP endpoints of the eventlog module.
func Routes(mux *http.ServeMux, registry *events.Registry) {
	mux.HandleFunc("GET /events", func(w http.ResponseWriter, r *http.Request) {
		recent := registry.Recent(60)
		result := make([]publication, len(recent))
		for i, p := range recent {
			result[i] = publication{
				ID:          p.ID,
				EventType:   p.EventType[strings.LastIndex(p.EventType, ".")+1:],
				Module:      moduleOf(p.ListenerID),
				Listener:    listenerOf(p.ListenerID),
				Event:       json.RawMessage(p.SerializedEvent),
				Status:      p.Status,
				PublishedAt: p.PublicationDate,
				CompletedAt: p.CompletionDate,
				Attempts:    p.CompletionAttempts,
			}
		}
		web.WriteJSON(w, result)
	})

	// Hands all failed publications to their listeners again.
	mux.HandleFunc("POST /events/resubmit", func(w http.ResponseWriter, r *http.Request) {
		registry.ResubmitFailed()
	})
}

// moduleOf turns "lweiss.dev/northbean/inventory/internal.OrderPlacedListener.On" into "inventory".
// Without a separator after the prefix it returns the rest unchanged instead of panicking.
func moduleOf(listenerID string) string {
	rest := strings.TrimPrefix(listenerID, basePackage)
	if i := strings.IndexAny(rest, "/."); i >= 0 {
		return rest[:i]
	}
	return rest
}

// listenerOf turns "lweiss.dev/northbean/inventory/internal.OrderPlacedListener.On" into "OrderPlacedListener.On".
func listenerOf(listenerID string) string {
	parts := strings.Split(listenerID[strings.LastIndex(listenerID, "/")+1:], ".")
	return strings.Join(parts[max(0, len(parts)-2):], ".")
}
