package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

const (
	DefaultLimit = 50
	MaxLimit     = 1000
)

type JourneyHandler struct {
	DB *sql.DB
}

type JourneyResponse struct {
	Data   []OrderedRow      `json:"data"`
	Count  int               `json:"count"`
	Params map[string]string `json:"params"`
}

func (h *JourneyHandler) GetJourneys(w http.ResponseWriter, r *http.Request) {
	// Log request
	clientIP := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}
	log.Printf("GET %s from %s", r.URL.RequestURI(), clientIP)

	w.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	query := r.URL.Query()

	// Collect params for response
	params := make(map[string]string)
	for key, values := range query {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	// Build SQL query using QueryBuilder
	qb := NewQueryBuilder("SELECT * FROM public.estimatedvehiclejourney WHERE 1=1")

	// Add conditions
	if id := query.Get("id"); id != "" {
		qb.AddCondition("id", id)
	}
	qb.AddSingleOrMultipleCondition("vehicleref", query.Get("vehicle_ref"))
	qb.AddSingleOrMultipleCondition("datasource", query.Get("data_source"))
	qb.AddSingleOrMultipleCondition("lineref", query.Get("line_ref"))

	if after := query.Get("after"); after != "" {
		qb.AddComparisonCondition("id", "<", after)
	}

	// Add ORDER BY
	qb.Query += " ORDER BY id DESC"

	// Add LIMIT
	limit := ParseLimit(query.Get("limit"), DefaultLimit, MaxLimit)
	qb.AddLimit(limit)

	// Execute query
	rows, err := h.DB.Query(qb.Query, qb.Args...)
	if err != nil {
		http.Error(w, `{"error":"Database query failed","code":500}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Scan rows
	results, err := ScanRowsToOrderedMaps(rows)
	if err != nil {
		http.Error(w, `{"error":"Failed to scan rows","code":500}`, http.StatusInternalServerError)
		return
	}

	// Build response
	response := JourneyResponse{
		Data:   results,
		Count:  len(results),
		Params: params,
	}

	// Return JSON response
	json.NewEncoder(w).Encode(response)
}
