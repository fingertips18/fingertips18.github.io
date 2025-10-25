package v1

import "time"

type ProjectDTO struct {
	Id          string    `json:"id"`
	Preview     string    `json:"preview"`
	BlurHash    string    `json:"blur_hash"`
	Title       string    `json:"title"`
	SubTitle    string    `json:"sub_title"`
	Description string    `json:"description"`
	Stack       []string  `json:"stack"`
	Type        string    `json:"type"`
	Link        string    `json:"link"`
	EducationID string    `json:"education_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProjectRequest struct {
	Preview     string   `json:"preview"`
	BlurHash    string   `json:"blur_hash"`
	Title       string   `json:"title"`
	SubTitle    string   `json:"sub_title"`
	Description string   `json:"description"`
	Stack       []string `json:"stack"`
	Type        string   `json:"type"`
	Link        string   `json:"link"`
	EducationID string   `json:"education_id,omitempty"`
}

type ProjectFilterRequest struct {
	Page          int32  `json:"page"`
	PageSize      int32  `json:"page_size"`
	SortBy        string `json:"sort_by"`
	SortAscending bool   `json:"sort_ascending"`
	Type          string `json:"type"`
}

type SchoolPeriodDTO struct {
	Link        string    `json:"link,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Logo        string    `json:"logo"`
	BlurHash    string    `json:"blur_hash"`
	Honor       string    `json:"honor,omitempty"`
	StartDate   time.Time `json:"start_date" example:"2020-09-01T00:00:00Z"`
	EndDate     time.Time `json:"end_date" example:"2024-06-01T00:00:00Z"`
}

type CreateEducationRequest struct {
	MainSchool    SchoolPeriodDTO   `json:"main_school"`
	SchoolPeriods []SchoolPeriodDTO `json:"school_periods,omitempty"`
	Level         string            `json:"level"`
}

type EducationDTO struct {
	Id            string            `json:"id"`
	MainSchool    SchoolPeriodDTO   `json:"main_school"`
	SchoolPeriods []SchoolPeriodDTO `json:"school_periods,omitempty"`
	Projects      []ProjectDTO      `json:"projects,omitempty"`
	Level         string            `json:"level"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
}

type EducationFilterRequest struct {
	Page          int32  `json:"page"`
	PageSize      int32  `json:"page_size"`
	SortBy        string `json:"sort_by"`
	SortAscending bool   `json:"sort_ascending"`
}

type IDResponse struct {
	Id string `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
