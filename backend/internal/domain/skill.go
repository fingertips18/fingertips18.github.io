package domain

import (
	"errors"
	"time"
)

type CreateSkill struct {
	Icon      string    `json:"icon"`
	HexColor  string    `json:"hex_color"`
	Label     string    `json:"label"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s CreateSkill) ValidatePayload() error {
	if s.Icon == "" {
		return errors.New("icon missing")
	}
	if s.HexColor == "" {
		return errors.New("hex color missing")
	}
	if s.Label == "" {
		return errors.New("label missing")
	}
	if s.Category == "" {
		return errors.New("category missing")
	}

	return nil
}
