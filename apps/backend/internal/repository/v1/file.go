package v1

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/database"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/utils"
	"github.com/jackc/pgx/v5"
)

type FileRepository interface {
	FindByParent(ctx context.Context, parentTable, parentID string, role domain.FileRole) ([]domain.File, error)
	Create(ctx context.Context, file *domain.File) (string, error)
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*domain.File, error)
}

type FileRepositoryConfig struct {
	DatabaseAPI database.DatabaseAPI
	FileTable   string

	timeProvider domain.TimeProvider
}

type fileRepository struct {
	fileTable    string
	databaseAPI  database.DatabaseAPI
	timeProvider domain.TimeProvider
}

// NewFileRepository creates and returns a configured FileRepository.
//
// It accepts a FileRepositoryConfig and constructs an internal
// fileRepository backed by cfg.FileTable and cfg.DatabaseAPI.
// If cfg.timeProvider is nil, the repository defaults to using time.Now
// as the time provider. The returned value implements the
// FileRepository interface and is never nil.
func NewFileRepository(cfg FileRepositoryConfig) FileRepository {
	timeProvider := cfg.timeProvider
	if timeProvider == nil {
		timeProvider = time.Now
	}

	return &fileRepository{
		fileTable:    cfg.FileTable,
		databaseAPI:  cfg.DatabaseAPI,
		timeProvider: timeProvider,
	}
}

// FindByParent retrieves all files for a specific parent entity and role.
// It queries the database for files matching the provided parentTable, parentID, and role.
// The results are ordered by created_at in descending order.
//
// Parameters:
//   - ctx: The context for controlling cancellation and deadlines.
//   - parentTable: The name of the parent table (e.g., "projects", "education").
//   - parentID: The unique identifier of the parent record.
//   - role: The file role to filter by (e.g., FileRoleImage, FileRoleDocument).
//
// Returns:
//   - []domain.File: A slice of files matching the criteria (may be empty).
//   - error: An error if the query fails, scanning fails, or row iteration encounters an issue.
func (r *fileRepository) FindByParent(ctx context.Context, parentTable, parentID string, role domain.FileRole) ([]domain.File, error) {
	if parentTable == "" {
		return nil, errors.New("failed to find files: parentTable missing")
	}
	if parentID == "" {
		return nil, errors.New("failed to find files: parentID missing")
	}
	if role == "" {
		return nil, errors.New("failed to find files: role missing")
	}

	query := fmt.Sprintf(
		`SELECT id, parent_table, parent_id, role, name, url, type, size, created_at, updated_at
        FROM %s
        WHERE parent_table = $1 AND parent_id = $2 AND role = $3
        ORDER BY created_at DESC`,
		r.fileTable,
	)

	rows, err := r.databaseAPI.Query(ctx, query, parentTable, parentID, role)
	if err != nil {
		return nil, fmt.Errorf("failed to query files by parent: %w", err)
	}
	defer rows.Close()

	var files []domain.File
	for rows.Next() {
		var file domain.File
		err := rows.Scan(
			&file.ID,
			&file.ParentTable,
			&file.ParentID,
			&file.Role,
			&file.Name,
			&file.URL,
			&file.Type,
			&file.Size,
			&file.CreatedAt,
			&file.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file: %w", err)
		}

		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return files, nil
}

// Create creates a new file record in the repository and returns its generated ID.
// It validates the provided File payload, generates a unique ID, sets CreatedAt and
// UpdatedAt timestamps from the repository's time provider, and inserts the record
// into the configured file table.
//
// Parameters:
//   - ctx: The context for controlling cancellation and deadlines.
//   - file: The file record to create. Must not be nil and must pass validation.
//
// Returns:
//   - string: The newly created file's ID.
//   - error: An error if validation fails, database insertion fails, or the returned ID is empty.
//
// The provided file object's CreatedAt and UpdatedAt fields are updated when the method succeeds.
func (r *fileRepository) Create(ctx context.Context, file *domain.File) (string, error) {
	if file == nil {
		return "", errors.New("failed to validate file: payload is nil")
	}

	if err := file.ValidatePayload(); err != nil {
		return "", fmt.Errorf("failed to validate file: %w", err)
	}

	id := utils.GenerateKey()
	now := r.timeProvider()

	file.CreatedAt = now
	file.UpdatedAt = now

	query := fmt.Sprintf(
		`INSERT INTO %s
        (id, parent_table, parent_id, role, name, url, type, size, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING id`,
		r.fileTable,
	)

	var returnedID string
	err := r.databaseAPI.QueryRow(
		ctx,
		query,
		id,
		file.ParentTable,
		file.ParentID,
		file.Role,
		file.Name,
		file.URL,
		file.Type,
		file.Size,
		file.CreatedAt,
		file.UpdatedAt,
	).Scan(&returnedID)

	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	if returnedID == "" {
		return "", errors.New("invalid file returned: ID missing")
	}

	return returnedID, nil
}

// Delete removes a file from the database by its ID.
// It returns an error if the deletion fails or if no file with the given ID exists.
//
// Parameters:
//   - ctx: The context for controlling cancellation and deadlines.
//   - id: The unique identifier of the file to delete.
//
// Returns:
//   - error: An error if the operation fails or if no file is found with the specified ID (pgx.ErrNoRows).
func (r *fileRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("failed to delete file: ID missing")
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", r.fileTable)

	cmdTag, err := r.databaseAPI.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

// FindByID retrieves a file by its unique identifier from the repository's database.
// The ctx is used for cancellation and deadlines. The id must be non-empty; if it is
// empty, FindByID returns an error indicating a missing ID. On success, it returns a
// pointer to a validated domain.File populated from the matching database row.
// If no row matches the provided id, the returned error will wrap pgx.ErrNoRows.
// Any other query/scan or validation failures are returned wrapped to provide context.
//
// Parameters:
//   - ctx: The context for controlling cancellation and deadlines.
//   - id: The unique identifier of the file to retrieve.
//
// Returns:
//   - *domain.File: A pointer to the file record if found.
//   - error: An error if the file is not found, the query fails, or validation fails.
func (r *fileRepository) FindByID(ctx context.Context, id string) (*domain.File, error) {
	if id == "" {
		return nil, errors.New("failed to get file: ID missing")
	}

	var file domain.File

	query := fmt.Sprintf(
		`SELECT id, parent_table, parent_id, role, name, url, type, size, created_at, updated_at
        FROM %s
        WHERE id = $1`,
		r.fileTable,
	)

	err := r.databaseAPI.QueryRow(ctx, query, id).Scan(
		&file.ID,
		&file.ParentTable,
		&file.ParentID,
		&file.Role,
		&file.Name,
		&file.URL,
		&file.Type,
		&file.Size,
		&file.CreatedAt,
		&file.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("failed to get file: %w", err)
		}
		return nil, fmt.Errorf("failed to scan file: %w", err)
	}

	if err := file.ValidateResponse(); err != nil {
		return nil, fmt.Errorf("invalid file returned: %w", err)
	}

	return &file, nil
}
