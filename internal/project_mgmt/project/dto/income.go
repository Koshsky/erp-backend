package dto

import "time"

type CreateProjectRequest struct {
	Code      string    `json:"code"       example:"КО_001"`
	StartDate time.Time `json:"start_date" example:"2026-01-01T00:00:00Z"`
	EndDate   time.Time `json:"end_date"   example:"2026-02-01T00:00:00Z"`
	Priority  int       `json:"priority"   example:"2"`
}

type UpdateProjectRequest struct {
	Code      *string    `json:"code"       example:"1"`
	StartDate *time.Time `json:"start_date" example:"2026-01-01T00:00:00Z"`
	EndDate   *time.Time `json:"end_date"   example:"2026-02-01T00:00:00Z"`
	Priority  *int       `json:"priority"   example:"2"`
}
