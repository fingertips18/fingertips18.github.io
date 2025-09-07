package domain

import "time"

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
