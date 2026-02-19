package repository

import (
	"strings"

	"github.com/lib/pq"
)

const (
	uniqueViolation = "23505"
)

func isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == uniqueViolation
	}
	return strings.Contains(err.Error(), "duplicate key")
}
