package handlers

import (
	"BrandScoutCitatnik/internal/models"
	"BrandScoutCitatnik/internal/repository"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	Repo repository.QuoteRepository
}

func New(repo repository.QuoteRepository) Handler {
	return Handler{Repo: repo}
}

func (h *Handler) CreateQuote(w http.ResponseWriter, r *http.Request) {
	var newQuote models.Quote

	if err := json.NewDecoder(r.Body).Decode(&newQuote); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	result, err := h.Repo.Add(newQuote)

	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) GetQuotes(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")

	var quotes []models.Quote
	var err error
	if author == "" {
		quotes, err = h.Repo.GetAll()
	} else {
		quotes, err = h.Repo.GetByAuthor(author)
	}

	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)
}

func (h *Handler) GetRandomQuotes(w http.ResponseWriter, r *http.Request) {

	quotes, err := h.Repo.GetRandom()

	if err != nil {

		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "No quotes available", http.StatusNotFound)
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quotes)

}

func (h *Handler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	part := strings.Split(r.URL.Path, "/")
	if len(part) != 3 {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(part[2])
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	err = h.Repo.Delete(id)

	if err != nil {

		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "No quotes available", http.StatusNotFound)
		} else {
			http.Error(w, "internal error", http.StatusInternalServerError)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
