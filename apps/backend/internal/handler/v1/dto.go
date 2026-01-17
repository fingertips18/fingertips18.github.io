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

type FileDTO struct {
	Name     string  `json:"name"`
	Size     int64   `json:"size"`
	Type     string  `json:"type"`
	CustomID *string `json:"custom_id,omitempty"`
}

type ImageUploadRequestDTO struct {
	Files              []FileDTO `json:"files"`
	ACL                *string   `json:"acl,omitempty"`
	Metadata           any       `json:"metadata,omitempty"`
	ContentDisposition *string   `json:"content_disposition,omitempty"`
}

type ImageUploadFileDTO struct {
	Key                string         `json:"key"`
	FileName           string         `json:"file_name"`
	FileType           string         `json:"file_type"`
	FileUrl            string         `json:"file_url"`
	ContentDisposition string         `json:"content_disposition"`
	PollingJwt         string         `json:"polling_jwt"`
	PollingUrl         string         `json:"polling_url"`
	CustomId           *string        `json:"custom_id,omitempty"`
	URL                string         `json:"url"` // signed URL to upload
	Fields             map[string]any `json:"fields,omitempty"`
}

type ImageUploadResponseDTO struct {
	File ImageUploadFileDTO `json:"file"`
}

type IDResponse struct {
	Id string `json:"id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
