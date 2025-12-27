package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"ti1be/config"
	"ti1be/handlers"

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
	w.Header().Set("Content-Type", "application/json")
	response := StatusResponse{Status: "ok"}
	json.NewEncoder(w).Encode(response)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	response := ErrorResponse{
		Error: "Not Found",
		Code:  404,
	}
	json.NewEncoder(w).Encode(response)
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

	mux := http.NewServeMux()

	// Register status endpoint
	mux.HandleFunc("/status", statusHandler)

	// Register journey endpoints
	mux.HandleFunc("/journey/get", journeyHandler.GetJourneys)

	// Wrap with custom 404 handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if path matches any registered routes
		if r.URL.Path != "/status" && !strings.HasPrefix(r.URL.Path, "/journey/") {
			notFoundHandler(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
