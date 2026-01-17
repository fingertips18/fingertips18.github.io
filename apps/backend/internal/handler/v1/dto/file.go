package dto

import "time"

type CreateFileRequest struct {
	ParentTable string `json:"parent_table" validate:"required"`
	ParentID    string `json:"parent_id" validate:"required"`
	Role        string `json:"role" validate:"required"`
	Name        string `json:"name" validate:"required"`
	URL         string `json:"url" validate:"required,url"`
	Type        string `json:"type" validate:"required"`
	Size        int64  `json:"size" validate:"gte=0"`
}

type FileResponse struct {
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
