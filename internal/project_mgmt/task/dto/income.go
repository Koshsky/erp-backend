package dto

import "github.com/Koshsky/erp-backend/internal/common/date"

type UpdateTaskRequest struct {
	ProcessID *int64     `json:"process_id" example:"1"`
	OwnerID   *int64     `json:"owner_id"   example:"1"`
	Title     *string    `json:"title"      example:"Пуско-наладочные работы"`
	StartDate *date.Date `json:"start_date" example:"2026-01-01"              format:"date"`
	EndDate   *date.Date `json:"end_date"   example:"2026-02-01"              format:"date"`
}

type CreateTaskRequest struct {
	ProcessID int64     `json:"process_id" example:"1"`
	OwnerID   *int64    `json:"owner_id"   example:"1"`
	Title     string    `json:"title"      example:"Пуско-наладочные работы"`
	StartDate date.Date `json:"start_date" example:"2026-01-01"              format:"date"`
	EndDate   date.Date `json:"end_date"   example:"2026-02-01"              format:"date"`
}
