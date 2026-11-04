package order

import (
	"net/http"

	"example.com/shop/internal/web"
)

type placeOrderRequest struct {
	Lines []struct {
		Product  *string `json:"product"`
		Quantity *int    `json:"quantity"`
	} `json:"lines"`
}

type placeOrderResponse struct {
	OrderID string `json:"orderId"`
}

// Routes adds the HTTP endpoints of the order module.
func Routes(mux *http.ServeMux, orders *Service) {
	mux.HandleFunc("POST /orders", func(w http.ResponseWriter, r *http.Request) {
		var request placeOrderRequest
		if !web.ReadJSON(w, r, &request) {
			return
		}
		lines := make([]OrderLine, 0, len(request.Lines))
		for _, l := range request.Lines {
			// Like Kotlin's non-null properties: a line without product or quantity is a bad request.
			if l.Product == nil || l.Quantity == nil {
				web.BadRequest(w, "Each line needs a product and a quantity")
				return
			}
			lines = append(lines, OrderLine{Product: *l.Product, Quantity: *l.Quantity})
		}

		id, err := orders.Place(lines)
		if err != nil {
			web.WriteError(w, err)
			return
		}
		web.WriteJSON(w, placeOrderResponse{OrderID: id})
	})
}
