package dto

import "time"

type ProjectDTO struct {
	Id          string    `json:"id"`
	BlurHash    string    `json:"blurhash"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"sub_title"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Type        string    `json:"type"`
	Link        string    `json:"link"`
	EducationID string    `json:"education_id,omitempty"`
	Previews    []FileDTO `json:"previews"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProjectRequest struct {
	Previews    []CreateFileRequest `json:"previews"`
	BlurHash    string              `json:"blurhash"`
	Title       string              `json:"title"`
	Subtitle    string              `json:"sub_title"`
	Description string              `json:"description"`
	Tags        []string            `json:"tags"`
	Type        string              `json:"type"`
	Link        string              `json:"link"`
	EducationID string              `json:"education_id,omitempty"`
}

type ProjectFilterRequest struct {
	Page          int32  `json:"page"`
	PageSize      int32  `json:"page_size"`
	SortBy        string `json:"sort_by"`
	SortAscending bool   `json:"sort_ascending"`
	Type          string `json:"type"`
}
