package v1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/database"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	dto "github.com/fingertips18/fingertips18.github.io/backend/internal/handler/v1/dto"
	v1 "github.com/fingertips18/fingertips18.github.io/backend/internal/repository/v1"
	"github.com/jackc/pgx/v5"
)

type FileHandler interface {
	http.Handler
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request, id string)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request, id string)
	DeleteByParent(w http.ResponseWriter, r *http.Request, parentTable, parentID string)
	ListByParent(w http.ResponseWriter, r *http.Request)
	Upload(w http.ResponseWriter, r *http.Request)
}

type FileServiceConfig struct {
	DatabaseAPI          database.DatabaseAPI
	UploadthingSecretKey string

	fileRepo v1.FileRepository
}

type fileServiceHandler struct {
	fileRepo v1.FileRepository
}

// NewFileServiceHandler creates and returns a new instance of FileHandler.
// It accepts a FileServiceConfig, which may include a custom file repository.
// If no repository is provided in the config, it initializes a default FileRepository
// using the provided DatabaseAPI and a default table name.
// Returns a FileHandler implementation.
func NewFileServiceHandler(cfg FileServiceConfig) FileHandler {
	fileRepo := cfg.fileRepo
	if fileRepo == nil {
		fileRepo = v1.NewFileRepository(
			v1.FileRepositoryConfig{
				DatabaseAPI:          cfg.DatabaseAPI,
				FileTable:            "File",
				UploadthingSecretKey: cfg.UploadthingSecretKey,
			},
		)
	}

	return &fileServiceHandler{
		fileRepo: fileRepo,
	}
}

// ServeHTTP handles HTTP requests for file-related endpoints.
//
// It supports the following routes:
//   - GET    /files?parent_table=...&parent_id=...&role=...  : List files by parent
//   - POST   /file                                           : Create a new file record
//   - GET    /file/{id}                                      : Retrieve a file by its ID
//   - DELETE /file/{id}                                      : Delete a file by its ID
//
// For unsupported methods or unknown routes, it responds with appropriate HTTP error codes.
func (h *fileServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")

	switch {
	case path == "/file/upload":
		h.Upload(w, r)
		return

	// GET /files?parent_table=...&parent_id=...&role=...
	case path == "/files":
		switch r.Method {
		case http.MethodGet:
			h.ListByParent(w, r)
		case http.MethodDelete:
			parentTable := r.URL.Query().Get("parent_table")
			parentID := r.URL.Query().Get("parent_id")
			h.DeleteByParent(w, r, parentTable, parentID)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return

	// POST / PUT /file
	case path == "/file":
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		case http.MethodPut:
			h.Update(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		return

	// GET / DELETE /file/{id}
	case strings.HasPrefix(path, "/file/"):
		id := strings.TrimPrefix(path, "/file/")

		if r.Method != http.MethodGet && r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if id == "" {
			http.Error(w, "File ID is required", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			h.Get(w, r, id)
		default: // Must be DELETE due to earlier check
			h.Delete(w, r, id)
		}
		return

	// Unknown route
	default:
		http.NotFound(w, r)
		return
	}
}

// Create handles HTTP POST requests to create a new file record.
// It expects a JSON payload in the request body representing file metadata.
// On success, it responds with a JSON object containing the new file's ID and a status message.
// If the request method is not POST, the JSON is invalid, or file creation fails, it responds with an appropriate HTTP error.
//
// @Security ApiKeyAuth
// @Summary Create a file record
// @Description Creates a new file metadata record from the provided JSON payload. Returns the created file ID.
// @Tags file
// @Accept json
// @Produce json
// @Param file body dto.CreateFileRequest true "File payload"
// @Success 201 {object} IDResponse "File ID"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /file [post]
func (h *fileServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	file := &domain.File{
		ParentTable: domain.ParentTable(req.ParentTable),
		ParentID:    req.ParentID,
		Role:        domain.FileRole(req.Role),
		Name:        req.Name,
		URL:         req.URL,
		Type:        req.Type,
		Size:        req.Size,
	}

	id, err := h.fileRepo.Create(r.Context(), *file)
	if err != nil {
		msg := err.Error()
		if len(msg) > 0 {
			msg = strings.ToUpper(msg[:1]) + msg[1:]
		}
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	resp := IDResponse{
		Id: id,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(buf.Bytes())
}

// Get handles HTTP GET requests to retrieve a file by its ID.
// It expects the file ID as a URL path parameter.
// On success, it responds with a JSON representation of the file.
// If the ID is invalid, the file is not found, or retrieval fails, it responds with an appropriate HTTP error.
//
// @Security ApiKeyAuth
// @Summary Get a file by ID
// @Description Retrieves a file record by its unique identifier.
// @Tags file
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} dto.FileDTO "File details"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /file/{id} [get]
func (h *fileServiceHandler) Get(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	file, err := h.fileRepo.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.FileDTO{
		ID:          file.ID,
		ParentTable: string(file.ParentTable),
		ParentID:    file.ParentID,
		Role:        string(file.Role),
		Name:        file.Name,
		URL:         file.URL,
		Type:        file.Type,
		Size:        file.Size,
		CreatedAt:   file.CreatedAt,
		UpdatedAt:   file.UpdatedAt,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// Update handles HTTP PUT requests to update an existing file record.
// It expects a JSON payload in the request body with the file ID and updated fields.
// On success, it responds with a JSON object containing the updated file details.
// If the JSON is invalid, the file ID is missing, or the update fails, it responds with an appropriate HTTP error.
//
// @Security ApiKeyAuth
// @Summary Update a file record
// @Description Updates an existing file metadata record. The file ID must be provided.
// @Tags file
// @Accept json
// @Produce json
// @Param file body dto.FileDTO true "Updated file data"
// @Success 200 {object} dto.FileDTO "Updated file details"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse "Not Found"
// @Failure 500 {object} ErrorResponse
// @Router /file [put]
func (h *fileServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.FileDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		http.Error(w, "Invalid: ID missing", http.StatusBadRequest)
		return
	}

	file := &domain.File{
		ID:          req.ID,
		ParentTable: domain.ParentTable(req.ParentTable),
		ParentID:    req.ParentID,
		Role:        domain.FileRole(req.Role),
		Name:        req.Name,
		URL:         req.URL,
		Type:        req.Type,
		Size:        req.Size,
	}

	updatedFile, err := h.fileRepo.Update(r.Context(), *file)
	if err != nil {
		msg := err.Error()
		if len(msg) > 0 {
			msg = strings.ToUpper(msg[:1]) + msg[1:]
		}
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if updatedFile == nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	resp := dto.FileDTO{
		ID:          updatedFile.ID,
		ParentTable: string(updatedFile.ParentTable),
		ParentID:    updatedFile.ParentID,
		Role:        string(updatedFile.Role),
		Name:        updatedFile.Name,
		URL:         updatedFile.URL,
		Type:        updatedFile.Type,
		Size:        updatedFile.Size,
		CreatedAt:   updatedFile.CreatedAt,
		UpdatedAt:   updatedFile.UpdatedAt,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// Delete handles HTTP DELETE requests to remove a file by its ID.
// It expects the file ID as a URL path parameter.
// On success, it responds with status "ok".
// If the ID is invalid, the file is not found, or deletion fails, it responds with an appropriate HTTP error.
//
// @Security ApiKeyAuth
// @Summary Delete a file by ID
// @Description Deletes a file record by its unique identifier.
// @Tags file
// @Produce json
// @Param id path string true "File ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /file/{id} [delete]
func (h *fileServiceHandler) Delete(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		http.Error(w, "Invalid file ID", http.StatusBadRequest)
		return
	}

	err := h.fileRepo.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Successful delete → 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// DeleteByParent handles HTTP DELETE requests to remove all files associated with a parent entity.
// It expects the parent table name and parent ID as parameters.
// On success, it responds with status 204 No Content.
// If required parameters are missing or the deletion fails, it responds with an appropriate HTTP error.
//
// @Security ApiKeyAuth
// @Summary Delete files by parent
// @Description Deletes all file records associated with a specific parent entity.
// @Tags file
// @Produce json
// @Param parent_table query string true "Parent table name"
// @Param parent_id query string true "Parent ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /files [delete]
func (h *fileServiceHandler) DeleteByParent(w http.ResponseWriter, r *http.Request, parentTable, parentID string) {
	if parentTable == "" {
		http.Error(w, "Invalid: missing parentTable", http.StatusBadRequest)
		return
	}

	if parentID == "" {
		http.Error(w, "Invalid: missing parentID", http.StatusBadRequest)
		return
	}

	err := h.fileRepo.DeleteByParent(r.Context(), parentTable, parentID)
	if err != nil {
		http.Error(w, "Failed to delete files on table "+parentTable+" with ID "+parentID+":"+err.Error(), http.StatusInternalServerError)
		return
	}

	// Successful delete → 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// ListByParent handles HTTP GET requests to retrieve files by parent entity and role.
// It expects query parameters: parentTable, parentId, and role.
// On success, it responds with a JSON array of files matching the criteria.
// If required parameters are missing or invalid, it responds with an appropriate HTTP error.
//
// @Security ApiKeyAuth
// @Summary List files by parent
// @Description Retrieves all files for a specific parent entity and role.
// @Tags file
// @Produce json
// @Param parent_table query string true "Parent table name"
// @Param parent_id query string true "Parent ID"
// @Param role query string true "File role (image)"
// @Success 200 {array} dto.FileDTO "List of files"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /files [get]
func (h *fileServiceHandler) ListByParent(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	parentTable := query.Get("parent_table")
	parentID := query.Get("parent_id")
	role := query.Get("role")

	if parentTable == "" {
		http.Error(w, "parent_table query parameter is required", http.StatusBadRequest)
		return
	}
	if parentID == "" {
		http.Error(w, "parent_id query parameter is required", http.StatusBadRequest)
		return
	}
	if role == "" {
		http.Error(w, "role query parameter is required", http.StatusBadRequest)
		return
	}

	files, err := h.fileRepo.FindByParent(r.Context(), parentTable, parentID, domain.FileRole(role))
	if err != nil {
		http.Error(w, "Failed to retrieve files: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fileResponses := make([]dto.FileDTO, 0, len(files))
	for _, file := range files {
		fileResponses = append(fileResponses, dto.FileDTO{
			ID:          file.ID,
			ParentTable: string(file.ParentTable),
			ParentID:    file.ParentID,
			Role:        string(file.Role),
			Name:        file.Name,
			URL:         file.URL,
			Type:        file.Type,
			Size:        file.Size,
			CreatedAt:   file.CreatedAt,
			UpdatedAt:   file.UpdatedAt,
		})
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(fileResponses); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(buf.Bytes())
}

// Upload handles HTTP POST requests to upload file files.
// It expects a JSON request body containing file metadata and upload configuration.
// The method validates the HTTP method, decodes the request body, converts DTOs to domain objects,
// and delegates the upload operation to the file repository.
// On success, it returns a 202 Accepted status with a JSON response containing the uploaded URL.
// On failure, it returns appropriate HTTP error status codes with error messages.
//
// @Security ApiKeyAuth
// @Summary Upload a file
// @Description Handles file upload with the supplied metadata and returns the Uploadthing URL of the stored file.
// @Tags file
// @Accept json
// @Produce json
// @Param fileUpload body dto.FileUploadRequestDTO true "File upload payload"
// @Success 202 {object} dto.FileUploadedResponseDTO "File upload URL"
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /file/upload [post]
func (h *fileServiceHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed: only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	var req dto.FileUploadRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	var files []domain.FileUpload
	for _, f := range req.Files {
		files = append(files, domain.FileUpload{
			Name:     f.Name,
			Size:     f.Size,
			Type:     f.Type,
			CustomID: f.CustomID,
		})
	}
	upload := domain.FileUploadRequest{
		Files:              files,
		ACL:                req.ACL,
		Metadata:           req.Metadata,
		ContentDisposition: req.ContentDisposition,
	}

	uploaded, err := h.fileRepo.Upload(r.Context(), &upload)
	if err != nil {
		// The error in the repo is comprehensive enough
		// Ensure that the first letter is capitalize
		msg := err.Error()
		if len(msg) > 0 {
			msg = strings.ToUpper(msg[:1]) + msg[1:]
		}

		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	file := dto.FileUploadedDTO{
		Key:                uploaded.Key,
		FileName:           uploaded.FileName,
		FileType:           uploaded.FileType,
		FileUrl:            uploaded.FileUrl,
		ContentDisposition: uploaded.ContentDisposition,
		PollingJwt:         uploaded.PollingJwt,
		PollingUrl:         uploaded.PollingUrl,
		CustomId:           uploaded.CustomId,
		URL:                uploaded.URL,
		Fields:             uploaded.Fields,
	}

	resp := dto.FileUploadedResponseDTO{
		File: file,
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		http.Error(w, "Failed to write response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(buf.Bytes())
}
