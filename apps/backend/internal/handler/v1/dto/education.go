package dto

import "time"

type SchoolPeriodDTO struct {
	Link        string    `json:"link,omitempty"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Logo        string    `json:"logo"`
	BlurHash    string    `json:"blurhash"`
	Honor       string    `json:"honor,omitempty"`
	StartDate   time.Time `json:"start_date" example:"2020-09-01T00:00:00Z"`
	EndDate     time.Time `json:"end_date" example:"2024-06-01T00:00:00Z"`
}

type CreateEducationRequest struct {
	MainSchool    SchoolPeriodDTO   `json:"main_school"`
	SchoolPeriods []SchoolPeriodDTO `json:"school_periods,omitempty"`
	Level         string            `json:"level" example:"elementary"`
}

type UpdateEducationRequest struct {
	Id            string            `json:"id"`
	MainSchool    SchoolPeriodDTO   `json:"main_school"`
	SchoolPeriods []SchoolPeriodDTO `json:"school_periods,omitempty"`
	Level         string            `json:"level" example:"elementary"`
}

type UpdateEducationResponse struct {
	Id            string            `json:"id"`
	MainSchool    SchoolPeriodDTO   `json:"main_school"`
	SchoolPeriods []SchoolPeriodDTO `json:"school_periods,omitempty"`
	Level         string            `json:"level"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
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
