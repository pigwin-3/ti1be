package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"ti1be/config"
	"ti1be/handlers"
	"ti1be/pages"

	// fjern senere
	"github.com/joho/godotenv"
)

type StatusResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

var db *sql.DB

func statusHandler(w http.ResponseWriter, r *http.Request) {
	lrw, logResponse := handlers.LogRequestWithWriter(w, r)
	defer logResponse()

	lrw.Header().Set("Content-Type", "application/json")
	response := StatusResponse{Status: "running"}
	json.NewEncoder(lrw).Encode(response)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	lrw, logResponse := handlers.LogRequestWithWriter(w, r)
	defer logResponse()

	lrw.Header().Set("Content-Type", "application/json")
	lrw.WriteHeader(http.StatusNotFound)
	response := ErrorResponse{
		Error: "Not Found",
		Code:  404,
	}
	json.NewEncoder(lrw).Encode(response)
}

func main() {
	// fjern senere
	godotenv.Load()

	// Connect to PostgreSQL
	var err error
	db, err = config.ConnectToPostgreSQL()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer config.DisconnectFromPostgreSQL(db)

	// Create handlers
	journeyHandler := &handlers.JourneyHandler{DB: db}
	callsHandler := &handlers.CallsHandler{DB: db}

	mux := http.NewServeMux()

	// Serve home page
	mux.HandleFunc("/", pages.HomeHandler)

	// Serve favicon
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/favicon.ico")
	})

	// Register API v1 endpoints
	mux.HandleFunc("/api/v1/status", statusHandler)
	mux.HandleFunc("/api/v1/journey/get", journeyHandler.GetJourneys)
	mux.HandleFunc("/api/v1/journey/calls", journeyHandler.GetJourneyCalls)
	mux.HandleFunc("/api/v1/calls/get", callsHandler.GetCalls)

	// Use mux directly (it handles 404s for unregistered routes)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
