package service

import (
	"go-mux-trivia/internal/database"
	"go-mux-trivia/internal/models"
	"go-mux-trivia/internal/repository"
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
	r repository.TriviaRepository
}

func NewTriviaBundleService() *triviaBundleService {
	return &triviaBundleService{
		r: repository.NewTriviaRepository(database.New()),
	}
}

func (s *triviaBundleService) GetTriviaBundle(id int) (*models.TriviaBundle, error) {
	bundle, err := s.r.GetTriviaBundle(id)
	return bundle, err
}

func (s *triviaBundleService) CreateTriviaBundle(tb *models.TriviaBundle) (int, error) {
	return s.r.CreateTriviaBundle(tb)
}

func (s *triviaBundleService) UpdateTriviaBundle(tb *models.TriviaBundle) error {
	return s.r.UpdateTriviaBundle(tb)
}

func (s *triviaBundleService) DeleteTriviaBundle(id int) error {
	return s.r.DeleteTriviaBundle(id)
}

func (s *triviaBundleService) GetTriviaBundles() ([]*models.TriviaBundle, error) {
	return s.r.GetTriviaBundles()
}

func (s *triviaBundleService) GetTriviaBundlesByCategory(category string) ([]*models.TriviaBundle, error) {
	return s.r.GetTriviaBundlesByCategory(category)
}
