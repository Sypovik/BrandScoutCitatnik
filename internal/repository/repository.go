package repository

import "BrandScoutCitatnik/internal/models"

type QuoteRepository interface {
	Add(quote models.Quote) (models.Quote, error)

	GetAll() ([]models.Quote, error)
	GetRandom() (*models.Quote, error)

	GetByAuthor(author string) ([]models.Quote, error)

	Delete(id int) error
}
