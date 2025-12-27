package handlers

import (
	"strconv"
	"strings"
)

// QueryBuilder helps build SQL queries with parameterized values
type QueryBuilder struct {
	Query    string
	Args     []interface{}
	ArgCount int
}

// NewQueryBuilder creates a new query builder with a base query
func NewQueryBuilder(baseQuery string) *QueryBuilder {
	return &QueryBuilder{
		Query:    baseQuery,
		Args:     []interface{}{},
		ArgCount: 1,
	}
}

// AddCondition adds a simple WHERE condition
func (qb *QueryBuilder) AddCondition(field string, value interface{}) {
	qb.Query += " AND " + field + " = $" + strconv.Itoa(qb.ArgCount)
	qb.Args = append(qb.Args, value)
	qb.ArgCount++
}

// AddComparisonCondition adds a WHERE condition with a custom operator
func (qb *QueryBuilder) AddComparisonCondition(field string, operator string, value interface{}) {
	qb.Query += " AND " + field + " " + operator + " $" + strconv.Itoa(qb.ArgCount)
	qb.Args = append(qb.Args, value)
	qb.ArgCount++
}

// AddInCondition adds an IN condition with multiple values
func (qb *QueryBuilder) AddInCondition(field string, values []string) {
	placeholders := []string{}
	for _, val := range values {
		placeholders = append(placeholders, "$"+strconv.Itoa(qb.ArgCount))
		qb.Args = append(qb.Args, strings.TrimSpace(val))
		qb.ArgCount++
	}
	qb.Query += " AND " + field + " IN (" + strings.Join(placeholders, ",") + ")"
}

// AddSingleOrMultipleCondition handles parameters that can be single or comma-separated
func (qb *QueryBuilder) AddSingleOrMultipleCondition(field string, value string) {
	if value == "" {
		return
	}

	values := strings.Split(value, ",")
	if len(values) == 1 {
		qb.AddCondition(field, value)
	} else {
		qb.AddInCondition(field, values)
	}
}

// AddSingleOrMultipleConditionWithQuotes handles parameters that need quoted field names (like "order")
func (qb *QueryBuilder) AddSingleOrMultipleConditionWithQuotes(field string, value string) {
	if value == "" {
		return
	}

	quotedField := `"` + field + `"`
	values := strings.Split(value, ",")
	if len(values) == 1 {
		qb.Query += " AND " + quotedField + " = $" + strconv.Itoa(qb.ArgCount)
		qb.Args = append(qb.Args, value)
		qb.ArgCount++
	} else {
		placeholders := []string{}
		for _, val := range values {
			placeholders = append(placeholders, "$"+strconv.Itoa(qb.ArgCount))
			qb.Args = append(qb.Args, strings.TrimSpace(val))
			qb.ArgCount++
		}
		qb.Query += " AND " + quotedField + " IN (" + strings.Join(placeholders, ",") + ")"
	}
}

// AddLimit adds a LIMIT clause
func (qb *QueryBuilder) AddLimit(limit int) {
	qb.Query += " LIMIT $" + strconv.Itoa(qb.ArgCount)
	qb.Args = append(qb.Args, limit)
	qb.ArgCount++
}

// ParseLimit parses a limit string with default and max values
func ParseLimit(limitStr string, defaultLimit, maxLimit int) int {
	if limitStr == "" {
		return defaultLimit
	}

	if parsedLimit, err := strconv.Atoi(limitStr); err == nil {
		if parsedLimit > maxLimit {
			return maxLimit
		}
		if parsedLimit < 1 {
			return defaultLimit
		}
		return parsedLimit
	}

	return defaultLimit
}
