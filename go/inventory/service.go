package inventory

import (
	"cmp"

	"lweiss.dev/northbean/internal/check"
	"lweiss.dev/northbean/internal/events"
	"lweiss.dev/northbean/inventory/internal"
)

// Service is the public API of the inventory module.
type Service struct {
	stock *internal.StockRepository
}

// New creates the inventory module and lets it listen to OrderPlaced.
func New(registry *events.Registry) *Service {
	stock := internal.NewStockRepository()
	events.Listen(registry, &internal.OrderPlacedListener{Stock: stock})
	return &Service{stock: stock}
}

func (s *Service) StockOf(product string) int { return s.stock.Get(product) }

func (s *Service) AllStock() map[string]int { return s.stock.All() }

func (s *Service) Restock(product string, quantity int) error {
	if err := cmp.Or(
		check.Require(s.stock.Contains(product), "Unknown product: "+product),
		check.Require(quantity >= 1 && quantity <= 100, "Quantity must be between 1 and 100"),
	); err != nil {
		return err
	}
	s.stock.Add(product, quantity)
	return nil
}
