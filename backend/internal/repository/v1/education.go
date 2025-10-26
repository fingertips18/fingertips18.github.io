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
	Update(ctx context.Context, education *domain.Education) (*domain.Education, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter domain.EducationFilter) ([]domain.Education, error)
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

// Create creates a new education record in the repository and returns its generated ID.
// It validates the provided Education payload, generates a unique ID, marshals the main
// school and school periods to JSON, sets CreatedAt and UpdatedAt timestamps from the
// repository's time provider, and inserts the record into the configured education table.
// The function returns the newly created record's ID, or a non-nil error if validation,
// JSON marshaling, database insertion, or returned-ID verification fails. The provided
// education object's CreatedAt and UpdatedAt fields are updated when the method succeeds.
func (r *educationRepository) Create(ctx context.Context, education *domain.Education) (string, error) {
	if education == nil {
		return "", errors.New("failed to validate education: payload is nil")
	}

	if err := education.ValidatePayload(); err != nil {
		return "", fmt.Errorf("failed to validate education: %w", err)
	}

	id := utils.GenerateKey()
	now := r.timeProvider()

	mainSchoolJSON, err := json.Marshal(education.MainSchool)
	if err != nil {
		return "", fmt.Errorf("failed to marshal main school: %w", err)
	}
	schoolPeriodsJSON, err := json.Marshal(education.SchoolPeriods)
	if err != nil {
		return "", fmt.Errorf("failed to marshal school periods: %w", err)
	}

	education.CreatedAt = now
	education.UpdatedAt = now

	query := fmt.Sprintf(
		`INSERT INTO %s
        (id, main_school, school_periods, level, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`,
		r.educationTable,
	)

	var returnedID string
	err = r.databaseAPI.QueryRow(
		ctx,
		query,
		id,
		mainSchoolJSON,
		schoolPeriodsJSON,
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

// Get retrieves an Education by its id from the repository.
// It returns a pointer to domain.Education on success.
//
// Get requires a non-empty id and will return an error if id == "".
// It uses the provided context for the database query so callers can control
// cancellation and timeouts.
//
// Expected database columns are: id, main_school, school_periods, level, created_at, updated_at.
// The main_school and school_periods columns are stored as JSON: main_school is unmarshaled
// into Education.MainSchool (required), and school_periods is unmarshaled into
// Education.SchoolPeriods only if present (empty JSON is allowed).
//
// Possible errors returned include:
//   - an error indicating the ID is missing when id == "".
//   - a not-found error (pgx.ErrNoRows) when no matching row exists.
//   - errors from scanning database row values.
//   - JSON unmarshal errors for main_school or school_periods.
//   - an error if the resulting domain.Education fails ValidateResponse().
//
// All returned errors are wrapped with contextual messages to aid debugging.
func (r *educationRepository) Get(ctx context.Context, id string) (*domain.Education, error) {
	if id == "" {
		return nil, fmt.Errorf("failed to get education: ID missing")
	}

	var (
		education                         domain.Education
		mainSchoolJSON, schoolPeriodsJSON []byte
	)

	query := fmt.Sprintf(
		`SELECT id, main_school, school_periods, level, created_at, updated_at
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

	if err := education.ValidateResponse(); err != nil {
		return nil, fmt.Errorf("invalid education returned: %w", err)
	}

	return &education, nil
}

// Update validates and persists changes to an existing Education record.
//
// It performs the following steps:
//   - Validates the provided education payload via ValidatePayload.
//   - Sets the UpdatedAt timestamp using the repository's time provider.
//   - Marshals JSON-serializable fields (MainSchool, SchoolPeriods, Projects).
//   - Executes an SQL UPDATE that writes MainSchool, SchoolPeriods, Projects, and Level,
//     sets UpdatedAt to the current timestamp, and RETURNs the updated row.
//   - Scans the returned row, unmarshals JSON columns back into the domain.Education,
//     and returns the updated object.
//
// Returns:
// - (*domain.Education, nil) on success with the updated record.
// - (nil, nil) if no row with the given id was found.
// - (nil, error) on validation, marshaling, database, or unmarshaling errors.
func (r *educationRepository) Update(ctx context.Context, education *domain.Education) (*domain.Education, error) {
	if education == nil {
		return nil, errors.New("failed to validate education: payload is nil")
	}

	if education.Id == "" {
		return nil, fmt.Errorf("failed to update education: ID missing")
	}

	if err := education.ValidatePayload(); err != nil {
		return nil, fmt.Errorf("failed to validate education: %w", err)
	}

	now := r.timeProvider()
	education.UpdatedAt = now

	mainSchoolJSON, err := json.Marshal(education.MainSchool)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal main school: %w", err)
	}
	schoolPeriodsJSON, err := json.Marshal(education.SchoolPeriods)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal school periods: %w", err)
	}

	var (
		mainSchoolBytes    []byte
		schoolPeriodsBytes []byte
		updatedEducation   domain.Education
	)

	query := fmt.Sprintf(
		`UPDATE %s
		SET main_school=$2,
			school_periods=$3,
			level=$4,
			updated_at=$5
		WHERE id=$1
		RETURNING id, main_school, school_periods, level, created_at, updated_at`,
		r.educationTable,
	)

	err = r.databaseAPI.QueryRow(
		ctx,
		query,
		education.Id,
		mainSchoolJSON,
		schoolPeriodsJSON,
		education.Level,
		education.UpdatedAt,
	).Scan(
		&updatedEducation.Id,
		&mainSchoolBytes,
		&schoolPeriodsBytes,
		&updatedEducation.Level,
		&updatedEducation.CreatedAt,
		&updatedEducation.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update education: %w", err)
	}

	if err := json.Unmarshal(mainSchoolBytes, &updatedEducation.MainSchool); err != nil {
		return nil, fmt.Errorf("failed to unmarshal main school: %w", err)
	}
	if len(schoolPeriodsBytes) > 0 {
		if err = json.Unmarshal(schoolPeriodsBytes, &updatedEducation.SchoolPeriods); err != nil {
			return nil, fmt.Errorf("failed to unmarshal school periods: %w", err)
		}
	}

	if err := updatedEducation.ValidateResponse(); err != nil {
		return nil, fmt.Errorf("invalid education returned: %w", err)
	}

	return &updatedEducation, nil
}

// Delete removes the education record with the given id from the repository.
// It validates the id is not empty, executes a SQL DELETE against the repository's
// education table, and returns an error if the deletion fails.
//
// The function returns:
//   - an error wrapping the underlying database error when the Exec fails,
//   - pgx.ErrNoRows when no row was deleted (i.e., the id does not exist),
//   - an error if the provided id is empty.
//
// Parameters:
//
//	ctx - context for cancellations, timeouts and request-scoped values.
//	id  - unique identifier of the education record to delete.
func (r *educationRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("failed to delete education: ID missing")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", r.educationTable)

	cmdTag, err := r.databaseAPI.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete education: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *educationRepository) List(ctx context.Context, filter domain.EducationFilter) ([]domain.Education, error) {
	// Set defaults if not provided
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 || filter.PageSize > 20 {
		filter.PageSize = 20
	}
	if filter.SortBy == nil {
		defaultSort := domain.CreatedAt
		filter.SortBy = &defaultSort
	}

	baseQuery := fmt.Sprintf(
		`SELECT id, main_school, school_periods, level, created_at, updated_at FROM %s`,
		r.educationTable,
	)

	// Validate SortBy against allowed columns
	allowedSortColumns := map[domain.SortBy]bool{
		domain.CreatedAt: true,
		domain.UpdatedAt: true,
	}
	if !allowedSortColumns[*filter.SortBy] {
		return nil, fmt.Errorf("invalid sort column: %s", *filter.SortBy)
	}

	// Add sorting
	var sortColumn string
	switch *filter.SortBy {
	case domain.CreatedAt:
		sortColumn = "created_at"
	case domain.UpdatedAt:
		sortColumn = "updated_at"
	default:
		return nil, fmt.Errorf("invalid sort column: %s", *filter.SortBy)
	}

	sortOrder := "ASC"
	if !filter.SortAscending {
		sortOrder = "DESC"
	}
	baseQuery += fmt.Sprintf(" ORDER BY %s %s", sortColumn, sortOrder)

	// Add pagination
	offset := (filter.Page - 1) * filter.PageSize
	baseQuery += " LIMIT $1 OFFSET $2"

	// Execute query
	rows, err := r.databaseAPI.Query(ctx, baseQuery, filter.PageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list education: %w", err)
	}
	defer rows.Close()

	education := []domain.Education{}
	for rows.Next() {
		var (
			ed                domain.Education
			mainSchoolJSON    []byte
			schoolPeriodsJSON []byte
		)

		err := rows.Scan(
			&ed.Id,
			&mainSchoolJSON,
			&schoolPeriodsJSON,
			&ed.Level,
			&ed.CreatedAt,
			&ed.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan education: %w", err)
		}

		if err = json.Unmarshal(mainSchoolJSON, &ed.MainSchool); err != nil {
			return nil, fmt.Errorf("failed to unmarshal main school: %w", err)
		}
		if len(schoolPeriodsJSON) > 0 {
			if err = json.Unmarshal(schoolPeriodsJSON, &ed.SchoolPeriods); err != nil {
				return nil, fmt.Errorf("failed to unmarshal school periods: %w", err)
			}
		}

		education = append(education, ed)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return education, nil
}
