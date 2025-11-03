package utils

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateKey_V7Success(t *testing.T) {
	// Override newV7 to always return a fixed UUID
	expected := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	newV7 = func() (uuid.UUID, error) {
		return expected, nil
	}
	defer func() { newV7 = uuid.NewV7 }() // Restore after test

	got := GenerateKey()
	assert.Equal(t, expected.String(), got)
}

func TestGenerateKey_V7FailsFallbackToV4(t *testing.T) {
	newV7 = func() (uuid.UUID, error) {
		return uuid.Nil, errors.New("force fail")
	}
	defer func() { newV7 = uuid.NewV7 }()

	got := GenerateKey()
	_, err := uuid.Parse(got)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil.String(), got)
}
