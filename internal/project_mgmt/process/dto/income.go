package dto

import "time"

type UpdateProcessRequest struct {
	Title     *string    `json:"title"      example:"Инсталляция"`
	StartDate *time.Time `json:"start_date" example:"2026-01-01T00:00:00Z"`
	EndDate   *time.Time `json:"end_date"   example:"2026-02-01T00:00:00Z"`
}

type CreateProcessRequest struct {
	ProjectID int64     `json:"project_id" example:"1"`
	Title     string    `json:"title"      example:"Инсталляция"`
	StartDate time.Time `json:"start_date" example:"2026-01-01T00:00:00Z"`
	EndDate   time.Time `json:"end_date"   example:"2026-02-01T00:00:00Z"`
}
