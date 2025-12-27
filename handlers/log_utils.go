package handlers

import (
	"log"
	"net/http"
	"time"
)

// ResponseWriter wrapper to capture status code
type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// LogRequest logs the HTTP request with client IP and URL
// Returns a function to log the response details
func LogRequest(r *http.Request) func(statusCode int, duration time.Duration) {
	clientIP := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}

	logResponse := func(statusCode int, duration time.Duration) {
		log.Printf("%s %s from %s - %d [%v]", r.Method, r.URL.RequestURI(), clientIP, statusCode, duration)
	}

	return logResponse
}

// LogRequestWithWriter wraps the ResponseWriter and returns both the wrapper and a cleanup function
func LogRequestWithWriter(w http.ResponseWriter, r *http.Request) (*LoggingResponseWriter, func()) {
	clientIP := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}

	startTime := time.Now()
	lrw := NewLoggingResponseWriter(w)

	cleanup := func() {
		duration := time.Since(startTime)
		log.Printf("%s %s from %s - %d [%v]", r.Method, r.URL.RequestURI(), clientIP, lrw.StatusCode, duration)
	}

	return lrw, cleanup
}
