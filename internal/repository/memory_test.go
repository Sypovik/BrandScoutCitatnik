package repository

import (
	"BrandScoutCitatnik/internal/models"
	"errors"
	"testing"
)

// setup создает и инициализирует новое in-memory хранилище для каждого теста
func setup() QuoteRepository {
	repo := NewQuoteRepositoryMemory()
	repo.Add(models.Quote{Author: "Confucius", Quote: "Life is simple."})
	repo.Add(models.Quote{Author: "Aristotle", Quote: "Knowing yourself is wisdom."})
	repo.Add(models.Quote{Author: "Confucius", Quote: "Everything has beauty."})
	return repo
}

func TestMemoryRepository(t *testing.T) {
	t.Run("Add and GetAll", func(t *testing.T) {
		t.Parallel()

		repo := NewQuoteRepositoryMemory()
		q := models.Quote{Author: "Test", Quote: "Hello"}

		added, err := repo.Add(q)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		all, _ := repo.GetAll()
		if len(all) != 1 {
			t.Errorf("expected 1 quote, got %d", len(all))
		}
		if all[0] != added {
			t.Errorf("expected %+v, got %+v", added, all[0])
		}
	})

	t.Run("GetByAuthor", func(t *testing.T) {
		t.Parallel()

		repo := setup()

		confuciusQuotes, err := repo.GetByAuthor("Confucius")
		if err != nil {
			t.Fatalf("error on GetByAuthor: %v", err)
		}
		if len(confuciusQuotes) != 2 {
			t.Errorf("expected 2 quotes, got %d", len(confuciusQuotes))
		}
	})

	t.Run("GetRandom", func(t *testing.T) {
		t.Parallel()

		repo := setup()

		quote, err := repo.GetRandom()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if quote == nil {
			t.Fatal("expected non-nil quote")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		t.Parallel()

		repo := setup()
		all, _ := repo.GetAll()
		idToDelete := all[0].ID

		err := repo.Delete(idToDelete)
		if err != nil {
			t.Fatalf("failed to delete quote: %v", err)
		}

		allAfter, _ := repo.GetAll()
		if len(allAfter) != 2 {
			t.Errorf("expected 2 quotes after delete, got %d", len(allAfter))
		}
	})

	t.Run("EmptyQuotation", func(t *testing.T) {
		t.Parallel()

		repo := NewQuoteRepositoryMemory()

		err := repo.Delete(1)
		if err == nil {
			t.Fatal("Должна быть ошибка, попытка удалить пустой объект")
		}

		if !errors.Is(err, ErrNotFound) {
			t.Errorf("unexpected error: got %v, want %v", err.Error(), ErrNotFound.Error())
		}

		_, err = repo.GetRandom()

		if err == nil {
			t.Fatal("Должна быть ошибка, попытка выбрать пустой объект")
		}

		if !errors.Is(err, ErrNotFound) {
			t.Errorf("unexpected error: got %v, want %v", err.Error(), ErrNotFound.Error())
		}

	})
}
