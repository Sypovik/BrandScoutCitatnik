package repository

import (
	"BrandScoutCitatnik/internal/models"
	"testing"
)

func TestSaveAndFindAll(t *testing.T) {
	repo := NewMemoryQuoteRepository()

	q := models.Quote{
		Author: "Confucius",
		Quote:  "Life is simple, but we insist on making it complicated.",
	}

	_, err := repo.Add(q)
	if err != nil {
		t.Fatalf("unexpected error on Save: %v", err)
	}

	all, err := repo.GetAll()
	if err != nil {
		t.Fatalf("unexpected error on FindAll: %v", err)
	}

	if len(all) != 1 {
		t.Errorf("expected 1 quote, got %d", len(all))
	}

	if all[0].ID != 1 {
		t.Errorf("expected ID 1, got %d", all[0].ID)
	}
}

func TestFindByAuthor(t *testing.T) {
	repo := NewMemoryQuoteRepository()

	quotes := []models.Quote{
		{Author: "Confucius", Quote: "One"},
		{Author: "Aristotle", Quote: "Two"},
		{Author: "Confucius", Quote: "Three"},
	}

	for _, q := range quotes {
		_, _ = repo.Add(q)
	}

	found, err := repo.GetByAuthor("Confucius")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(found) != 2 {
		t.Errorf("expected 2 quotes, got %d", len(found))
	}
}

func TestFindRandom(t *testing.T) {
	repo := NewMemoryQuoteRepository()

	_, err := repo.GetRandom()
	if err == nil {
		t.Error("expected error on empty repository")
	}

	_, _ = repo.Add(models.Quote{Author: "A", Quote: "1"})
	_, _ = repo.Add(models.Quote{Author: "B", Quote: "2"})
	_, _ = repo.Add(models.Quote{Author: "C", Quote: "3"})
	_, _ = repo.Add(models.Quote{Author: "D", Quote: "4"})
	_, _ = repo.Add(models.Quote{Author: "D", Quote: "5"})

	q, err := repo.GetRandom()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q == nil {
		t.Error("expected non-nil quote")
	}
}

func TestDelete(t *testing.T) {
	repo := NewMemoryQuoteRepository()

	q := models.Quote{Author: "A", Quote: "Q"}
	_, _ = repo.Add(q)

	err := repo.Delete(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = repo.Delete(1)
	if err == nil {
		t.Error("expected error when deleting non-existent ID")
	}
}
