package order

import (
	"cmp"
	"slices"

	"example.com/shop/internal/check"
	"example.com/shop/internal/events"
	"example.com/shop/internal/uuid"
	"example.com/shop/order/internal"
)

// Service is the public API of the order module.
type Service struct {
	repository *internal.OrderRepository
	events     events.Publisher
}

// New creates the order module. It needs nothing from other modules, only a way to publish events.
func New(publisher events.Publisher) *Service {
	return &Service{repository: internal.NewOrderRepository(), events: publisher}
}

// Place stores the order and publishes OrderPlaced. It returns the ID of the new order.
func (s *Service) Place(lines []OrderLine) (string, error) {
	if err := cmp.Or(
		check.Require(len(lines) > 0, "An order needs at least one line"),
		check.Require(!slices.ContainsFunc(lines, func(l OrderLine) bool { return l.Quantity < 1 || l.Quantity > 100 }), "Quantity must be between 1 and 100"),
		check.Require(distinctProducts(lines), "Each product may appear only once"),
	); err != nil {
		return "", err
	}

	order := s.repository.Save(internal.Order{ID: uuid.New(), Lines: toInternal(lines)})

	// No call to inventory or notification here – just tell the world what happened.
	// Listeners run on their own goroutines once the order is stored.
	if err := s.events.Publish(OrderPlaced{OrderID: order.ID, Lines: lines}); err != nil {
		return "", err
	}
	return order.ID, nil
}

func distinctProducts(lines []OrderLine) bool {
	products := map[string]bool{}
	for _, l := range lines {
		products[l.Product] = true
	}
	return len(products) == len(lines)
}

func toInternal(lines []OrderLine) []internal.Line {
	result := make([]internal.Line, len(lines))
	for i, l := range lines {
		result[i] = internal.Line{Product: l.Product, Quantity: l.Quantity}
	}
	return result
}
