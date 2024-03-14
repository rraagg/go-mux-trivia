package server

import (
	"encoding/json"
	"go-mux-trivia/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", s.HelloWorldHandler)

	r.HandleFunc("/health", s.healthHandler)

	r.HandleFunc("/trivia/{id}", s.GetTriviaBundle).Methods("GET")

	r.HandleFunc("/create", s.CreateTriviaBundle).Methods("POST")

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(s.triviaBundleRepostiory.Health())

	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) GetTriviaBundle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	// Convert the id from string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call the service layer to get the TriviaBundle
	triviaBundle, err := s.triviaBundleService.GetTriviaBundle(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(triviaBundle)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) CreateTriviaBundle(w http.ResponseWriter, r *http.Request) {
	var tb models.TriviaBundle
	err := json.NewDecoder(r.Body).Decode(&tb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Default().Printf("Handler: Creating TriviaBundle: %s", tb.Question)

	id, err := s.triviaBundleService.CreateTriviaBundle(&tb)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(strconv.Itoa(id)))
}
