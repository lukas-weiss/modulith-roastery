// Package internal holds the order module's internals. Go itself enforces what Spring Modulith checks
// in the Kotlin version: only packages inside order/ may import order/internal.
package internal

import (
	"maps"
	"slices"
	"sync"
)

// Order is internal: other modules must not use it.
type Order struct {
	ID    string
	Lines []Line
}

// Line is a product and its quantity. It does not use order.OrderLine because order imports this
// package, and Go does not allow import cycles.
type Line struct {
	Product  string
	Quantity int
}

// OrderRepository keeps orders in memory to keep the example small.
type OrderRepository struct {
	mu     sync.Mutex
	orders map[string]Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{orders: map[string]Order{}}
}

func (r *OrderRepository) Save(order Order) Order {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	return order
}

func (r *OrderRepository) FindAll() []Order {
	r.mu.Lock()
	defer r.mu.Unlock()
	return slices.Collect(maps.Values(r.orders))
}
