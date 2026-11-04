package inventory_test

import (
	"log/slog"
	"testing"
	"time"

	"example.com/shop/internal/events"
	"example.com/shop/internal/uuid"
	"example.com/shop/inventory"
	"example.com/shop/order"
)

// Starts only the inventory module and feeds it an event, as if the order module had sent it.
func TestStockIsReducedWhenAnOrderIsPlaced(t *testing.T) {
	registry := events.NewRegistry(slog.Default())
	stock := inventory.New(registry)

	if err := registry.Publish(order.OrderPlaced{OrderID: uuid.New(), Lines: []order.OrderLine{{Product: "espresso-cups", Quantity: 2}}}); err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().Add(5 * time.Second)
	for stock.StockOf("espresso-cups") != 4 {
		if time.Now().After(deadline) {
			t.Fatalf("stock is %d, want 4", stock.StockOf("espresso-cups"))
		}
		time.Sleep(10 * time.Millisecond)
	}
}
