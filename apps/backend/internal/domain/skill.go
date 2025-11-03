package domain

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

type SkillCategory string

const (
	Frontend SkillCategory = "frontend"
	Backend  SkillCategory = "backend"
	Tools    SkillCategory = "tools"
	Others   SkillCategory = "others"
)

func (sc SkillCategory) isValid() bool {
	switch sc {
	case Frontend, Backend, Tools, Others:
		return true
	default:
		return false
	}
}

type Skill struct {
	Id        string        `json:"id"`
	Icon      string        `json:"icon"`
	HexColor  string        `json:"hex_color"`
	Label     string        `json:"label"`
	Category  SkillCategory `json:"category"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type SkillFilter struct {
	Page          int32
	PageSize      int32
	SortBy        *SortBy
	SortAscending bool
	Category      *SkillCategory
}

var hexColorRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{3}([0-9A-Fa-f]{3})?$`)

func (s Skill) ValidatePayload() error {
	if s.Icon == "" {
		return errors.New("icon missing")
	}
	if s.HexColor == "" {
		return errors.New("hex color missing")
	}
	if !hexColorRegex.MatchString(s.HexColor) {
		return errors.New("hex color must be in format #RGB or #RRGGBB")
	}
	if s.Label == "" {
		return errors.New("label missing")
	}
	if s.Category == "" {
		return errors.New("category missing")
	}

	if !s.Category.isValid() {
		return fmt.Errorf("category invalid = %s", s.Category)
	}

	return nil
}

func (s Skill) ValidateResponse() error {
	if s.Id == "" {
		return errors.New("ID missing")
	}

	if err := s.ValidatePayload(); err != nil {
		return err
	}

	if s.CreatedAt.IsZero() {
		return errors.New("createdAt missing")
	}

	if s.UpdatedAt.IsZero() {
		return errors.New("updatedAt missing")
	}

	if s.UpdatedAt.Before(s.CreatedAt) {
		return errors.New("updatedAt before createdAt")
	}

	return nil
}
