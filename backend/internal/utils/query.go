package utils

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
)

// GetQueryInt32 retrieves the value associated with the given key from the provided url.Values,
// attempts to convert it to an int32, and returns the result. If the key is not present or the
// value cannot be converted to an integer, the provided default value (def) is returned.
func GetQueryInt32(q url.Values, key string, def int32) int32 {
	v := q.Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return def
	}
	return int32(n)
}

// GetQuerySortBy retrieves and validates a sort-by parameter from the provided URL query values.
// It expects the value associated with the given key to match one of the allowed domain.SortBy values
// (e.g., domain.CreatedAt, domain.UpdatedAt). If the key is not present or the value is empty, it returns nil.
// If the value is valid, it returns a pointer to the corresponding domain.SortBy value.
// Otherwise, it returns an error indicating an invalid sort-by value.
func GetQuerySortBy(q url.Values, key string) (*domain.SortBy, error) {
	v := q.Get(key)
	if v == "" {
		return nil, nil
	}

	s := domain.SortBy(v)
	switch s {
	case domain.CreatedAt, domain.UpdatedAt:
		return &s, nil
	default:
		return nil, fmt.Errorf("invalid sort by value: %q", v)
	}
}

// GetQueryBool retrieves a boolean value from the provided url.Values map using the specified key.
// If the key is not present or the value cannot be parsed as a boolean, the default value 'def' is returned.
// Accepted boolean values are as defined by strconv.ParseBool (e.g., "1", "t", "T", "TRUE", "true", "True" for true).
func GetQueryBool(q url.Values, key string, def bool) bool {
	v := q.Get(key)
	if v == "" {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}
