package domain

import (
	"errors"
	"fmt"
	"time"
)

type EducationLevel string

const (
	Elementary       EducationLevel = "elementary"
	JuniorHighSchool EducationLevel = "junior-high-school"
	SeniorHighSchool EducationLevel = "senior-high-school"
	College          EducationLevel = "college"
)

type SchoolPeriod struct {
	Link        string    `json:"link,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Logo        string    `json:"logo"`
	BlurHash    string    `json:"blurhash"`
	Honor       string    `json:"honor,omitempty"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type Education struct {
	Id            string
	MainSchool    SchoolPeriod
	SchoolPeriods []SchoolPeriod
	Level         EducationLevel
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type EducationFilter struct {
	Page          int32
	PageSize      int32
	SortBy        *SortBy
	SortAscending bool
}

func (el EducationLevel) isValid() bool {
	switch el {
	case Elementary, JuniorHighSchool, SeniorHighSchool, College:
		return true
	default:
		return false
	}
}

func (s SchoolPeriod) Validate() error {
	if s.Name == "" {
		return errors.New("name missing")
	}
	if s.Description == "" {
		return errors.New("description missing")
	}
	if s.Logo == "" {
		return errors.New("logo missing")
	}
	if s.BlurHash == "" {
		return errors.New("blurHash missing")
	}
	if s.StartDate.IsZero() {
		return errors.New("start date missing")
	}
	if s.EndDate.IsZero() {
		return errors.New("end date missing")
	}
	if !s.EndDate.After(s.StartDate) {
		return errors.New("end date must be after start date")
	}
	return nil
}

func (e Education) ValidatePayload() error {
	if e.MainSchool == (SchoolPeriod{}) {
		return errors.New("main school missing")
	}
	if err := e.MainSchool.Validate(); err != nil {
		return fmt.Errorf("main school %w", err)
	}

	for i, sp := range e.SchoolPeriods {
		if sp == (SchoolPeriod{}) {
			return fmt.Errorf("school period[%d] is empty", i)
		}
		if err := sp.Validate(); err != nil {
			return fmt.Errorf("school period[%d] %w", i, err)
		}
	}

	if e.Level == "" {
		return errors.New("level missing")
	}
	if !e.Level.isValid() {
		return fmt.Errorf("level invalid = %s", e.Level)
	}

	return nil
}

func (e Education) ValidateResponse() error {
	if e.Id == "" {
		return errors.New("ID missing")
	}

	if err := e.ValidatePayload(); err != nil {
		return err
	}

	if e.CreatedAt.IsZero() {
		return errors.New("createdAt missing")
	}

	if e.UpdatedAt.IsZero() {
		return errors.New("updatedAt missing")
	}

	if e.UpdatedAt.Before(e.CreatedAt) {
		return errors.New("updatedAt before createdAt")
	}

	return nil
}
