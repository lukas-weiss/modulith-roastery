package order

// OrderLine is one line of an order: a product and how many of it. Part of the order module's public API.
type OrderLine struct {
	Product  string `json:"product"`
	Quantity int    `json:"quantity"`
}

// OrderPlaced is the domain event published by the order module.
// It lives in the module's base package, so it is part of the module's public API
// and other modules may use it.
type OrderPlaced struct {
	OrderID string      `json:"orderId"`
	Lines   []OrderLine `json:"lines"`
}
