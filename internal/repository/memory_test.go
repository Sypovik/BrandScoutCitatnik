package repository

import (
	"BrandScoutCitatnik/internal/models"
	"errors"
	"testing"
)

func setup() QuoteRepository {
	repo := NewQuoteRepositoryMemory()
	repo.Add(models.Quote{Author: "Конфуций", Quote: "Жизнь проста."})
	repo.Add(models.Quote{Author: "Аристотель", Quote: "Познать себя - это мудрость."})
	repo.Add(models.Quote{Author: "Конфуций", Quote: "Во всем есть красота."})
	return repo
}

func TestMemoryRepository(t *testing.T) {
	t.Run("Add and GetAll", func(t *testing.T) {
		t.Parallel()

		repo := NewQuoteRepositoryMemory()
		q := models.Quote{Author: "Тест", Quote: "Привет"}

		added, err := repo.Add(q)
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}

		all, _ := repo.GetAll()
		if len(all) != 1 {
			t.Errorf("ожидалась 1 цитата, получено %d", len(all))
		}
		if all[0] != added {
			t.Errorf("ожидалось %+v, получено %+v", added, all[0])
		}
	})

	t.Run("GetByAuthor", func(t *testing.T) {
		t.Parallel()

		repo := setup()

		confuciusQuotes, err := repo.GetByAuthor("Конфуций")
		if err != nil {
			t.Fatalf("ошибка при GetByAuthor: %v", err)
		}
		if len(confuciusQuotes) != 2 {
			t.Errorf("ожидалось 2 цитаты, получено %d", len(confuciusQuotes))
		}
	})

	t.Run("GetRandom", func(t *testing.T) {
		t.Parallel()

		repo := setup()

		quote, err := repo.GetRandom()
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
		if quote == nil {
			t.Fatal("ожидалась не nil цитата")
		}
	})

	t.Run("Delete", func(t *testing.T) {
		t.Parallel()

		repo := setup()
		all, _ := repo.GetAll()
		idToDelete := all[0].ID

		err := repo.Delete(idToDelete)
		if err != nil {
			t.Fatalf("не удалось удалить цитату: %v", err)
		}

		allAfter, _ := repo.GetAll()
		if len(allAfter) != 2 {
			t.Errorf("ожидалось 2 цитаты после удаления, получено %d", len(allAfter))
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
			t.Errorf("неожиданная ошибка: получено %v, ожидалось %v", err.Error(), ErrNotFound.Error())
		}

		_, err = repo.GetRandom()

		if err == nil {
			t.Fatal("Должна быть ошибка, попытка выбрать пустой объект")
		}

		if !errors.Is(err, ErrNotFound) {
			t.Errorf("неожиданная ошибка: получено %v, ожидалось %v", err.Error(), ErrNotFound.Error())
		}

	})
}
