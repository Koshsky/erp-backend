package domain

import "time"

type Process struct {
	ID        int64     `json:"id"`
	OwnerID   *int64    `json:"owner_id"`
	ProjectID int64     `json:"project_id"`
	Title     string    `json:"title"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	SortOrder int       `json:"order"`
}
