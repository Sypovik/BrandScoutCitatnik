package router

import (
	"BrandScoutCitatnik/internal/handlers"
	"BrandScoutCitatnik/internal/models"
	"BrandScoutCitatnik/internal/repository"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupRouter() http.Handler {
	repo := repository.NewQuoteRepositoryMemory()
	handler := handlers.New(repo)
	return Init(handler)
}

func TestRouterQuotesEndpoints(t *testing.T) {
	t.Run("POST /quotes и GET /quotes", func(t *testing.T) {
		router := setupRouter()

		// Добавим цитату
		body := `{"author":"Marcus Aurelius", "quote":"You have power over your mind."}`
		req := httptest.NewRequest(http.MethodPost, "/quotes", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Ожидался статус 201, получен %d", w.Code)
		}

		// Получим все цитаты
		req = httptest.NewRequest(http.MethodGet, "/quotes", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Ожидался статус 200, получен %d", w.Code)
		}

		var quotes []models.Quote
		if err := json.NewDecoder(w.Body).Decode(&quotes); err != nil {
			t.Fatalf("Ошибка декодирования: %v", err)
		}
		if len(quotes) != 1 || quotes[0].Author != "Marcus Aurelius" {
			t.Error("Неверное содержимое цитаты")
		}
	})

	t.Run("GET /quotes?author=Seneca", func(t *testing.T) {
		router := setupRouter()

		// Добавим несколько цитат
		quotes := []models.Quote{
			{Author: "Seneca", Quote: "Time heals all wounds."},
			{Author: "Seneca", Quote: "Luck is preparation."},
			{Author: "Other", Quote: "Something else."},
		}
		for _, q := range quotes {
			body, _ := json.Marshal(q)
			req := httptest.NewRequest(http.MethodPost, "/quotes", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		req := httptest.NewRequest(http.MethodGet, "/quotes?author=Seneca", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var filtered []models.Quote
		json.NewDecoder(w.Body).Decode(&filtered)
		if len(filtered) != 2 {
			t.Errorf("Ожидалось 2 цитаты от Seneca, получено %d", len(filtered))
		}
	})
}

func TestRouterRandomQuote(t *testing.T) {
	t.Run("GET /quotes/random", func(t *testing.T) {
		router := setupRouter()

		// Добавим цитату
		req := httptest.NewRequest(http.MethodPost, "/quotes", strings.NewReader(`{"author":"Epictetus", "quote":"Endure and abstain."}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Получим случайную цитату
		req = httptest.NewRequest(http.MethodGet, "/quotes/random", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Ожидался статус 200, получен %d", w.Code)
		}
		var q models.Quote
		json.NewDecoder(w.Body).Decode(&q)
		if q.Author != "Epictetus" {
			t.Errorf("Ожидался Epictetus, получен %s", q.Author)
		}
	})
}

func TestRouterDeleteQuote(t *testing.T) {
	t.Run("Успешное удаление и повторная попытка", func(t *testing.T) {
		router := setupRouter()

		// Добавим цитату
		req := httptest.NewRequest(http.MethodPost, "/quotes", strings.NewReader(`{"author":"Seneca", "quote":"Delete me"}`))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Удалим
		req = httptest.NewRequest(http.MethodDelete, "/quotes/1", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Ожидался статус 204, получен %d", w.Code)
		}

		// Повторное удаление
		req = httptest.NewRequest(http.MethodDelete, "/quotes/1", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Ожидался статус 404, получен %d", w.Code)
		}
	})

	t.Run("Неверный ID", func(t *testing.T) {
		router := setupRouter()

		req := httptest.NewRequest(http.MethodDelete, "/quotes/abc", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Ожидался статус 400, получен %d", w.Code)
		}
	})
}

func TestRouterInvalidMethods(t *testing.T) {
	t.Run("PUT /quotes возвращает 405", func(t *testing.T) {
		router := setupRouter()

		req := httptest.NewRequest(http.MethodPut, "/quotes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Ожидался 405, получен %d", w.Code)
		}
	})

	t.Run("GET /quotes/1 (не DELETE) возвращает 404", func(t *testing.T) {
		router := setupRouter()

		req := httptest.NewRequest(http.MethodGet, "/quotes/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Ожидался статус 404, получен %d", w.Code)
		}
	})
}

func TestRouterMiddleware(t *testing.T) {
	t.Run("Middleware должен быть вызван", func(t *testing.T) {
		var called bool

		repo := repository.NewQuoteRepositoryMemory()
		handler := handlers.New(repo)
		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				next.ServeHTTP(w, r)
			})
		}

		router := Init(handler, middleware)

		req := httptest.NewRequest(http.MethodGet, "/quotes", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if !called {
			t.Error("Ожидалось, что middleware будет вызвано")
		}
	})
}
