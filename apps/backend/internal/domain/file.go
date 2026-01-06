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
	ID          uuid.UUID       `json:"id"`
	ParentTable ParentTableType `json:"parent_table"`
	ParentID    uuid.UUID       `json:"parent_id"`
	Role        string          `json:"role,omitempty"`
	Key         string          `json:"key"`
	URL         string          `json:"url"`
	Type        string          `json:"type"`
	Size        int64           `json:"size"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (pt ParentTableType) isValid() bool {
	switch pt {
	case ProjectTable, UserTable, EducationTable:
		return true
	default:
		return false
	}
}
