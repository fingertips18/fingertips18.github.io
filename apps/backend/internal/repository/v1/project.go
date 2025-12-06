package v1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/database"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/utils"
	"github.com/fingertips18/fingertips18.github.io/backend/pkg/metadata"
	"github.com/jackc/pgx/v5"
)

type ProjectRepository interface {
	Create(ctx context.Context, project *domain.Project) (string, error)
	Get(ctx context.Context, id string) (*domain.Project, error)
	Update(ctx context.Context, project *domain.Project) (*domain.Project, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter domain.ProjectFilter) ([]domain.Project, error)
	ListByEducationID(ctx context.Context, educationID string) ([]domain.Project, error)
	ListByEducationIDs(ctx context.Context, educationIDs []string) (map[string][]domain.Project, error)
}

type ProjectRepositoryConfig struct {
	DatabaseAPI  database.DatabaseAPI
	BlurHashAPI  metadata.BlurHashAPI
	ProjectTable string

	timeProvider domain.TimeProvider
}

type projectRepository struct {
	projectTable string
	databaseAPI  database.DatabaseAPI
	blurHashAPI  metadata.BlurHashAPI
	timeProvider domain.TimeProvider
}

// NewProjectRepository creates and returns a new instance of ProjectRepository.
// It initializes the repository with the provided configuration. If the pgxAPI or timeProvider
// are not supplied in the configuration, it sets them to default implementations.
//   - cfg: Configuration for the project repository, including database connection and table info.
//
// Returns a ProjectRepository implementation.
func NewProjectRepository(cfg ProjectRepositoryConfig) ProjectRepository {
	blurHashAPI := cfg.BlurHashAPI
	if blurHashAPI == nil {
		blurHashAPI = metadata.NewBlurHashAPI()
	}

	timeProvider := cfg.timeProvider
	if timeProvider == nil {
		timeProvider = time.Now
	}

	return &projectRepository{
		projectTable: cfg.ProjectTable,
		databaseAPI:  cfg.DatabaseAPI,
		blurHashAPI:  blurHashAPI,
		timeProvider: timeProvider,
	}
}

// Create validates and persists a new project, returning the newly generated project ID.
//
// The method performs the following steps:
// 1. Validates the provided project payload.
// 2. Generates a unique ID for the project (via utils.GenerateKey).
// 3. Sets CreatedAt and UpdatedAt on the project copy using the repository's timeProvider.
// 4. Inserts the project into the repository's configured project table using the provided context.
// 5. Scans and returns the inserted project's ID.
//
// The provided context is used for the database operation and may cancel or time out the request.
// Returns the new project ID on success. Returns an error if validation fails, the database query fails,
// or the database returns an empty/missing ID.
func (r *projectRepository) Create(ctx context.Context, project *domain.Project) (string, error) {
	if project == nil {
		return "", errors.New("failed to validate project: payload is nil")
	}

	if err := project.ValidatePayload(r.blurHashAPI); err != nil {
		return "", fmt.Errorf("failed to validate project: %w", err)
	}

	id := utils.GenerateKey()
	now := r.timeProvider()

	project.CreatedAt = now
	project.UpdatedAt = now

	query := fmt.Sprintf(
		`INSERT INTO %s
		(id, preview, blur_hash, title, sub_title, description, tags, type, link, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`,
		r.projectTable,
	)

	var returnedID string
	err := r.databaseAPI.QueryRow(
		ctx,
		query,
		id,
		project.Preview,
		project.BlurHash,
		project.Title,
		project.Subtitle,
		project.Description,
		project.Tags,
		project.Type,
		project.Link,
		project.CreatedAt,
		project.UpdatedAt,
	).Scan(&returnedID)

	if err != nil {
		return "", fmt.Errorf("failed to create project: %w", err)
	}

	if returnedID == "" {
		return "", errors.New("invalid project returned: ID missing")
	}

	return returnedID, nil
}

// Get retrieves a project by its unique identifier from the repository's database.
// The ctx is used for cancellation and deadlines. The id must be non-empty; if it is
// empty, Get returns an error indicating a missing ID. On success, it returns a
// pointer to a validated domain.Project populated from the matching database row.
// If no row matches the provided id, the returned error will wrap pgx.ErrNoRows.
// Any other query/scan or validation failures are returned wrapped to provide context.
func (r *projectRepository) Get(ctx context.Context, id string) (*domain.Project, error) {
	if id == "" {
		return nil, fmt.Errorf("failed to get project: ID missing")
	}

	var project domain.Project

	query := fmt.Sprintf(
		`SELECT id, preview, blur_hash, title, sub_title, description, tags, type, link, created_at, updated_at
		FROM %s
		WHERE id = $1`,
		r.projectTable,
	)

	err := r.databaseAPI.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&project.Id,
		&project.Preview,
		&project.BlurHash,
		&project.Title,
		&project.Subtitle,
		&project.Description,
		&project.Tags,
		&project.Type,
		&project.Link,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to get project: %w", err)
		}
		return nil, fmt.Errorf("failed to scan project: %w", err)
	}

	if err := project.ValidateResponse(r.blurHashAPI); err != nil {
		return nil, fmt.Errorf("invalid project returned: %w", err)
	}

	return &project, nil
}

// Update updates the project identified by project.Id in the repository.
// It validates the provided project payload, updates the UpdatedAt timestamp
// from the repository's time provider, and persists the project's fields to
// the database. The method returns the updated project as stored in the
// database. If no row matches the provided id, it returns (nil, nil).
// Validation errors or other database errors are returned (wrapped) to the
// caller. The provided context is used for database cancellation and timeouts.
func (r *projectRepository) Update(ctx context.Context, project *domain.Project) (*domain.Project, error) {
	if project == nil {
		return nil, errors.New("failed to validate project: payload is nil")
	}
	if project.Id == "" {
		return nil, fmt.Errorf("failed to update project: ID missing")
	}

	if err := project.ValidatePayload(r.blurHashAPI); err != nil {
		return nil, fmt.Errorf("failed to validate project: %w", err)
	}

	now := r.timeProvider()
	project.UpdatedAt = now

	var updatedProject domain.Project

	query := fmt.Sprintf(
		`UPDATE %s
		SET preview=$2,
			blur_hash=$3,
			title=$4,
			sub_title=$5,
			description=$6,
			tags=$7,
			type=$8,
			link=$9,
			updated_at=$10
		WHERE id=$1
		RETURNING id, preview, blur_hash, title, sub_title, description, tags, type, link, created_at, updated_at`,
		r.projectTable,
	)

	err := r.databaseAPI.QueryRow(
		ctx,
		query,
		project.Id,
		project.Preview,
		project.BlurHash,
		project.Title,
		project.Subtitle,
		project.Description,
		project.Tags,
		project.Type,
		project.Link,
		project.CreatedAt,
		project.UpdatedAt,
	).Scan(
		&updatedProject.Id,
		&updatedProject.Preview,
		&updatedProject.BlurHash,
		&updatedProject.Title,
		&updatedProject.Subtitle,
		&updatedProject.Description,
		&updatedProject.Tags,
		&updatedProject.Type,
		&updatedProject.Link,
		&updatedProject.CreatedAt,
		&updatedProject.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	if err := updatedProject.ValidateResponse(r.blurHashAPI); err != nil {
		return nil, fmt.Errorf("invalid project returned: %w", err)
	}

	return &updatedProject, nil
}

// Delete removes a project from the database by its ID.
// It returns an error if the deletion fails or if no project with the given ID exists.
//
// Parameters:
//   - ctx: The context for controlling cancellation and deadlines.
//   - id: The unique identifier of the project to delete.
//
// Returns:
//   - error: An error if the operation fails or if no project is found with the specified ID.
func (r *projectRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("failed to delete project: ID missing")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", r.projectTable)

	cmdTag, err := r.databaseAPI.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// List retrieves a paginated, optionally filtered and sorted slice of domain.Project from the repository.
// It takes a context for cancellation and a ProjectFilter that controls filtering, pagination and sorting.
// Defaults are applied when values are not provided: Page defaults to 1, PageSize defaults to 20 (and is capped at 20),
// and SortBy defaults to CreatedAt. If filter.Type is non-nil, results are restricted to that project type.
// Sorting is applied by the specified field in ascending order by default; set SortAscending to false for descending.
// Results are limited to PageSize with an offset of (Page-1)*PageSize.
// The query is executed with parameterized arguments to avoid SQL injection and the returned slice contains
// mapped domain.Project values. An error is returned if query execution, row scanning, or row iteration fails.
func (r *projectRepository) List(ctx context.Context, filter domain.ProjectFilter) ([]domain.Project, error) {
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
		`SELECT id, preview, blur_hash, title, sub_title, description, tags, type, link, created_at, updated_at FROM %s`,
		r.projectTable,
	)
	var conditions []string
	var args []any
	argIdx := 1

	// Add optional type filter
	if filter.Type != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIdx))
		args = append(args, *filter.Type)
		argIdx++
	}

	// Append WHERE clause if any filters exist
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add sorting
	sortOrder := "ASC"
	if !filter.SortAscending {
		sortOrder = "DESC"
	}
	baseQuery += fmt.Sprintf(" ORDER BY %s %s", *filter.SortBy, sortOrder)

	// Add pagination
	offset := (filter.Page - 1) * filter.PageSize
	baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.PageSize, offset)

	// Execute query
	rows, err := r.databaseAPI.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	projects := []domain.Project{}
	for rows.Next() {
		var project domain.Project

		err := rows.Scan(
			&project.Id,
			&project.Preview,
			&project.BlurHash,
			&project.Title,
			&project.Subtitle,
			&project.Description,
			&project.Tags,
			&project.Type,
			&project.Link,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return projects, nil
}

// ListProjectsByEducationID queries the database for all projects that belong to the
// provided educationID and returns them as a slice of domain.Project.
//
// Behavior:
//   - Executes a SELECT for id, preview, blur_hash, title, sub_title, description,
//     tags, type, link, education_id, created_at, updated_at from the project table
//     where education_id = $1.
//   - Results are ordered by created_at DESC.
//
// Parameters:
// - ctx: context for tracing/cancellation/timeout of the database operation.
// - educationID: the education identifier used as the single query parameter.
//
// Returns:
//   - ([]domain.Project, error): the slice of projects (possibly empty) and an error.
//     Errors are returned if the query fails, if scanning a row fails, or if rows
//     iteration reports an error. Errors are wrapped for context.
//
// Notes:
//   - The nullable education_id column is mapped into Project.EducationID only when
//     the SQL value is valid; otherwise the field is left zero-valued.
//   - The function defers closing rows and returns any rows.Err() after iteration.
func (r *projectRepository) ListByEducationID(ctx context.Context, educationID string) ([]domain.Project, error) {
	query := fmt.Sprintf(`
        SELECT id, preview, blur_hash, title, sub_title, description, tags, type, link, education_id, created_at, updated_at
        FROM %s
        WHERE education_id = $1
        ORDER BY created_at DESC
    `, r.projectTable)

	rows, err := r.databaseAPI.Query(ctx, query, educationID)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects by education_id: %w", err)
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var p domain.Project
		var educationID sql.NullString

		err := rows.Scan(
			&p.Id, &p.Preview, &p.BlurHash, &p.Title, &p.Subtitle,
			&p.Description, &p.Tags, &p.Type, &p.Link,
			&educationID, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}

		if educationID.Valid {
			p.EducationID = educationID.String
		}

		projects = append(projects, p)
	}

	return projects, rows.Err()
}

// ListByEducationIDs fetches projects for multiple education IDs in a single query.
// Returns a map where the key is the education ID and the value is a slice of projects.
func (r *projectRepository) ListByEducationIDs(ctx context.Context, educationIDs []string) (map[string][]domain.Project, error) {
	if len(educationIDs) == 0 {
		return make(map[string][]domain.Project), nil
	}

	placeholders := make([]string, len(educationIDs))
	args := make([]any, len(educationIDs))
	for i, id := range educationIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`
		SELECT id, preview, blur_hash, title, sub_title, description, tags, type, link, education_id, created_at, updated_at
		FROM %s
		WHERE education_id IN (%s)
		ORDER BY education_id, created_at DESC
	`, r.projectTable, strings.Join(placeholders, ", "))

	rows, err := r.databaseAPI.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to batch query projects: %w", err)
	}
	defer rows.Close()

	projectsByEducation := make(map[string][]domain.Project)
	for _, eduID := range educationIDs {
		projectsByEducation[eduID] = []domain.Project{}
	}

	for rows.Next() {
		var p domain.Project
		var educationID sql.NullString

		err := rows.Scan(
			&p.Id, &p.Preview, &p.BlurHash, &p.Title, &p.Subtitle,
			&p.Description, &p.Tags, &p.Type, &p.Link,
			&educationID, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}

		if educationID.Valid {
			if _, ok := projectsByEducation[educationID.String]; ok {
				p.EducationID = educationID.String
				projectsByEducation[educationID.String] = append(projectsByEducation[educationID.String], p)
			}
		}
	}

	return projectsByEducation, rows.Err()
}
