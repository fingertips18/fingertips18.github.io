package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
	v1 "github.com/fingertips18/fingertips18.github.io/backend/internal/repository/v1"
)

type ImageHandler interface {
	http.Handler
	Upload(w http.ResponseWriter, r *http.Request)
}

type ImageServiceConfig struct {
	UploadthingToken string

	imageRepo v1.ImageRepository
}

type imageServiceHandler struct {
	imageRepo v1.ImageRepository
}

// NewImageServiceHandler returns an ImageHandler configured from the provided cfg.
// If cfg.imageRepo is nil, a default v1.ImageRepository is created using
// cfg.UploadthingToken. The resulting ImageHandler is an *imageServiceHandler
// whose imageRepo field is set to the provided or constructed repository.
func NewImageServiceHandler(cfg ImageServiceConfig) ImageHandler {
	imageRepo := cfg.imageRepo
	if imageRepo == nil {
		imageRepo = v1.NewImageRepository(
			v1.ImageRepositoryConfig{
				UploadthingToken: cfg.UploadthingToken,
			},
		)
	}

	return &imageServiceHandler{
		imageRepo: imageRepo,
	}
}

// ServeHTTP handles HTTP requests for image operations.
// It routes requests based on the URL path after removing the "/image" prefix.
// Requests to "/image" are routed to the Upload handler.
// All other paths result in a 404 Not Found response.
func (h *imageServiceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/image")

	switch path {
	case "/upload":
		h.Upload(w, r)
	default:
		http.NotFound(w, r)
	}
}

// Upload handles HTTP POST requests to upload image files.
// It expects a JSON request body containing file metadata and upload configuration.
// The method validates the HTTP method, decodes the request body, converts DTOs to domain objects,
// and delegates the upload operation to the image repository.
// On success, it returns a 202 Accepted status with a JSON response containing the uploaded URL.
// On failure, it returns appropriate HTTP error status codes with error messages.
//
// @Security ApiKeyAuth
// @Summary Upload an image
// @Description Handles image upload with the supplied metadata and returns the Uploadthing URL of the stored image.
// @Tags image
// @Accept json
// @Produce json
// @Param imageUpload body UploadRequestDTO true "Image upload payload"
// @Success 202 {object} UploadResponseDTO "Confirmation message"
// @Failure 400 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /image/upload [post]
func (h *imageServiceHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed: only POST is supported", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var req UploadRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	var files []domain.Files
	for _, f := range req.Files {
		files = append(files, domain.Files{
			Name:     f.Name,
			Size:     f.Size,
			Type:     f.Type,
			CustomID: f.CustomID,
		})
	}
	upload := domain.UploadRequest{
		Files:              files,
		ACL:                req.ACL,
		Metadata:           req.Metadata,
		ContentDisposition: req.ContentDisposition,
	}

	url, err := h.imageRepo.Upload(r.Context(), &upload)
	if err != nil {
		http.Error(w, "Failed to upload image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp := UploadResponseDTO{
		URL: url,
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
