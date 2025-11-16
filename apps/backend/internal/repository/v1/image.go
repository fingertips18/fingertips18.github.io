package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/fingertips18/fingertips18.github.io/backend/internal/client"
	"github.com/fingertips18/fingertips18.github.io/backend/internal/domain"
)

type ImageRepository interface {
	Upload(ctx context.Context, image *domain.UploadRequest) (string, error)
}

type ImageRepositoryConfig struct {
	UploadthingToken string

	httpAPI client.HttpAPI
}

type imageRepository struct {
	uploadthingToken string
	httpAPI          client.HttpAPI
}

// NewImageRepository creates and returns an ImageRepository configured using the
// provided ImageRepositoryConfig. If cfg.httpAPI is nil, a default HTTP API
// client with a 30-second timeout will be created. The returned repository will
// use cfg.UploadthingToken for authenticated requests and the configured
// httpAPI for performing image-related operations.
func NewImageRepository(cfg ImageRepositoryConfig) ImageRepository {
	httpAPI := cfg.httpAPI
	if httpAPI == nil {
		httpAPI = client.NewHTTPAPI(30 * time.Second)
	}

	return &imageRepository{
		uploadthingToken: cfg.UploadthingToken,
		httpAPI:          httpAPI,
	}
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
//     repository's Uploadthing API key (r.uploadingthingToken). The HTTP request is executed
//     with the provided ctx.
//   - Treats any non-200 (OK) response as an error and includes the response status and body
//     in the returned error for diagnostics.
//   - Decodes the successful response into domain.UploadthingUploadResponse, validates it,
//     and returns the FileUrl of the first returned file (uploadResp.Data[0].FileUrl).
//   - All underlying errors are wrapped with context for easier debugging.
//
// Logging: the method logs the upload attempt and the resulting uploaded file URL on success.
//
// Return values:
// - string: the URL of the uploaded file on success, or an empty string on failure.
// - error: non-nil if validation, marshaling, network, decoding, or response validation fails.
func (r *imageRepository) Upload(ctx context.Context, image *domain.UploadRequest) (string, error) {
	// Validate request structure
	if err := image.Validate(); err != nil {
		return "", fmt.Errorf("failed to validate image: %w", err)
	}

	log.Println("Attempting to upload the image...")

	payload := *image

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
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Build HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.uploadthing.com/v6/uploadFiles",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Uploadthing-Api-Key", r.uploadthingToken)

	// Execute request
	resp, err := r.httpAPI.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 response
	if resp.StatusCode != http.StatusOK {
		respBody, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			log.Printf("failed to read error response body: %v", readErr)
		}
		return "", fmt.Errorf(
			"failed to upload image: status=%s message=%s",
			resp.Status,
			respBody,
		)
	}

	// Decode UploadThing response
	var uploadResp domain.UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return "", fmt.Errorf("failed to decode uploadthing response: %w", err)
	}

	// Make sure at least one file was returned
	if err := uploadResp.Validate(); err != nil {
		return "", fmt.Errorf("invalid uploadthing response: %w", err)
	}

	// Extract file URL
	fileUrl := uploadResp.Data[0].Data.URL

	log.Println("Image uploaded successfully:", fileUrl)

	return fileUrl, nil
}
