package dto

import "github.com/Koshsky/erp-backend/internal/common/date"

type ProjectResponse struct {
	ID        int64     `json:"id"`
	OwnerID   int64     `json:"owner_id"`
	Code      string    `json:"code"`
	StartDate date.Date `json:"start_date"`
	EndDate   date.Date `json:"end_date"`
	Priority  int       `json:"priority"`
}
