package domain

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ParentTable string

const (
	ProjectTable   ParentTable = "projects"
	UserTable      ParentTable = "users"
	EducationTable ParentTable = "educations"
	// Add other valid parent table names as needed
)

type FileRole string

const (
	Image FileRole = "image"
	// Add other valid file role names as needed
)

var mimeTypeRe = regexp.MustCompile(`^[a-z]+/[a-z0-9][a-z0-9\-\+\.]*$`)

// File represents a file attachment that can be associated with any parent entity.
// It uses a polymorphic association pattern via ParentTable and ParentID fields.
type File struct {
	ID          string      `json:"id"`
	ParentTable ParentTable `json:"parent_table"`
	ParentID    string      `json:"parent_id"`
	Role        FileRole    `json:"role"`
	Name        string      `json:"name"`
	URL         string      `json:"url"`
	Type        string      `json:"type"`
	Size        int64       `json:"size"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// isValid validates that the ParentTable is one of the allowed parent table types.
func (pt ParentTable) isValid() error {
	if strings.TrimSpace(string(pt)) == "" {
		return errors.New("parent_table missing")
	}

	switch pt {
	case ProjectTable, UserTable, EducationTable:
		return nil
	default:
		return errors.New("parent_table invalid")
	}
}

// isValid validates that the FileRole is one of the allowed file role types.
func (fr FileRole) isValid() error {
	if strings.TrimSpace(string(fr)) == "" {
		return errors.New("role missing")
	}

	switch fr {
	case Image:
		return nil
	default:
		return errors.New("role invalid")
	}
}

// isValidMimeType validates that the given MIME type string is in a valid format.
func isValidMimeType(mimeType string) error {
	if strings.TrimSpace(mimeType) == "" {
		return errors.New("type missing")
	}

	// Validate only the media-type part
	base := strings.TrimSpace(strings.SplitN(mimeType, ";", 2)[0])
	if mimeTypeRe.MatchString(strings.ToLower(base)) {
		return nil
	}
	return errors.New("invalid type")
}

// isValidUUID validates that the given id is a valid UUID format.
func isValidUUID(id string, label string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%s missing", label)
	}

	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%s invalid", label)
	}

	return nil
}

// isValidURL validates that the given value is a valid HTTP or HTTPS URL.
func isValidURL(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return errors.New("url missing")
	}

	u, err := url.ParseRequestURI(value)
	if err != nil {
		return errors.New("url invalid")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("url invalid")
	}
	if u.Host == "" {
		return errors.New("url invalid")
	}

	return nil
}

// ValidatePayload validates the File payload fields required for creating or updating a file.
func (f File) ValidatePayload() error {
	if err := f.ParentTable.isValid(); err != nil {
		return err
	}
	if err := isValidUUID(f.ParentID, "parent_id"); err != nil {
		return err
	}
	if err := f.Role.isValid(); err != nil {
		return err
	}
	if strings.TrimSpace(f.Name) == "" {
		return errors.New("name missing")
	}
	if err := isValidURL(f.URL); err != nil {
		return err
	}
	if err := isValidMimeType(f.Type); err != nil {
		return err
	}
	if f.Size <= 0 {
		return errors.New("size must be greater than 0")
	}
	return nil
}

// ValidateResponse validates all File fields including those set by the server (id, created_at, updated_at).
func (f File) ValidateResponse() error {
	if err := isValidUUID(f.ID, "id"); err != nil {
		return err
	}
	if err := f.ValidatePayload(); err != nil {
		return err
	}
	if f.CreatedAt.IsZero() {
		return errors.New("created_at missing")
	}
	if f.UpdatedAt.IsZero() {
		return errors.New("updated_at missing")
	}
	if f.CreatedAt.After(f.UpdatedAt) {
		return errors.New("created_at cannot be after updated_at")
	}
	return nil
}

// -------------------- UPLOADTHING types below --------------------

// FileUpload represents a single file to be uploaded via UploadThing.
type FileUpload struct {
	Name     string  `json:"name"`
	Size     int64   `json:"size"`
	Type     string  `json:"type"`
	CustomID *string `json:"customId,omitempty"`
}

// FileUploadRequest represents a request to upload files via UploadThing with optional metadata and configuration.
type FileUploadRequest struct {
	Files              []FileUpload `json:"files"`
	ACL                *string      `json:"acl,omitempty"`
	Metadata           any          `json:"metadata,omitempty"`
	ContentDisposition *string      `json:"contentDisposition,omitempty"`
}

// Validate validates that the FileUploadRequest has all required fields and valid values.
func (i FileUploadRequest) Validate() error {
	// Files cannot be empty
	if len(i.Files) == 0 {
		return errors.New("files missing")
	}

	for idx, f := range i.Files {
		if f.Name == "" {
			return errors.New("file[" + strconv.Itoa(idx) + "]: name missing")
		}
		if f.Size <= 0 {
			return errors.New("file[" + strconv.Itoa(idx) + "]: size invalid")
		}
		if f.Type == "" {
			return errors.New("file[" + strconv.Itoa(idx) + "]: type missing")
		}
		if err := isValidMimeType(f.Type); err != nil {
			return errors.New("file[" + strconv.Itoa(idx) + "]: " + err.Error())
		}
	}

	// UploadThing ACL accepts: private, public-read
	if i.ACL != nil && *i.ACL != "" {
		if *i.ACL != "public-read" && *i.ACL != "private" {
			return errors.New("acl must be 'public-read' or 'private'")
		}
	}

	// UploadThing ContentDisposition accepts: inline, attachment
	if i.ContentDisposition != nil && *i.ContentDisposition != "" {
		if *i.ContentDisposition != "inline" && *i.ContentDisposition != "attachment" {
			return errors.New("contentDisposition must be 'inline' or 'attachment'")
		}
	}

	return nil
}

// FileUploaded represents a file that has been successfully uploaded via UploadThing.
type FileUploaded struct {
	Key                string         `json:"key"`
	FileName           string         `json:"fileName"`
	FileType           string         `json:"fileType"`
	FileUrl            string         `json:"fileUrl"`
	ContentDisposition string         `json:"contentDisposition"`
	PollingJwt         string         `json:"pollingJwt"`
	PollingUrl         string         `json:"pollingUrl"`
	CustomId           *string        `json:"customId,omitempty"`
	URL                string         `json:"url"` // signed URL to upload
	Fields             map[string]any `json:"fields,omitempty"`
}

// FileUploadedResponse represents the response from UploadThing containing uploaded file information.
type FileUploadedResponse struct {
	Data []FileUploaded `json:"data"`
}

// Validate validates that the FileUploadedResponse has valid uploaded file data.
func (r FileUploadedResponse) Validate() error {
	if len(r.Data) == 0 {
		return errors.New("uploadthing: response returned no files")
	}

	for idx, f := range r.Data {
		prefix := fmt.Sprintf("uploadthing: data[%d]", idx)

		if f.Key == "" {
			return fmt.Errorf("%s.key missing", prefix)
		}
		if f.URL == "" {
			return fmt.Errorf("%s.url missing", prefix)
		}
	}

	return nil
}
