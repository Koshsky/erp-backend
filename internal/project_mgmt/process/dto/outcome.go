package dto

import "time"

type ProcessResponse struct {
	ID        int64     `json:"id"         example:"1"`
	OwnerID   int64     `json:"owner_id"   example:"1"`
	ProjectID int64     `json:"project_id" example:"1"`
	Title     string    `json:"title"      example:"Инсталляция"`
	StartDate time.Time `json:"start_date" example:"2026-01-01T00:00:00Z"`
	EndDate   time.Time `json:"end_date"   example:"2026-02-01T00:00:00Z"`
}
