package dto

import "time"

type CreateFileRequest struct {
	ParentTable string `json:"parent_table"`
	ParentID    string `json:"parent_id"`
	Role        string `json:"role"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	Size        int64  `json:"size"`
}

type FileDTO struct {
	ID          string    `json:"id"`
	ParentTable string    `json:"parent_table"`
	ParentID    string    `json:"parent_id"`
	Role        string    `json:"role"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Type        string    `json:"type"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type FileUploadDTO struct {
	Name     string  `json:"name"`
	Size     int64   `json:"size"`
	Type     string  `json:"type"`
	CustomID *string `json:"custom_id,omitempty"`
}

type FileUploadRequestDTO struct {
	Files              []FileUploadDTO `json:"files"`
	ACL                *string         `json:"acl,omitempty"`
	Metadata           any             `json:"metadata,omitempty"`
	ContentDisposition *string         `json:"content_disposition,omitempty"`
}

type FileUploadedDTO struct {
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

type FileUploadedResponseDTO struct {
	File FileUploadedDTO `json:"file"`
}
