package domain

import "time"

type Project struct {
	ID        int64     `json:"id"`
	OwnerID   *int64    `json:"owner_id"`
	Code      string    `json:"code"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Priority  int       `json:"priority"`
}

// AutoCreatedCounts is what the auto-create trigger (V8) created for a project
// on insert: process/task/assignment counts (all zero when the template is
// disabled or empty).
type AutoCreatedCounts struct {
	Processes   int64
	Tasks       int64
	Assignments int64
}
