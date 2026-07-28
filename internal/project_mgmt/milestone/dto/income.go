package dto

import "time"

type UpdateMilestoneRequest struct {
	Title   *string    `json:"title"`
	Content *string    `json:"content"`
	Date    *time.Time `json:"date" time_format:"2006-01-02"`
}

type CreateMilestoneRequest struct {
	ProcessID int64     `json:"process_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Date      time.Time `json:"date" time_format:"2006-01-02"`
}
