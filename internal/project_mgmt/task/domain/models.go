package domain

import "time"

type Task struct {
	ID        int64     `json:"id"`
	ProcessID int64     `json:"process_id"`
	OwnerID   *int64    `json:"owner_id"`
	Title     string    `json:"title"`
	Color     *string   `json:"color"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	SortOrder int       `json:"order"`
}
