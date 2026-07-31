package dto

import "time"

type UpdateMilestoneRequest struct {
	Title   *string    `json:"title"   example:"Телевидение"`
	Content *string    `json:"content" example:"Приедут с России1"`
	Date    *time.Time `json:"date"    example:"2026-01-01T00:00:00Z"`
}

type CreateMilestoneRequest struct {
	ProcessID int64     `json:"process_id" example:"1"`
	Title     string    `json:"title"      example:"Телевидение"`
	Content   string    `json:"content"    example:"Приедут с России1"`
	Date      time.Time `json:"date"       example:"2026-01-01T00:00:00Z"`
}
