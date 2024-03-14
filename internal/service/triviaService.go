package service

import (
	"go-mux-trivia/internal/database"
	"go-mux-trivia/internal/models"
	"go-mux-trivia/internal/repository"
)

type TriviaService interface {
	GetTriviaBundle(id int) (*models.TriviaBundle, error)
	GetTriviaBundlesByCategory(category string) ([]*models.TriviaBundle, error)
	GetTriviaBundles() ([]*models.TriviaBundle, error)
	CreateTriviaBundle(triviaBundle *models.TriviaBundle) (int, error)
	UpdateTriviaBundle(triviaBundle *models.TriviaBundle) error
	DeleteTriviaBundle(id int) error
}

type triviaService struct {
	r repository.TriviaRepository
}

func NewTriviaService() *triviaService {
	return &triviaService{
		r: repository.NewTriviaRepository(database.New()),
	}
}

func (s *triviaService) GetTriviaBundle(id int) (*models.TriviaBundle, error) {
	bundle, err := s.r.GetTriviaBundle(id)
	return bundle, err
}

func (s *triviaService) CreateTriviaBundle(tb *models.TriviaBundle) (int, error) {
	return s.r.CreateTriviaBundle(tb)
}

func (s *triviaService) UpdateTriviaBundle(tb *models.TriviaBundle) error {
	return s.r.UpdateTriviaBundle(tb)
}

func (s *triviaService) DeleteTriviaBundle(id int) error {
	return s.r.DeleteTriviaBundle(id)
}

func (s *triviaService) GetTriviaBundles() ([]*models.TriviaBundle, error) {
	return s.r.GetTriviaBundles()
}

func (s *triviaService) GetTriviaBundlesByCategory(category string) ([]*models.TriviaBundle, error) {
	return s.r.GetTriviaBundlesByCategory(category)
}
