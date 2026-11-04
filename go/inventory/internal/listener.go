package internal

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"example.com/shop/order"
)

// OrderPlacedListener reduces the stock when an order is placed. It may use order.OrderPlaced,
// because that is part of the order module's public API.
type OrderPlacedListener struct {
	Stock *StockRepository
}

// On runs on its own goroutine, after the order has been stored (see events.Listen).
func (l *OrderPlacedListener) On(_ context.Context, event order.OrderPlaced) error {
	quantities := map[string]int{}
	for _, line := range event.Lines {
		quantities[line.Product] = line.Quantity
	}
	if err := l.Stock.Reduce(quantities); err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("Inventory: reduced stock for order %s: %s", event.OrderID, describe(event.Lines)))
	return nil
}

// describe writes lines as "3 x espresso-beans, 2 x espresso-cups".
func describe(lines []order.OrderLine) string {
	parts := make([]string, len(lines))
	for i, line := range lines {
		parts[i] = fmt.Sprintf("%d x %s", line.Quantity, line.Product)
	}
	return strings.Join(parts, ", ")
}
