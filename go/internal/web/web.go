// Package web holds the small JSON helpers the HTTP handlers of all modules share.
package web

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"lweiss.dev/northbean/internal/check"
)

// ReadJSON decodes the request body into target. On failure it answers 400 and returns false.
func ReadJSON(w http.ResponseWriter, r *http.Request, target any) bool {
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20)).Decode(target); err != nil {
		BadRequest(w, "Invalid JSON: "+err.Error())
		return false
	}
	return true
}

// WriteJSON answers 200 with value as JSON.
func WriteJSON(w http.ResponseWriter, value any) {
	writeJSON(w, http.StatusOK, value)
}

// BadRequest answers 400 with {"error": message}.
func BadRequest(w http.ResponseWriter, message string) {
	writeJSON(w, http.StatusBadRequest, map[string]string{"error": message})
}

// WriteError answers 400 for a check.ArgumentError and 500 for anything else.
func WriteError(w http.ResponseWriter, err error) {
	if argument, ok := errors.AsType[*check.ArgumentError](err); ok {
		BadRequest(w, argument.Message)
		return
	}
	slog.Error("Request failed", "error", err)
	writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
