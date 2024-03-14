package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-mux-trivia/internal/models"
	"log"
	"time"
)

type TriviaRepository interface {
	GetTriviaBundle(id int) (*models.TriviaBundle, error)
	GetTriviaBundlesByCategory(category string) ([]*models.TriviaBundle, error)
	GetTriviaBundles() ([]*models.TriviaBundle, error)
	CreateTriviaBundle(triviaBundle *models.TriviaBundle) (int, error)
	UpdateTriviaBundle(triviaBundle *models.TriviaBundle) error
	DeleteTriviaBundle(id int) error
	Health() map[string]string
}

type triviaRepository struct {
	db *sql.DB
}

func NewTriviaRepository(db *sql.DB) TriviaRepository {
	return &triviaRepository{
		db: db,
	}
}

func (t *triviaRepository) GetTriviaBundle(id int) (*models.TriviaBundle, error) {
	var tb models.TriviaBundle
	err := t.db.QueryRow("SELECT * FROM trivia_bundle WHERE id = $1", id).
		Scan(&tb.ID, &tb.Question, &tb.Category, &tb.ShowAnswer)
	if err != nil {
		return nil, err
	}

	rows, err := t.db.Query("SELECT * FROM answer WHERE trivia_bundle_id = $1", id)
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

func (t *triviaRepository) GetTriviaBundlesByCategory(category string) ([]*models.TriviaBundle, error) {
	rows, err := t.db.Query("SELECT * FROM trivia_bundle WHERE category = $1", category)
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

func (t *triviaRepository) GetTriviaBundles() ([]*models.TriviaBundle, error) {
	rows, err := t.db.Query("SELECT * FROM trivia_bundle")
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

func (t *triviaRepository) CreateTriviaBundle(triviaBundle *models.TriviaBundle) (int, error) {
	log.Default().Printf("Creating trivia bundle: %s", triviaBundle.Question)
	tx, err := t.db.Begin()
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

func (t *triviaRepository) UpdateTriviaBundle(triviaBundle *models.TriviaBundle) error {
	_, err := t.db.Exec(
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

func (t *triviaRepository) DeleteTriviaBundle(id int) error {
	_, err := t.db.Exec("DELETE FROM trivia_bundle WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (t *triviaRepository) Close() error {
	return t.db.Close()
}

func (t *triviaRepository) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := t.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}
