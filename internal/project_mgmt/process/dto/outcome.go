package dto

import "github.com/Koshsky/erp-backend/pkg/date"

type ProcessResponse struct {
	ID        int64     `json:"id"         example:"1"`
	OwnerID   *int64    `json:"owner_id"   example:"1"`
	ProjectID int64     `json:"project_id" example:"1"`
	Title     string    `json:"title"      example:"Инсталляция"`
	StartDate date.Date `json:"start_date" example:"2026-01-01"  format:"date"`
	EndDate   date.Date `json:"end_date"   example:"2026-02-01"  format:"date"`
}
