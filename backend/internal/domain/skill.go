package domain

import (
	"errors"
	"regexp"
	"time"
)

type Skill struct {
	Id        string    `json:"id"`
	Icon      string    `json:"icon"`
	HexColor  string    `json:"hex_color"`
	Label     string    `json:"label"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var hexColorRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{3}([0-9A-Fa-f]{3})?$`)

func (s Skill) ValidatePayload() error {
	if s.Icon == "" {
		return errors.New("icon missing")
	}
	if s.HexColor == "" {
		return errors.New("hex color missing")
	} else {
		if !hexColorRegex.MatchString(s.HexColor) {
			return errors.New("hex color must be in format #RGB or #RRGGBB")
		}
	}
	if s.Label == "" {
		return errors.New("label missing")
	}
	if s.Category == "" {
		return errors.New("category missing")
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

	return nil
}
