package handlers

import (
	"BrandScoutCitatnik/internal/models"
	"BrandScoutCitatnik/internal/repository"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func setupHandler() Handler {
	repo := repository.NewQuoteRepositoryMemory()
	return New(repo)
}

func TestCreateQuote(t *testing.T) {
	t.Run("Успешное создание", func(t *testing.T) {
		h := setupHandler()
		quote := models.Quote{Author: "Тест", Quote: "Привет мир"}
		body, _ := json.Marshal(quote)

		req := httptest.NewRequest("POST", "/quotes", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.CreateQuote(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Ожидался статус 201, получен %d", w.Code)
		}

		var resp models.Quote
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Ошибка декодирования ответа: %v", err)
		}
		if resp.Author != quote.Author || resp.Quote != quote.Quote {
			t.Error("Полученная цитата не соответствует ожидаемой")
		}
	})

	t.Run("Неверный JSON", func(t *testing.T) {
		h := setupHandler()
		body := []byte("{invalid}")

		req := httptest.NewRequest("POST", "/quotes", bytes.NewReader(body))
		w := httptest.NewRecorder()

		h.CreateQuote(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Ожидался статус 400, получен %d", w.Code)
		}
	})
}

func TestGetQuotes(t *testing.T) {
	h := setupHandler()
	h.Repo.Add(models.Quote{Author: "Автор1", Quote: "Цитата1"})
	h.Repo.Add(models.Quote{Author: "Автор1", Quote: "Цитата2"})
	h.Repo.Add(models.Quote{Author: "Автор2", Quote: "Цитата3"})

	t.Run("Получение всех цитат", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/quotes", nil)
		w := httptest.NewRecorder()

		h.GetQuotes(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Ожидался статус 200, получен %d", w.Code)
		}

		var quotes []models.Quote
		if err := json.Unmarshal(w.Body.Bytes(), &quotes); err != nil {
			t.Fatalf("Ошибка декодирования ответа: %v", err)
		}
		if len(quotes) != 3 {
			t.Errorf("Ожидалось 3 цитаты, получено %d", len(quotes))
		}
	})

	t.Run("Фильтрация по автору", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/quotes?author=Автор1", nil)
		w := httptest.NewRecorder()

		h.GetQuotes(w, req)

		var quotes []models.Quote
		if err := json.Unmarshal(w.Body.Bytes(), &quotes); err != nil {
			t.Fatalf("Ошибка декодирования ответа: %v", err)
		}
		if len(quotes) != 2 {
			t.Errorf("Ожидалось 2 цитаты, получено %d", len(quotes))
		}
	})
}

func TestRandomQuotes(t *testing.T) {
	t.Run("Успешное получение случайной цитаты", func(t *testing.T) {
		h := setupHandler()
		h.Repo.Add(models.Quote{Author: "Автор1", Quote: "Цитата1"})
		h.Repo.Add(models.Quote{Author: "Автор2", Quote: "Цитата2"})

		req := httptest.NewRequest("GET", "/quotes/random", nil)
		w := httptest.NewRecorder()

		h.GetRandomQuotes(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Ожидался статус 200, получен %d", w.Code)
		}

		var quote models.Quote
		if err := json.Unmarshal(w.Body.Bytes(), &quote); err != nil {
			t.Fatalf("Ошибка декодирования ответа: %v", err)
		}
		if quote.Author == "" || quote.Quote == "" {
			t.Error("Получена пустая цитата")
		}
	})

	t.Run("Пустое хранилище", func(t *testing.T) {
		h := setupHandler()

		req := httptest.NewRequest("GET", "/quotes/random", nil)
		w := httptest.NewRecorder()

		h.GetRandomQuotes(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Ожидался статус 404, получен %d", w.Code)
		}
	})
}

func TestDeleteQuote(t *testing.T) {
	t.Run("Успешное удаление", func(t *testing.T) {
		h := setupHandler()
		added, _ := h.Repo.Add(models.Quote{Author: "Удаляемый", Quote: "Цитата"})
		idStr := strconv.Itoa(added.ID)

		req := httptest.NewRequest("DELETE", "/quotes/"+idStr, nil)
		w := httptest.NewRecorder()

		h.DeleteQuote(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Ожидался статус 204, получен %d", w.Code)
		}
	})

	t.Run("Неверный ID", func(t *testing.T) {
		h := setupHandler()

		req := httptest.NewRequest("DELETE", "/quotes/abc", nil)
		w := httptest.NewRecorder()

		h.DeleteQuote(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Ожидался статус 400, получен %d", w.Code)
		}
	})

	t.Run("Несуществующий ID", func(t *testing.T) {
		h := setupHandler()

		req := httptest.NewRequest("DELETE", "/quotes/999", nil)
		w := httptest.NewRecorder()

		h.DeleteQuote(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Ожидался статус 404, получен %d", w.Code)
		}
	})

	t.Run("Некорректный путь - меньше частей", func(t *testing.T) {
		h := setupHandler()

		req := httptest.NewRequest("DELETE", "/quotes", nil)
		w := httptest.NewRecorder()

		h.DeleteQuote(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Ожидался статус 400, получен %d", w.Code)
		}
	})

	t.Run("Некорректный путь - больше частей", func(t *testing.T) {
		h := setupHandler()

		req := httptest.NewRequest("DELETE", "/quotes/123/extra", nil) // len == 4
		w := httptest.NewRecorder()

		h.DeleteQuote(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Ожидался статус 400, получен %d", w.Code)
		}
	})
}
