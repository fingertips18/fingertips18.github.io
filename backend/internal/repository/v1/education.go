package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/database"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/utils"
	"github.com/jackc/pgx/v5"
)

type EducationRepository interface {
	Create(ctx context.Context, education *domain.Education) (string, error)
	Get(ctx context.Context, id string) (*domain.Education, error)
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
func (r *educationRepository) Create(ctx context.Context, education *domain.Education) (string, error) {
	if err := education.ValidatePayload(); err != nil {
		return "", fmt.Errorf("failed to validate education: %w", err)
	}

	id := utils.GenerateKey()
	now := r.timeProvider()

	mainSchoolJSON, _ := json.Marshal(education.MainSchool)
	schoolPeriodsJSON, _ := json.Marshal(education.SchoolPeriods)
	projectsJSON, _ := json.Marshal(education.Projects)

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
		mainSchoolJSON,
		schoolPeriodsJSON,
		projectsJSON,
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

func (r *educationRepository) Get(ctx context.Context, id string) (*domain.Education, error) {
	if id == "" {
		return nil, fmt.Errorf("failed to get education: ID missing")
	}

	var education domain.Education

	var mainSchoolJSON []byte
	var schoolPeriodsJSON []byte
	var projectsJSON []byte

	query := fmt.Sprintf(
		`SELECT id, main_school, school_periods, projects, level, created_at, updated_at
		FROM %s
		WHERE id = $1`,
		r.educationTable,
	)

	err := r.databaseAPI.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&education.Id,
		&mainSchoolJSON,
		&schoolPeriodsJSON,
		&projectsJSON,
		&education.Level,
		&education.CreatedAt,
		&education.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("education not found: %w", err)
		}
		return nil, fmt.Errorf("failed to scan education: %w", err)
	}

	if err = json.Unmarshal(mainSchoolJSON, &education.MainSchool); err != nil {
		return nil, fmt.Errorf("failed to unmarshal main school: %w", err)
	}
	if len(schoolPeriodsJSON) > 0 {
		if err = json.Unmarshal(schoolPeriodsJSON, &education.SchoolPeriods); err != nil {
			return nil, fmt.Errorf("failed to unmarshal school periods: %w", err)
		}
	}
	if len(projectsJSON) > 0 {
		if err = json.Unmarshal(projectsJSON, &education.Projects); err != nil {
			return nil, fmt.Errorf("failed to unmarshal projects: %w", err)
		}
	}

	if err := education.ValidateResponse(); err != nil {
		return nil, fmt.Errorf("invalid education returned: %w", err)
	}

	return &education, nil
}
