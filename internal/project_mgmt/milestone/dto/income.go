package dto

import "github.com/Koshsky/erp-backend/internal/common/date"

type UpdateMilestoneRequest struct {
	Title   *string    `json:"title"   example:"Телевидение"`
	Content *string    `json:"content" example:"Приедут с России1"`
	Date    *date.Date `json:"date"    example:"2026-01-01"        format:"date"`
}

type CreateMilestoneRequest struct {
	ProcessID int64     `json:"process_id" example:"1"`
	Title     string    `json:"title"      example:"Телевидение"`
	Content   string    `json:"content"    example:"Приедут с России1"`
	Date      date.Date `json:"date"       example:"2026-01-01"        format:"date"`
}
