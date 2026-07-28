package dto

import "time"

type ProcessResponse struct {
	ID        int64     `json:"id"`
	OwnerID   int64     `json:"owner_id"`
	ProjectID int64     `json:"project_id"`
	Title     string    `json:"title"`
	StartDate time.Time `json:"start_date" time_format:"2006-01-02"`
	EndDate   time.Time `json:"end_date" time_format:"2006-01-02"`
}
