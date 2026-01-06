package domain

import (
	"time"

	"github.com/google/uuid"
)

type ParentTableType string

const (
	ProjectTable   ParentTableType = "projects"
	UserTable      ParentTableType = "users"
	EducationTable ParentTableType = "educations"
	// Add other valid parent table names as needed
)

// File represents a file attachment that can be associated with any parent entity.
// It uses a polymorphic association pattern via ParentTable and ParentID fields.
type File struct {
	ID          uuid.UUID
	ParentTable ParentTableType
	ParentID    uuid.UUID
	Role        string
	Key         string
	URL         string
	Type        string
	Size        int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (pt ParentTableType) isValid() bool {
	switch pt {
	case ProjectTable, UserTable, EducationTable:
		return true
	default:
		return false
	}
}
