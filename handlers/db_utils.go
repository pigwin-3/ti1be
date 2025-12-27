package handlers

import (
	"database/sql"
	"encoding/json"
	"strings"
)

// OrderedRow ensures 'id' is always first in JSON output
type OrderedRow map[string]interface{}

func (o OrderedRow) MarshalJSON() ([]byte, error) {
	// Manually build JSON with id first
	var buf strings.Builder
	buf.WriteString("{")

	// Write id first if it exists
	first := true
	if id, ok := o["id"]; ok {
		idBytes, err := json.Marshal(id)
		if err != nil {
			return nil, err
		}
		buf.WriteString(`"id":`)
		buf.Write(idBytes)
		first = false
	}

	// Write all other fields alphabetically
	keys := make([]string, 0, len(o))
	for k := range o {
		if k != "id" {
			keys = append(keys, k)
		}
	}

	// Sort remaining keys
	sortKeys(keys)

	for _, k := range keys {
		if !first {
			buf.WriteString(",")
		}
		first = false

		keyBytes, err := json.Marshal(k)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)
		buf.WriteString(":")

		valBytes, err := json.Marshal(o[k])
		if err != nil {
			return nil, err
		}
		buf.Write(valBytes)
	}

	buf.WriteString("}")
	return []byte(buf.String()), nil
}

// Simple string sort helper
func sortKeys(keys []string) {
	// Bubble sort (simple for small arrays)
	n := len(keys)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if keys[j] > keys[j+1] {
				keys[j], keys[j+1] = keys[j+1], keys[j]
			}
		}
	}
}

// ScanRowsToOrderedMaps scans SQL rows into a slice of OrderedRow maps
func ScanRowsToOrderedMaps(rows *sql.Rows) ([]OrderedRow, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := []OrderedRow{}

	for rows.Next() {
		// Create a slice of interface{}'s to represent each column
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the result into the column pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Create a map for this row
		rowMap := make(OrderedRow)
		for i, col := range columns {
			val := values[i]
			// Handle different types
			switch v := val.(type) {
			case []byte:
				// Try to parse as JSON if it looks like JSON
				str := string(v)
				if (strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}")) ||
					(strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]")) {
					var jsonVal interface{}
					if err := json.Unmarshal(v, &jsonVal); err == nil {
						rowMap[col] = jsonVal
					} else {
						rowMap[col] = str
					}
				} else {
					rowMap[col] = str
				}
			case nil:
				rowMap[col] = nil
			default:
				// Keep original type (int64, float64, bool, string, etc.)
				rowMap[col] = val
			}
		}
		results = append(results, rowMap)
	}

	return results, rows.Err()
}
