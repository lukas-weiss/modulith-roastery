// Package notification is the module that "sends" a confirmation for each order.
package notification

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"example.com/shop/internal/events"
	"example.com/shop/order"
)

// NotificationListener reacts to OrderPlaced. (The name repeats the package to match the Kotlin version.)
type NotificationListener struct{}

// New creates the notification module and lets it listen to OrderPlaced.
func New(registry *events.Registry) {
	events.Listen(registry, &NotificationListener{})
}

func (l *NotificationListener) On(_ context.Context, event order.OrderPlaced) error {
	// A real app would send an email here.
	parts := make([]string, len(event.Lines))
	for i, line := range event.Lines {
		parts[i] = fmt.Sprintf("%d x %s", line.Quantity, line.Product)
	}
	slog.Info(fmt.Sprintf("Notification: order %s confirmed (%s)", event.OrderID, strings.Join(parts, ", ")))
	return nil
}
