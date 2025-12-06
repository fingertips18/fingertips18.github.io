package v1

import "time"

type ProjectDTO struct {
	Id          string    `json:"id"`
	Preview     string    `json:"preview"`
	BlurHash    string    `json:"blurhash"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"sub_title"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Type        string    `json:"type"`
	Link        string    `json:"link"`
	EducationID string    `json:"education_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProjectRequest struct {
	Preview     string   `json:"preview"`
	BlurHash    string   `json:"blurhash"`
	Title       string   `json:"title"`
	Subtitle    string   `json:"sub_title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
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
