package inventory

import (
	"net/http"

	"lweiss.dev/northbean/internal/web"
)

type restockRequest struct {
	Product  *string `json:"product"`
	Quantity *int    `json:"quantity"`
}

// Routes adds the HTTP endpoints of the inventory module.
func Routes(mux *http.ServeMux, inventory *Service) {
	mux.HandleFunc("GET /inventory", func(w http.ResponseWriter, r *http.Request) {
		web.WriteJSON(w, inventory.AllStock())
	})

	mux.HandleFunc("POST /inventory/restock", func(w http.ResponseWriter, r *http.Request) {
		var request restockRequest
		if !web.ReadJSON(w, r, &request) {
			return
		}
		if request.Product == nil || request.Quantity == nil {
			web.BadRequest(w, "A restock needs a product and a quantity")
			return
		}
		if err := inventory.Restock(*request.Product, *request.Quantity); err != nil {
			web.WriteError(w, err)
			return
		}
		web.WriteJSON(w, inventory.AllStock())
	})
}
