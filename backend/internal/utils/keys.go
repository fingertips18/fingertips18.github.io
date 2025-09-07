package utils

import (
	"github.com/google/uuid"
)

// GenerateKey generates a new UUID string. It attempts to create a UUID version 7,
// and falls back to a version 4 UUID if version 7 generation fails.
// The returned string is a standard 36-character UUID.
func GenerateKey() string {
	key, err := uuid.NewV7()
	if err != nil {
		key = uuid.New() // fallback to v4
	}
	return key.String() // full 36-character UUID
}
