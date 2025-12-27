package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

const (
	DefaultCallsGetLimit = 10
	MaxCallsGetLimit     = 1000
)

type CallsHandler struct {
	DB *sql.DB
}

type CallsResponse struct {
	Data   []OrderedRow      `json:"data"`
	Count  int               `json:"count"`
	Params map[string]string `json:"params"`
}

func (h *CallsHandler) GetCalls(w http.ResponseWriter, r *http.Request) {
	lrw, logResponse := LogRequestWithWriter(w, r)
	defer logResponse()

	lrw.Header().Set("Content-Type", "application/json")

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
	qb := NewQueryBuilder("SELECT * FROM calls WHERE 1=1")

	// Add conditions
	qb.AddSingleOrMultipleCondition("id", query.Get("id"))
	qb.AddSingleOrMultipleCondition("estimatedvehiclejourney", query.Get("estimatedvehiclejourney"))

	// Handle "order" field with quotes since it's a reserved keyword
	orderValue := query.Get("order")
	if orderValue != "" {
		// Use quoted field name for order
		qb.AddSingleOrMultipleConditionWithQuotes("order", orderValue)
	}

	qb.AddSingleOrMultipleCondition("stoppointref", query.Get("stoppointref"))

	// Add ORDER BY
	qb.Query += " ORDER BY id ASC"

	// Add LIMIT
	limit := ParseLimit(query.Get("limit"), DefaultCallsGetLimit, MaxCallsGetLimit)
	qb.AddLimit(limit)

	// Execute query
	rows, err := h.DB.Query(qb.Query, qb.Args...)
	if err != nil {
		http.Error(lrw, `{"error":"Database query failed","code":500}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Scan rows
	results, err := ScanRowsToOrderedMaps(rows)
	if err != nil {
		http.Error(lrw, `{"error":"Failed to scan rows","code":500}`, http.StatusInternalServerError)
		return
	}

	// Build response
	response := CallsResponse{
		Data:   results,
		Count:  len(results),
		Params: params,
	}

	// Return JSON response
	json.NewEncoder(lrw).Encode(response)
}
