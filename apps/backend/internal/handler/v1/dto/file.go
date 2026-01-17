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
