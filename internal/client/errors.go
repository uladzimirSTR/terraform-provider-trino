package client

import "fmt"

type QueryError struct {
	Message   string
	SQLState  string
	ErrorCode int
	ErrorName string
	ErrorType string
}

func (e QueryError) Error() string {
	return fmt.Sprintf(
		"trino query failed: %s [%s, code=%d, type=%s, sqlstate=%s]",
		e.Message,
		e.ErrorName,
		e.ErrorCode,
		e.ErrorType,
		e.SQLState,
	)
}
