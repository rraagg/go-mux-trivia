package database

import (
	"context"
	"database/sql"
	"fmt"
	"go-mux-trivia/internal/models"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Health() map[string]string
	GetTriviaBundle(id int) (*models.TriviaBundle, error)
	GetTriviaBundlesByCategory(category string) ([]*models.TriviaBundle, error)
	GetTriviaBundles() ([]*models.TriviaBundle, error)
	CreateTriviaBundle(triviaBundle *models.TriviaBundle) (int, error)
	UpdateTriviaBundle(triviaBundle *models.TriviaBundle) error
	DeleteTriviaBundle(id int) error
}

type service struct {
	db *sql.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() Service {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username,
		password,
		host,
		port,
		database,
	)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS trivia_bundle (
			id SERIAL PRIMARY KEY,
			question TEXT NOT NULL,
			category TEXT NOT NULL,
			show_answer BOOLEAN NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS answer (
			id SERIAL PRIMARY KEY,
			trivia_bundle_id INTEGER REFERENCES trivia_bundle(id),
			answer_text TEXT NOT NULL,
			is_correct BOOLEAN NOT NULL
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	s := &service{db: db}
	return s
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) GetTriviaBundle(id int) (*models.TriviaBundle, error) {
	var tb models.TriviaBundle
	err := s.db.QueryRow("SELECT * FROM trivia_bundle WHERE id = $1", id).
		Scan(&tb.ID, &tb.Question, &tb.Category, &tb.ShowAnswer)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query("SELECT * FROM answer WHERE trivia_bundle_id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a models.Answer
		err := rows.Scan(&a.ID, &a.TriviaBundleID, &a.AnswerText, &a.IsCorrect)
		if err != nil {
			return nil, err
		}
		tb.Answers = append(tb.Answers, a)
	}

	return &tb, nil
}

func (s *service) GetTriviaBundlesByCategory(category string) ([]*models.TriviaBundle, error) {
	rows, err := s.db.Query("SELECT * FROM trivia_bundle WHERE category = $1", category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bundles []*models.TriviaBundle
	for rows.Next() {
		var tb models.TriviaBundle
		err := rows.Scan(
			&tb.ID,
			&tb.Question,
			&tb.Category,
			&tb.ShowAnswer,
		)
		if err != nil {
			return nil, err
		}
		bundles = append(bundles, &tb)
	}
	return bundles, nil
}

func (s *service) GetTriviaBundles() ([]*models.TriviaBundle, error) {
	rows, err := s.db.Query("SELECT * FROM trivia_bundle")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bundles []*models.TriviaBundle
	for rows.Next() {
		var tb models.TriviaBundle
		err := rows.Scan(
			&tb.ID,
			&tb.Question,
			&tb.Category,
			&tb.ShowAnswer,
		)
		if err != nil {
			return nil, err
		}
		bundles = append(bundles, &tb)
	}
	return bundles, nil
}

func (s *service) CreateTriviaBundle(triviaBundle *models.TriviaBundle) (int, error) {
	log.Default().Printf("Creating trivia bundle: %s", triviaBundle.Question)
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	err = tx.QueryRow(
		"INSERT INTO trivia_bundle (question, category, show_answer) VALUES ($1, $2, $3) RETURNING id",
		triviaBundle.Question,
		triviaBundle.Category,
		triviaBundle.ShowAnswer,
	).
		Scan(&id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, answer := range triviaBundle.Answers {
		_, err := tx.Exec(
			"INSERT INTO answer (trivia_bundle_id, answer_text, is_correct) VALUES ($1, $2, $3) RETURNING id",
			id,
			answer.AnswerText,
			answer.IsCorrect)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *service) UpdateTriviaBundle(triviaBundle *models.TriviaBundle) error {
	_, err := s.db.Exec(
		"UPDATE trivia_bundle SET question = $1, category = $2, show_answer = $3 WHERE id = $7",
		triviaBundle.Question,
		triviaBundle.Category,
		triviaBundle.ShowAnswer,
		triviaBundle.ID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) DeleteTriviaBundle(id int) error {
	_, err := s.db.Exec("DELETE FROM trivia_bundle WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Close() error {
	return s.db.Close()
}
