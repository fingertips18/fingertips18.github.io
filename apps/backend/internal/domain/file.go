package domain

import (
	"errors"
	"time"
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

func (pt ParentTable) isValid() bool {
	switch pt {
	case ProjectTable, UserTable, EducationTable:
		return true
	default:
		return false
	}
}

func (fr FileRole) isValid() bool {
	switch fr {
	case Image:
		return true
	default:
		return false
	}
}

func (f File) ValidatePayload() error {
	if f.ParentTable == "" {
		return errors.New("parentTable missing")
	}
	if !f.ParentTable.isValid() {
		return errors.New("parentTable invalid")
	}
	if f.ParentID == "" {
		return errors.New("parentId missing")
	}
	if f.Role == "" {
		return errors.New("role missing")
	}
	if !f.Role.isValid() {
		return errors.New("role invalid")
	}
	if f.Name == "" {
		return errors.New("name missing")
	}
	if f.URL == "" {
		return errors.New("url missing")
	}
	if f.Type == "" {
		return errors.New("type missing")
	}
	if f.Size <= 0 {
		return errors.New("size must be greater than 0")
	}
	return nil
}

func (f File) ValidateResponse() error {
	if f.ID == "" {
		return errors.New("id missing")
	}
	if err := f.ValidatePayload(); err != nil {
		return err
	}
	if f.CreatedAt.IsZero() {
		return errors.New("createdAt missing")
	}
	if f.UpdatedAt.IsZero() {
		return errors.New("updatedAt missing")
	}
	if f.CreatedAt.After(f.UpdatedAt) {
		return errors.New("createdAt cannot be after updatedAt")
	}
	return nil
}
