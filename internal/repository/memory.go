package repository

import (
	"BrandScoutCitatnik/internal/models"
	"errors"
	"math/rand"
	"sync"
)

type MemoryQuoteRepository struct {
	sync.RWMutex
	data   map[int]models.Quote
	nextID int
}

func NewMemoryQuoteRepository() *MemoryQuoteRepository {
	return &MemoryQuoteRepository{
		data:   make(map[int]models.Quote),
		nextID: 1,
	}
}

func (m *MemoryQuoteRepository) Add(quote models.Quote) (models.Quote, error) {
	m.Lock()
	defer m.Unlock()

	quote.ID = m.nextID
	m.data[m.nextID] = quote
	m.nextID++

	return quote, nil
}

func (m *MemoryQuoteRepository) GetAll() ([]models.Quote, error) {
	m.RLock()
	defer m.RUnlock()
	var arrQuotes []models.Quote
	// arrQuotes := make([]models.Quote, 0, len(m.data))

	for _, i := range m.data {
		arrQuotes = append(arrQuotes, i)
	}
	return arrQuotes, nil
}

// !!!Не работает потому, что при удалении, в данном подходе, ID может не быть
// !!!Нужно смотреть какие ID есть и среди них выводить случайный
// !!!Добавить зерно сида

func (m *MemoryQuoteRepository) GetRandom() (*models.Quote, error) {
	m.RLock()
	defer m.RUnlock()

	if len(m.data) == 0 {
		return nil, errors.New("нет цитат")
	}

	keys := make([]int, 0, len(m.data))

	for _, key := range m.data {
		keys = append(keys, key.ID)
	}

	randomID := keys[rand.Intn(len(keys))]
	quote := (m.data[randomID])

	return &quote, nil
}

func (m *MemoryQuoteRepository) GetByAuthor(author string) ([]models.Quote, error) {
	m.RLock()
	defer m.RUnlock()
	var authors []models.Quote
	for _, i := range m.data {
		if i.Author == author {
			authors = append(authors, i)
		}
	}
	return authors, nil
}

func (m *MemoryQuoteRepository) Delete(id int) error {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.data[id]; !ok {
		return errors.New("quote not found")
	}
	delete(m.data, id)
	return nil
}
