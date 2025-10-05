package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type ProjectType string

const (
	Web    ProjectType = "web"
	Mobile ProjectType = "mobile"
	Game   ProjectType = "game"
)

type Project struct {
	Id          string      `json:"id"`
	Preview     string      `json:"preview"`
	BlurHash    string      `json:"blur_hash"`
	Title       string      `json:"title"`
	SubTitle    string      `json:"sub_title"`
	Description string      `json:"description"`
	Stack       []string    `json:"stack"`
	Type        ProjectType `json:"type"`
	Link        string      `json:"link"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type CreateProject struct {
	Preview     string      `json:"preview"`
	BlurHash    string      `json:"blur_hash"`
	Title       string      `json:"title"`
	SubTitle    string      `json:"sub_title"`
	Description string      `json:"description"`
	Stack       []string    `json:"stack"`
	Type        ProjectType `json:"type"`
	Link        string      `json:"link"`
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

func (p Project) ValidatePayload() error {
	if p.Preview == "" {
		return errors.New("preview missing")
	}
	if p.BlurHash == "" {
		return errors.New("blurHash missing")
	}
	if p.Title == "" {
		return errors.New("title missing")
	}
	if p.SubTitle == "" {
		return errors.New("subTitle missing")
	}
	if p.Description == "" {
		return errors.New("description missing")
	}
	if len(p.Stack) == 0 {
		return errors.New("stack missing")
	}
	for i, item := range p.Stack {
		if strings.TrimSpace(item) == "" {
			return fmt.Errorf("stack[%d] is empty", i)
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

func (p Project) ValidateResponse() error {
	if p.Id == "" {
		return errors.New("ID missing")
	}

	if err := p.ValidatePayload(); err != nil {
		return err
	}

	if p.CreatedAt.IsZero() {
		return errors.New("createdAt missing")
	}

	if p.UpdatedAt.IsZero() {
		return errors.New("updatedAt missing")
	}

	return nil
}
