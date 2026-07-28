package dto

import "time"

type UpdateProcessRequest struct {
	Title     *string    `json:"title"`
	StartDate *time.Time `json:"start_date" time_format:"2006-01-02"`
	EndDate   *time.Time `json:"end_date" time_format:"2006-01-02"`
}

type CreateProcessRequest struct {
	ProjectID int64     `json:"project_id"`
	Title     string    `json:"title"`
	StartDate time.Time `json:"start_date" time_format:"2006-01-02"`
	EndDate   time.Time `json:"end_date" time_format:"2006-01-02"`
}
