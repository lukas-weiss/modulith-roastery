// Package internal holds the inventory module's internals. Only packages inside inventory/ may import it.
package internal

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"
)

// StockRepository keeps the stock in memory to keep the example small.
type StockRepository struct {
	mu    sync.Mutex
	stock map[string]int
}

func NewStockRepository() *StockRepository {
	return &StockRepository{stock: map[string]int{
		"crema-beans":    12,
		"espresso-beans": 8,
		"espresso-cups":  6,
	}}
}

func (r *StockRepository) Get(product string) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.stock[product]
}

// Reduce takes all quantities out of stock, or none: if one product is short, nothing changes and it
// returns an error. The listener then fails and its event publication stays incomplete.
func (r *StockRepository) Reduce(quantities map[string]int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	var short []string
	for _, product := range slices.Sorted(maps.Keys(quantities)) {
		if quantity := quantities[product]; r.stock[product] < quantity {
			short = append(short, fmt.Sprintf("%s (%d in stock, %d ordered)", product, r.stock[product], quantity))
		}
	}
	if len(short) > 0 {
		return fmt.Errorf("Not enough in stock: %s", strings.Join(short, ", "))
	}
	for product, quantity := range quantities {
		r.stock[product] -= quantity
	}
	return nil
}

func (r *StockRepository) Add(product string, quantity int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stock[product] += quantity
}

func (r *StockRepository) Contains(product string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.stock[product]
	return ok
}

// All returns a copy of the stock. encoding/json writes map keys sorted, like Kotlin's toSortedMap().
func (r *StockRepository) All() map[string]int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return maps.Clone(r.stock)
}
