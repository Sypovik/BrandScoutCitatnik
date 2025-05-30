package router

import (
	"BrandScoutCitatnik/internal/handlers"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// func Init(h *handlers.Handler, mw ...Middleware) http.Handler {
// func Init(h *handlers.Handler, mw ...Middleware) http.Handler {
func Init(h handlers.Handler, mw ...Middleware) http.Handler {

	mux := http.NewServeMux()

	// POST и GET /quotes
	mux.HandleFunc("/quotes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreateQuote(w, r)
		case http.MethodGet:
			h.GetQuotes(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// GET /quotes/random
	mux.HandleFunc("/quotes/random", h.GetRandomQuotes)

	// DELETE /quotes/{id}
	mux.HandleFunc("/quotes/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			h.DeleteQuote(w, r)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	})

	var handler http.Handler = mux
	for _, m := range mw {
		handler = m(handler)
	}

	return handler
}
