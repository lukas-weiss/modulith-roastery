package events

import (
	"context"
	"errors"
	"log/slog"
	"testing"
)

type greeting struct{ Name string }

type flaky struct{ failures int }

func (f *flaky) On(_ context.Context, g greeting) error {
	if f.failures > 0 {
		f.failures--
		return errors.New("not today")
	}
	return nil
}

type broken struct{}

func (broken) On(context.Context, greeting) error { panic("boom") }

func TestFailedPublicationStaysUntilResubmitted(t *testing.T) {
	registry := NewRegistry(slog.Default())
	Listen(registry, &flaky{failures: 1})

	if err := registry.Publish(greeting{"Ada"}); err != nil {
		t.Fatal(err)
	}
	registry.Wait()

	p := registry.Recent(10)[0]
	if p.Status != Failed || p.CompletionAttempts != 1 || p.CompletionDate != nil {
		t.Fatalf("after publish: %+v", p)
	}
	if p.ListenerID != "example.com/shop/internal/events.flaky.On" || p.EventType != "example.com/shop/internal/events.greeting" {
		t.Fatalf("ids: %+v", p)
	}
	if p.SerializedEvent != `{"Name":"Ada"}` {
		t.Fatalf("serialized: %s", p.SerializedEvent)
	}

	registry.ResubmitFailed()
	registry.Wait()

	p = registry.Recent(10)[0]
	if p.Status != Completed || p.CompletionAttempts != 2 || p.CompletionDate == nil {
		t.Fatalf("after resubmit: %+v", p)
	}
}

func TestPanickingListenerFails(t *testing.T) {
	registry := NewRegistry(slog.Default())
	Listen(registry, broken{})

	_ = registry.Publish(greeting{"Ada"})
	registry.Wait()

	if p := registry.Recent(10)[0]; p.Status != Failed {
		t.Fatalf("status = %s", p.Status)
	}
}

func TestEventsWithoutListenersAreNotStored(t *testing.T) {
	registry := NewRegistry(slog.Default())
	_ = registry.Publish(greeting{"Ada"})
	if n := len(registry.Recent(10)); n != 0 {
		t.Fatalf("stored %d publications", n)
	}
}
