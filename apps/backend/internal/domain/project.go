package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fingertips18/fingertips18.github.io/backend/pkg/metadata"
)

type ProjectType string

const (
	Web    ProjectType = "web"
	Mobile ProjectType = "mobile"
	Game   ProjectType = "game"
)

type Project struct {
	Id          string      `json:"id"`
	BlurHash    string      `json:"blurhash"`
	Title       string      `json:"title"`
	Subtitle    string      `json:"sub_title"`
	Description string      `json:"description"`
	Tags        []string    `json:"tags"`
	Type        ProjectType `json:"type"`
	Link        string      `json:"link"`
	EducationID string      `json:"education_id,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type ProjectFilter struct {
	Page          int32
	PageSize      int32
	SortBy        *SortBy
	SortAscending bool
	Type          *ProjectType
}

type ProjectIDResponse struct {
	ID string `json:"id"`
}

func (pt ProjectType) isValid() bool {
	switch pt {
	case Web, Mobile, Game:
		return true
	default:
		return false
	}
}

func (p Project) ValidatePayload(blurHashAPI metadata.BlurHashAPI) error {
	if p.BlurHash == "" {
		return errors.New("blurHash missing")
	}
	if !blurHashAPI.IsValid(p.BlurHash) {
		return errors.New("blurHash invalid")
	}
	if p.Title == "" {
		return errors.New("title missing")
	}
	if p.Subtitle == "" {
		return errors.New("subTitle missing")
	}
	if p.Description == "" {
		return errors.New("description missing")
	}
	if len(p.Tags) == 0 {
		return errors.New("tags missing")
	}
	for i, item := range p.Tags {
		if strings.TrimSpace(item) == "" {
			return fmt.Errorf("tag[%d] is empty", i)
		}
	}
	if p.Type == "" {
		return errors.New("type missing")
	} else if !p.Type.isValid() {
		return fmt.Errorf("type invalid = %s", p.Type)
	}
	if p.Link == "" {
		return errors.New("link missing")
	}

	return nil
}

func (p Project) ValidateResponse(blurHashAPI metadata.BlurHashAPI) error {
	if p.Id == "" {
		return errors.New("ID missing")
	}

	if err := p.ValidatePayload(blurHashAPI); err != nil {
		return err
	}

	if p.CreatedAt.IsZero() {
		return errors.New("createdAt missing")
	}

	if p.UpdatedAt.IsZero() {
		return errors.New("updatedAt missing")
	}

	if p.UpdatedAt.Before(p.CreatedAt) {
		return errors.New("updatedAt before createdAt")
	}

	return nil
}
