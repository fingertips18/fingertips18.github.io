package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/database"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/utils"
	"github.com/jackc/pgx/v5"
)

type ProjectRepository interface {
	Create(ctx context.Context, createConfig domain.Project) (string, error)
	Get(ctx context.Context, id string) (*domain.Project, error)
	Update(ctx context.Context, project domain.Project) (*domain.Project, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter domain.ProjectFilter) ([]domain.Project, error)
}

type ProjectRepositoryConfig struct {
	Database     database.Database
	ProjectTable string

	timeProvider domain.TimeProvider
}

type projectRepository struct {
	projectTable string
	database     database.Database
	timeProvider domain.TimeProvider
}

// NewProjectRepository creates and returns a new instance of ProjectRepository.
// It initializes the repository with the provided configuration. If the pgxAPI or timeProvider
// are not supplied in the configuration, it sets them to default implementations.
//   - cfg: Configuration for the project repository, including database connection and table info.
//
// Returns a ProjectRepository implementation.
func NewProjectRepository(cfg ProjectRepositoryConfig) ProjectRepository {
	timeProvider := cfg.timeProvider
	if timeProvider == nil {
		timeProvider = time.Now
	}

	return &projectRepository{
		projectTable: cfg.ProjectTable,
		database:     cfg.Database,
		timeProvider: timeProvider,
	}
}

// Create inserts a new project record into the database and returns the generated project ID.
// It marshals the project's stack to JSON, generates a unique ID, and sets the creation and update timestamps.
// Returns the generated project ID on success, or an error if the operation fails.
func (r *projectRepository) Create(ctx context.Context, project domain.Project) (string, error) {
	stackJSON, err := json.Marshal(project.Stack)
	if err != nil {
		return "", fmt.Errorf("failed to marshal project stack: %w", err)
	}

	id := utils.GenerateKey()

	now := r.timeProvider()

	project.CreatedAt = now
	project.UpdatedAt = now

	query := fmt.Sprintf(
		`INSERT INTO %s
		(id, preview, blur_hash, title, sub_title, description, stack, type, link, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`,
		r.projectTable,
	)

	var returnedID string
	err = r.database.Pool.QueryRow(
		ctx,
		query,
		id,
		project.Preview,
		project.BlurHash,
		project.Title,
		project.SubTitle,
		project.Description,
		stackJSON,
		project.Type,
		project.Link,
		project.CreatedAt,
		project.UpdatedAt,
	).Scan(&returnedID)

	if err != nil {
		return "", fmt.Errorf("failed to create project: %w", err)
	}

	return returnedID, nil
}

// Get retrieves a project by its ID from the database.
// It returns the corresponding domain.Project object if found, or an error otherwise.
// The stack field is expected to be stored as JSON in the database and is unmarshaled into the Project struct.
// Parameters:
//   - ctx: context for controlling cancellation and deadlines.
//   - id: the unique identifier of the project to retrieve.
//
// Returns:
//   - *domain.Project: the project object if found.
//   - error: an error if the project could not be retrieved or unmarshaled.
func (r *projectRepository) Get(ctx context.Context, id string) (*domain.Project, error) {
	var project domain.Project
	stackJSON := []byte{}

	query := fmt.Sprintf(
		`SELECT id, preview, blur_hash, title, sub_title, description, stack, type, link, created_at, updated_at
		FROM %s
		WHERE id = $1`,
		r.projectTable,
	)

	err := r.database.Pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&project.Id,
		&project.Preview,
		&project.BlurHash,
		&project.Title,
		&project.SubTitle,
		&project.Description,
		&stackJSON,
		&project.Type,
		&project.Link,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // not found
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Decode stack JSON
	if err = json.Unmarshal(stackJSON, &project.Stack); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stack: %w", err)
	}

	return &project, nil
}

// Update updates an existing project in the database with the provided project data.
// It marshals the project's stack to JSON, sets the updated timestamp, and executes
// an UPDATE SQL statement. The updated project is returned with its fields populated
// from the database. Returns an error if marshalling, database update, or unmarshalling fails.
//
// Parameters:
//   - ctx: context.Context for controlling cancellation and deadlines.
//   - project: domain.Project containing the updated project data.
//
// Returns:
//   - *domain.Project: pointer to the updated project.
//   - error: error if the update fails or data cannot be marshaled/unmarshaled.
func (r *projectRepository) Update(ctx context.Context, project domain.Project) (*domain.Project, error) {
	stackJSON, err := json.Marshal(project.Stack)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal project stack: %w", err)
	}

	now := r.timeProvider()

	project.UpdatedAt = now

	var updatedProject domain.Project
	updatedStackJSON := []byte{}

	query := fmt.Sprintf(
		`UPDATE %s
		SET preview=$2,
			blur_hash=$3,
			title=$4,
			sub_title=$5,
			description=$6,
			stack=$7,
			type=$8,
			link=$9,
			created_at=$10,
			updated_at=$11
		WHERE id=$1
		RETURNING id, preview, blur_hash, title, sub_title, description, stack, type, link, created_at, updated_at`,
		r.projectTable,
	)

	err = r.database.Pool.QueryRow(
		ctx,
		query,
		project.Id,
		project.Preview,
		project.BlurHash,
		project.Title,
		project.SubTitle,
		project.Description,
		stackJSON,
		project.Type,
		project.Link,
		project.CreatedAt,
		project.UpdatedAt,
	).Scan(
		&updatedProject.Id,
		&updatedProject.Preview,
		&updatedProject.BlurHash,
		&updatedProject.Title,
		&updatedProject.SubTitle,
		&updatedProject.Description,
		&updatedStackJSON,
		&updatedProject.Type,
		&updatedProject.Link,
		&updatedProject.CreatedAt,
		&updatedProject.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // not found
		}
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	// Decode stack JSON
	if err = json.Unmarshal(updatedStackJSON, &updatedProject.Stack); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stack: %w", err)
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

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", r.projectTable)

	cmdTag, err := r.database.Pool.Exec(
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

// List retrieves a list of projects from the database based on the provided filter criteria.
// It supports filtering by project type, sorting by a specified field and order, and paginating results.
// If filter values for page, page size, or sort field are not provided, sensible defaults are applied.
// The function returns a slice of domain.Project and an error if the query or data processing fails.
//
// Parameters:
//   - ctx: The context for controlling cancellation and deadlines.
//   - filter: The ProjectFilter containing pagination, sorting, and optional type filter.
//
// Returns:
//   - []domain.Project: A slice of projects matching the filter criteria.
//   - error: An error if the query fails or data cannot be processed.
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
		`SELECT id, preview, blur_hash, title, sub_title, description, stack, type, link, created_at, updated_at FROM %s`,
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
	rows, err := r.database.Pool.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects: %w", err)
	}
	defer rows.Close()

	var projects []domain.Project
	for rows.Next() {
		var project domain.Project
		var stackJSON []byte

		err := rows.Scan(
			&project.Id,
			&project.Preview,
			&project.BlurHash,
			&project.Title,
			&project.SubTitle,
			&project.Description,
			&stackJSON,
			&project.Type,
			&project.Link,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}

		if err = json.Unmarshal(stackJSON, &project.Stack); err != nil {
			return nil, fmt.Errorf("failed to unmarshal project stack: %w", err)
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return projects, nil
}
