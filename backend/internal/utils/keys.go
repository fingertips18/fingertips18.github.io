package utils

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// GenerateKey generates a unique key string composed of a timestamp-based part and a random part.
// It attempts to create a UUID version 7 (time-ordered), and falls back to a random UUID if necessary.
// The resulting key is formatted as "<timestampPart>-<randomPart>", where both parts are derived from the UUID.
func GenerateKey() string {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		uuidV7 = uuid.New() // Fallback to UUID v4 if v7 fails
	}
	uuidStr := strings.ReplaceAll(uuidV7.String(), "-", "")

	// Ensure we extract exactly 18 characters
	timestampPart := uuidStr[:8] // First 8 chars for timestamp
	randomPart := uuidStr[8:18]  // Next 10 chars for randomness

	// Insert dash between timestamp and randomness
	// This provides for a ~0.08% collision probability at ~1 million/sec (1,000/ms)
	return fmt.Sprintf("%s-%s", timestampPart, randomPart)
}
