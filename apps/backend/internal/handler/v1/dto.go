package v1

import "time"

type CreateSkillRequest struct {
	Icon     string `json:"icon"`
	HexColor string `json:"hex_color"`
	Label    string `json:"label"`
	Category string `json:"category"`
}

type UpdateSkillRequest struct {
	Id       string `json:"id"`
	Icon     string `json:"icon"`
	HexColor string `json:"hex_color"`
	Label    string `json:"label"`
	Category string `json:"category"`
}

type UpdateSkillResponse struct {
	Id        string    `json:"id"`
	Icon      string    `json:"icon"`
	HexColor  string    `json:"hex_color"`
	Label     string    `json:"label"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SkillDTO struct {
	Id        string    `json:"id"`
	Icon      string    `json:"icon"`
	HexColor  string    `json:"hex_color"`
	Label     string    `json:"label"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SkillFilterRequest struct {
	Page          int32  `json:"page"`
	PageSize      int32  `json:"page_size"`
	SortBy        string `json:"sort_by"`
	SortAscending bool   `json:"sort_ascending"`
	Category      string `json:"category"`
}

type IDResponse struct {
	Id string `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
