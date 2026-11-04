package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// The whole shop over HTTP: the same flow as the README describes for curl and the web UI.
func TestShortStockFailsInventoryUntilResubmitted(t *testing.T) {
	mux, registry := newShop()
	server := httptest.NewServer(mux)
	defer server.Close()

	// 9 espresso cups, but only 6 are in stock: inventory fails for the whole order.
	var placed struct{ OrderID string }
	call(t, server, "POST", "/orders", `{"lines":[{"product":"espresso-beans","quantity":3},{"product":"espresso-cups","quantity":9}]}`, 200, &placed)
	registry.Wait()

	stock := map[string]int{}
	call(t, server, "GET", "/inventory", "", 200, &stock)
	if stock["espresso-beans"] != 8 || stock["espresso-cups"] != 6 {
		t.Fatalf("stock changed although inventory failed: %v", stock)
	}
	if got := statuses(t, server, placed.OrderID); got["inventory"] != "FAILED 1" || got["notification"] != "COMPLETED 1" {
		t.Fatalf("statuses = %v", got)
	}

	call(t, server, "POST", "/inventory/restock", `{"product":"espresso-cups","quantity":10}`, 200, &stock)
	call(t, server, "POST", "/events/resubmit", "", 200, nil)
	registry.Wait()

	call(t, server, "GET", "/inventory", "", 200, &stock)
	if stock["espresso-beans"] != 5 || stock["espresso-cups"] != 7 {
		t.Fatalf("stock after resubmit = %v", stock)
	}
	if got := statuses(t, server, placed.OrderID); got["inventory"] != "COMPLETED 2" || got["notification"] != "COMPLETED 1" {
		t.Fatalf("statuses after resubmit = %v", got)
	}
}

func TestBadOrdersAreRejected(t *testing.T) {
	mux, _ := newShop()
	server := httptest.NewServer(mux)
	defer server.Close()

	for _, body := range []string{
		`{"lines":[]}`,
		`{"lines":[{"product":"crema-beans","quantity":0}]}`,
		`{"lines":[{"product":"crema-beans","quantity":101}]}`,
		`{"lines":[{"product":"crema-beans","quantity":1},{"product":"crema-beans","quantity":2}]}`,
		`{"lines":[{"product":"crema-beans"}]}`,
		`not json`,
	} {
		var answer struct{ Error string }
		call(t, server, "POST", "/orders", body, 400, &answer)
		if answer.Error == "" {
			t.Errorf("%s: no error message", body)
		}
	}
	call(t, server, "POST", "/inventory/restock", `{"product":"tea","quantity":10}`, 400, nil)
	call(t, server, "POST", "/inventory/restock", `{"product":"crema-beans","quantity":101}`, 400, nil)
}

// statuses maps module to "STATUS attempts" for the publications of one order.
func statuses(t *testing.T, server *httptest.Server, orderID string) map[string]string {
	var publications []struct {
		Module   string
		Listener string
		Status   string
		Attempts int
		Event    struct{ OrderID string }
	}
	call(t, server, "GET", "/events", "", 200, &publications)
	result := map[string]string{}
	for _, p := range publications {
		if p.Event.OrderID == orderID {
			result[p.Module] = fmt.Sprintf("%s %d", p.Status, p.Attempts)
		}
	}
	return result
}

func call(t *testing.T, server *httptest.Server, method, path, body string, status int, target any) {
	t.Helper()
	request, _ := http.NewRequest(method, server.URL+path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != status {
		t.Fatalf("%s %s %s: status %d, want %d", method, path, body, response.StatusCode, status)
	}
	if target != nil {
		if err := json.NewDecoder(response.Body).Decode(target); err != nil {
			t.Fatalf("%s %s: %v", method, path, err)
		}
	}
}
