package dto

import "time"

type TaskResponse struct {
	ID        int64     `json:"id"`
	ProcessID int64     `json:"process_id"`
	Title     string    `json:"title"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"  `
}
