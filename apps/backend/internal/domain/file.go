package domain

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
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

func isValidMimeType(mimeType string) error {
	if strings.TrimSpace(mimeType) == "" {
		return errors.New("type missing")
	}

	// RFC 6838 compliant MIME type pattern
	// Format: type/subtype with optional parameters
	pattern := `^[a-z]+/[a-z0-9][a-z0-9\-\+\.]*$`
	matched, err := regexp.MatchString(pattern, strings.ToLower(mimeType))
	if err != nil {
		return errors.New("invalid type")
	}

	if matched {
		return nil
	} else {
		return errors.New("invalid type")
	}
}

func isValidUUID(id string, label string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("%s missing", label)
	}

	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%s invalid", label)
	}

	return nil
}

func isValidURL(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("url missing")
	}
	if _, err := url.Parse(value); err != nil {
		return errors.New("url invalid")
	}

	return nil
}

func (f File) ValidatePayload() error {
	if strings.TrimSpace(string(f.ParentTable)) == "" {
		return errors.New("parent_table missing")
	}
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

func (f File) ValidateResponse() error {
	if err := isValidUUID(f.ParentID, "id"); err != nil {
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
