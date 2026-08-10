package dto

import "github.com/Koshsky/erp-backend/pkg/date"

type CreateProjectRequest struct {
	OwnerID   *int64    `json:"owner_id"   example:"1"`
	Code      string    `json:"code"       example:"КО_001"`
	StartDate date.Date `json:"start_date" example:"2026-01-01" format:"date"`
	EndDate   date.Date `json:"end_date"   example:"2026-02-01" format:"date"`
	Priority  int       `json:"priority"   example:"2"`
}

type UpdateProjectRequest struct {
	OwnerID   *int64     `json:"owner_id"   example:"1"`
	Code      *string    `json:"code"       example:"1"`
	StartDate *date.Date `json:"start_date" example:"2026-01-01" format:"date"`
	EndDate   *date.Date `json:"end_date"   example:"2026-02-01" format:"date"`
	Priority  *int       `json:"priority"   example:"2"`
}
