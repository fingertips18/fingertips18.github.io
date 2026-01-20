package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/client"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/database"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/utils"
	"github.com/jackc/pgx/v5"
)

type FileRepository interface {
	FindByParent(ctx context.Context, parentTable, parentID string, role domain.FileRole) ([]domain.File, error)
	Create(ctx context.Context, file domain.File) (string, error)
	Update(ctx context.Context, fileUpdate domain.File) (*domain.File, error)
	Delete(ctx context.Context, id string) error
	DeleteByParent(ctx context.Context, parentTable string, parentID string) error
	FindByID(ctx context.Context, id string) (*domain.File, error)
	Upload(ctx context.Context, file *domain.FileUploadRequest) (*domain.FileUploaded, error)
}

type FileRepositoryConfig struct {
	DatabaseAPI          database.DatabaseAPI
	FileTable            string
	UploadthingSecretKey string

	httpAPI      client.HttpAPI
	timeProvider domain.TimeProvider
}

type fileRepository struct {
	fileTable            string
	databaseAPI          database.DatabaseAPI
	uploadthingSecretKey string
	httpAPI              client.HttpAPI
	timeProvider         domain.TimeProvider
}

// NewFileRepository creates and returns a configured FileRepository.
//
// It accepts a FileRepositoryConfig and constructs an internal
// fileRepository backed by cfg.FileTable and cfg.DatabaseAPI.
// If cfg.timeProvider is nil, the repository defaults to using time.Now
// as the time provider. The returned value implements the
// FileRepository interface and is never nil.
func NewFileRepository(cfg FileRepositoryConfig) FileRepository {
	httpAPI := cfg.httpAPI
	if httpAPI == nil {
		httpAPI = client.NewHTTPAPI(30 * time.Second)
	}

	timeProvider := cfg.timeProvider
	if timeProvider == nil {
		timeProvider = time.Now
	}

	return &fileRepository{
		fileTable:            cfg.FileTable,
		databaseAPI:          cfg.DatabaseAPI,
		uploadthingSecretKey: cfg.UploadthingSecretKey,
		httpAPI:              httpAPI,
		timeProvider:         timeProvider,
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
// Note: Since the file is passed by value, the caller's struct is not modified with timestamps.
func (r *fileRepository) Create(ctx context.Context, file domain.File) (string, error) {
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

// Update updates an existing file record in the repository.
// It validates the provided File payload, ensures the ID is present, sets the UpdatedAt
// timestamp from the repository's time provider, and updates the record in the configured
// file table. The method returns the updated file record with all fields populated from
// the database.
//
// Parameters:
//   - ctx: The context for controlling cancellation and deadlines.
//   - fileUpdate: The file record with updates. Must have a non-empty ID and pass validation.
//
// Returns:
//   - *domain.File: A pointer to the updated file record if successful, or nil if not found.
//   - error: An error if validation fails, the database update fails, or the returned file is invalid.
func (r *fileRepository) Update(ctx context.Context, fileUpdate domain.File) (*domain.File, error) {
	if fileUpdate.ID == "" {
		return nil, fmt.Errorf("failed to update file: ID missing")
	}

	if err := fileUpdate.ValidatePayload(); err != nil {
		return nil, fmt.Errorf("failed to validate file: %w", err)
	}

	now := r.timeProvider()
	fileUpdate.UpdatedAt = now

	var updatedFile domain.File

	query := fmt.Sprintf(
		`UPDATE %s
		SET parent_table=$2,
			parent_id=$3,
			role=$4,
			name=$5,
			url=$6,
			type=$7,
			size=$8,
			updated_at=$9
		WHERE id=$1
		RETURNING id, parent_table, parent_id, role, name, url, type, size, created_at, updated_at`,
		r.fileTable,
	)

	err := r.databaseAPI.QueryRow(
		ctx,
		query,
		fileUpdate.ID,
		fileUpdate.ParentTable,
		fileUpdate.ParentID,
		fileUpdate.Role,
		fileUpdate.Name,
		fileUpdate.URL,
		fileUpdate.Type,
		fileUpdate.Size,
		fileUpdate.UpdatedAt,
	).Scan(
		&updatedFile.ID,
		&updatedFile.ParentTable,
		&updatedFile.ParentID,
		&updatedFile.Role,
		&updatedFile.Name,
		&updatedFile.URL,
		&updatedFile.Type,
		&updatedFile.Size,
		&updatedFile.CreatedAt,
		&updatedFile.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update file: %w", err)
	}

	if err := updatedFile.ValidateResponse(); err != nil {
		return nil, fmt.Errorf("invalid file returned: %w", err)
	}

	return &updatedFile, nil
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

// DeleteByParent removes all files associated with a specific parent entity from the database.
// It validates that both parentTable and parentID are provided, then deletes all matching records.
//
// Parameters:
//   - ctx: The context for controlling cancellation and deadlines.
//   - parentTable: The name of the parent table (e.g., "projects", "education").
//   - parentID: The unique identifier of the parent record.
//
// Returns:
//   - error: An error if validation fails or the database deletion fails.
func (r *fileRepository) DeleteByParent(ctx context.Context, parentTable string, parentID string) error {
	if parentTable == "" {
		return errors.New("failed to delete files: parentTable missing")
	}

	if parentID == "" {
		return errors.New("failed to delete files: parentID missing")
	}

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE parent_table = $1 AND parent_id = $2",
		r.fileTable,
	)

	_, err := r.databaseAPI.Exec(ctx, query, parentTable, parentID)
	if err != nil {
		return fmt.Errorf("failed to delete files by parent: %w", err)
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

// Upload uploads the provided UploadthingUploadRequest to the UploadThing service and
// returns the URL of the uploaded file or an error.
//
// Behavior:
//   - Validates the incoming request via image.Validate() and returns an error if invalid.
//   - Applies default fallback values when not provided (ACL => "public-read",
//     ContentDisposition => "inline").
//   - Marshals the request to JSON and issues an HTTP POST to
//     "https://api.uploadthing.com/v6/uploadFiles" using the repository's httpAPI and the
//     repository's Uploadthing API key (r.uploadthingSecretKey). The HTTP request is executed
//     with the provided ctx.
//   - Treats any non-200 (OK) response as an error and includes the response status and body
//     in the returned error for diagnostics.
//   - Decodes the successful response into domain.FileUploadedResponse, validates it,
//     and returns the file metadata of the first returned file (uploadResp.Data[0]).
//   - All underlying errors are wrapped with context for easier debugging.
//
// Logging: the method logs the upload attempt and the resulting uploaded file URL on success.
//
// Return values:
// - *domain.FileUploadedResponse: the file metadata on success, or nil on failure.
// - error: non-nil if validation, marshaling, network, decoding, or response validation fails.
func (r *fileRepository) Upload(ctx context.Context, file *domain.FileUploadRequest) (*domain.FileUploaded, error) {
	// Validate request structure
	if err := file.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate file: %w", err)
	}

	log.Println("Attempting to upload the file...")

	payload := *file

	// Default fallback values
	if payload.ACL == nil {
		acl := "public-read"
		payload.ACL = &acl
	}
	if payload.ContentDisposition == nil {
		contentDisposition := "inline"
		payload.ContentDisposition = &contentDisposition
	}

	// Marshal request payload
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Build HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.uploadthing.com/v6/uploadFiles",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Uploadthing-Api-Key", r.uploadthingSecretKey)

	// Execute request
	resp, err := r.httpAPI.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 response
	if resp.StatusCode != http.StatusOK {
		respBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			log.Printf("failed to read error response body: %v", readErr)
		}
		return nil, fmt.Errorf(
			"failed to upload file: status=%s message=%s",
			resp.Status,
			respBody,
		)
	}

	// Decode UploadThing success response
	var uploadResp domain.FileUploadedResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, fmt.Errorf("failed to decode uploadthing response: %w", err)
	}

	// Make sure at least one file was returned
	if err := uploadResp.Validate(); err != nil {
		return nil, fmt.Errorf("invalid uploadthing response: %w", err)
	}

	// Extract file URL
	data := uploadResp.Data[0]

	log.Println("File uploaded successfully:", data.FileName)

	return &data, nil
}
