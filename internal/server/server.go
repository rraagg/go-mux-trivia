package server

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"go-mux-trivia/internal/database"
	"go-mux-trivia/internal/repository"
	"go-mux-trivia/internal/service"
)

type Server struct {
	port                   int
	db                     *sql.DB
	templates              *template.Template
	triviaBundleService    service.TriviaBundleService
	triviaBundleRepostiory repository.TriviaRepository
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:                   port,
		db:                     database.New(),
		templates:              template.Must(template.ParseGlob("./internal/templates/*")),
		triviaBundleService:    service.NewTriviaBundleService(),
		triviaBundleRepostiory: repository.NewTriviaRepository(database.New()),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
