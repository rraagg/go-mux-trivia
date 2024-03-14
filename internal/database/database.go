package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

var (
	db_name  = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() *sql.DB {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username,
		password,
		host,
		port,
		db_name,
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

	return db
}
