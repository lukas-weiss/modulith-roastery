package order_test

import (
	"context"
	"log/slog"
	"reflect"
	"testing"

	"lweiss.dev/northbean/internal/events"
	"lweiss.dev/northbean/order"
)

type recorder struct{ events chan order.OrderPlaced }

func (r *recorder) On(_ context.Context, event order.OrderPlaced) error {
	r.events <- event
	return nil
}

// Starts only the order module – inventory and notification are not loaded.
func TestPlacingAnOrderPublishesOrderPlaced(t *testing.T) {
	registry := events.NewRegistry(slog.Default())
	received := &recorder{events: make(chan order.OrderPlaced, 1)}
	events.Listen(registry, received)
	orders := order.New(registry)

	lines := []order.OrderLine{{Product: "crema-beans", Quantity: 2}, {Product: "espresso-cups", Quantity: 1}}
	id, err := orders.Place(lines)
	if err != nil {
		t.Fatal(err)
	}

	event := <-received.events
	if event.OrderID != id || !reflect.DeepEqual(event.Lines, lines) {
		t.Fatalf("got %+v", event)
	}
}
