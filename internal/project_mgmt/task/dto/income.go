package dto

import "time"

type UpdateTaskRequest struct {
	Title     *string    `json:"title"      example:"Пуско-наладочные работы"`
	StartDate *time.Time `json:"start_date" example:"2026-01-01T00:00:00Z"`
	EndDate   *time.Time `json:"end_date"   example:"2026-02-01T00:00:00Z"`
}

type CreateTaskRequest struct {
	ProcessID int64     `json:"process_id" example:"1"`
	Title     string    `json:"title"      example:"Пуско-наладочные работы"`
	StartDate time.Time `json:"start_date" example:"2026-01-01T00:00:00Z"`
	EndDate   time.Time `json:"end_date"   example:"2026-02-01T00:00:00Z"`
}
