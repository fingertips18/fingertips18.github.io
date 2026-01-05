package domain

type File struct {
	ID          uuid.UUID
	ParentTable string
	ParentID    uuid.UUID
	Role        string
	Key         string
	URL         string
	Type        string
	Size        int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
