package service

import (
	"go-mux-trivia/internal/database"
	"go-mux-trivia/internal/models"
)

type TriviaBundleService interface {
	GetTriviaBundle(id int) (*models.TriviaBundle, error)
	GetTriviaBundlesByCategory(category string) ([]*models.TriviaBundle, error)
	GetTriviaBundles() ([]*models.TriviaBundle, error)
	CreateTriviaBundle(triviaBundle *models.TriviaBundle) (int, error)
	UpdateTriviaBundle(triviaBundle *models.TriviaBundle) error
	DeleteTriviaBundle(id int) error
}

type triviaBundleService struct {
	db database.Service
}

func NewTriviaBundleService(db database.Service) TriviaBundleService {
	return &triviaBundleService{
		db: db,
	}
}

func (s *triviaBundleService) GetTriviaBundle(id int) (*models.TriviaBundle, error) {
	bundle, err := s.db.GetTriviaBundle(id)
	return bundle, err
}

func (s *triviaBundleService) CreateTriviaBundle(tb *models.TriviaBundle) (int, error) {
	return s.db.CreateTriviaBundle(tb)
}

func (s *triviaBundleService) UpdateTriviaBundle(tb *models.TriviaBundle) error {
	return s.db.UpdateTriviaBundle(tb)
}

func (s *triviaBundleService) DeleteTriviaBundle(id int) error {
	return s.db.DeleteTriviaBundle(id)
}

func (s *triviaBundleService) GetTriviaBundles() ([]*models.TriviaBundle, error) {
	return s.db.GetTriviaBundles()
}

func (s *triviaBundleService) GetTriviaBundlesByCategory(category string) ([]*models.TriviaBundle, error) {
	return s.db.GetTriviaBundlesByCategory(category)
}
