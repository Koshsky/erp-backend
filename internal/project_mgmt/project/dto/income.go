package dto

import "time"

type CreateProjectRequest struct {
	Code      string    `json:"code"`
	StartDate time.Time `json:"start_date" time_format:"2006-01-02"`
	EndDate   time.Time `json:"end_date"   time_format:"2006-01-02"`
	Priority  int       `json:"priority"`
}

type UpdateProjectRequest struct {
	Code      *string    `json:"code"`
	StartDate *time.Time `json:"start_date" time_format:"2006-01-02"`
	EndDate   *time.Time `json:"end_date"   time_format:"2006-01-02"`
	Priority  *int       `json:"priority"`
}
