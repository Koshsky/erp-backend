package dto

import "github.com/Koshsky/erp-backend/internal/common/date"

type TaskResponse struct {
	ID        int64     `json:"id"`
	ProcessID int64     `json:"process_id"`
	Title     string    `json:"title"`
	StartDate date.Date `json:"start_date"`
	EndDate   date.Date `json:"end_date"`
}
