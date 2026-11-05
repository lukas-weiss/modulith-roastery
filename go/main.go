// The online shop of a small coffee roastery, split into modules that talk through events.
//
// Every direct sub-package of lweiss.dev/northbean is a module (internal/ holds the shared plumbing):
//
//   - order        – accepts orders and publishes order.OrderPlaced
//   - inventory    – reacts to OrderPlaced and reduces the stock
//   - notification – reacts to OrderPlaced and "sends" a confirmation
//   - eventlog     – shows the stored event publications in the web UI
//
// order knows nothing about the other modules. They only talk through events.
package main

import (
	"cmp"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"

	"lweiss.dev/northbean/eventlog"
	"lweiss.dev/northbean/internal/events"
	"lweiss.dev/northbean/inventory"
	"lweiss.dev/northbean/notification"
	"lweiss.dev/northbean/order"
)

//go:embed static
var static embed.FS

// newShop creates the modules and wires them up. Spring does this for the Kotlin version.
func newShop() (*http.ServeMux, *events.Registry) {
	registry := events.NewRegistry(slog.Default())

	orders := order.New(registry)
	stock := inventory.New(registry)
	notification.New(registry)

	mux := http.NewServeMux()
	order.Routes(mux, orders)
	inventory.Routes(mux, stock)
	eventlog.Routes(mux, registry)
	files, _ := fs.Sub(static, "static")
	mux.Handle("GET /", http.FileServerFS(files))
	return mux, registry
}

func main() {
	mux, _ := newShop()
	addr := ":" + cmp.Or(os.Getenv("PORT"), "8080")
	slog.Info("Shop started on http://localhost" + addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("Server stopped", "error", err)
		os.Exit(1)
	}
}
