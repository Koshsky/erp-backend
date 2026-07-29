package domain

import "time"

type Milestone struct {
	ID        int64     `json:"id"`
	ProcessID int64     `json:"process_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Date      time.Time `json:"date"`
}
