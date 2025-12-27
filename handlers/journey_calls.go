package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

const (
	DefaultCallsLimit = 200
	MaxCallsLimit     = 1000
)

type JourneyCallsResponse struct {
	Journey OrderedRow        `json:"journey"`
	Calls   []OrderedRow      `json:"calls"`
	Count   int               `json:"count"`
	Params  map[string]string `json:"params"`
}

func (h *JourneyHandler) GetJourneyCalls(w http.ResponseWriter, r *http.Request) {
	lrw, logResponse := LogRequestWithWriter(w, r)
	defer logResponse()

	lrw.Header().Set("Content-Type", "application/json")

	// Parse query parameters
	query := r.URL.Query()
	idStr := query.Get("id")

	// Collect params for response
	params := make(map[string]string)
	for key, values := range query {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	// Validate id parameter
	if idStr == "" {
		lrw.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(lrw).Encode(map[string]interface{}{
			"error": "Missing required parameter: id",
			"code":  400,
		})
		return
	}

	// Parse id as integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		lrw.WriteHeader(http.StatusNotFound)
		json.NewEncoder(lrw).Encode(map[string]interface{}{
			"error": "Invalid id parameter: must be an integer",
			"code":  404,
		})
		return
	}

	// First, get the journey
	journeyQuery := "SELECT * FROM estimatedvehiclejourney WHERE id = $1"
	journeyRows, err := h.DB.Query(journeyQuery, id)
	if err != nil {
		lrw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(lrw).Encode(map[string]interface{}{
			"error": "Database query failed",
			"code":  500,
		})
		return
	}
	defer journeyRows.Close()

	journeyResults, err := ScanRowsToOrderedMaps(journeyRows)
	if err != nil {
		lrw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(lrw).Encode(map[string]interface{}{
			"error": "Failed to scan journey",
			"code":  500,
		})
		return
	}

	if len(journeyResults) == 0 {
		lrw.WriteHeader(http.StatusNotFound)
		json.NewEncoder(lrw).Encode(map[string]interface{}{
			"error": "Journey not found",
			"code":  404,
		})
		return
	}

	journey := journeyResults[0]

	// Parse limit for calls
	limit := ParseLimit(query.Get("limit"), DefaultCallsLimit, MaxCallsLimit)

	// Get the calls
	callsQuery := `SELECT * FROM calls WHERE estimatedvehiclejourney = $1 ORDER BY "order" ASC LIMIT $2`
	callsRows, err := h.DB.Query(callsQuery, id, limit)
	if err != nil {
		lrw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(lrw).Encode(map[string]interface{}{
			"error": "Failed to query calls",
			"code":  500,
		})
		return
	}
	defer callsRows.Close()

	callsResults, err := ScanRowsToOrderedMaps(callsRows)
	if err != nil {
		lrw.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(lrw).Encode(map[string]interface{}{
			"error": "Failed to scan calls",
			"code":  500,
		})
		return
	}

	// Build response
	response := JourneyCallsResponse{
		Journey: journey,
		Calls:   callsResults,
		Count:   len(callsResults),
		Params:  params,
	}

	// Return JSON response
	json.NewEncoder(lrw).Encode(response)
}
