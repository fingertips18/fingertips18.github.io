package v1

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/database"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/utils"
)

type EducationRepository interface {
	Create(ctx context.Context, education domain.Education) (string, error)
}

type EducationRepositoryConfig struct {
	DatabaseAPI    database.DatabaseAPI
	EducationTable string

	timeProvider domain.TimeProvider
}

type educationRepository struct {
	educationTable string
	databaseAPI    database.DatabaseAPI
	timeProvider   domain.TimeProvider
}

// NewEducationRepository creates and returns a configured EducationRepository.
//
// It accepts an EducationRepositoryConfig and constructs an internal
// educationRepository backed by cfg.EducationTable and cfg.DatabaseAPI.
// If cfg.timeProvider is nil, the repository defaults to using time.Now
// as the time provider. The returned value implements the
// EducationRepository interface and is never nil.
func NewEducationRepository(cfg EducationRepositoryConfig) EducationRepository {
	timeProvider := cfg.timeProvider
	if timeProvider == nil {
		timeProvider = time.Now
	}

	return &educationRepository{
		educationTable: cfg.EducationTable,
		databaseAPI:    cfg.DatabaseAPI,
		timeProvider:   timeProvider,
	}
}

// Create validates and persists a new Education record in the repository.
// It performs the following steps:
//   - Validates the provided domain.Education via education.ValidatePayload() and returns an error if validation fails.
//   - Generates a new unique ID (via utils.GenerateKey) for the record.
//   - Sets CreatedAt and UpdatedAt to the current time obtained from the repository's timeProvider.
//   - Executes an INSERT against the configured education table using the repository's databaseAPI and expects the database to RETURN the inserted id.
//
// On success, Create returns the newly created record's ID. It returns an error if validation fails, if the database INSERT or Scan fails, or if the database returns an empty ID.
func (r *educationRepository) Create(ctx context.Context, education domain.Education) (string, error) {
	if err := education.ValidatePayload(); err != nil {
		return "", fmt.Errorf("failed to validate education: %w", err)
	}

	id := utils.GenerateKey()
	now := r.timeProvider()

	education.CreatedAt = now
	education.UpdatedAt = now

	query := fmt.Sprintf(
		`INSERT INTO %s
        (id, main_school, school_periods, projects, level, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`,
		r.educationTable,
	)

	var returnedID string
	err := r.databaseAPI.QueryRow(
		ctx,
		query,
		id,
		education.MainSchool,
		education.SchoolPeriods,
		education.Projects,
		education.Level,
		education.CreatedAt,
		education.UpdatedAt,
	).Scan(&returnedID)

	if err != nil {
		return "", fmt.Errorf("failed to create education: %w", err)
	}

	if returnedID == "" {
		return "", errors.New("invalid education returned: ID missing")
	}

	return returnedID, nil
}
