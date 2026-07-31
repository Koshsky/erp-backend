package dto

import "time"

type ProjectResponse struct {
	ID        int64     `json:"id"`
	OwnerID   int64     `json:"owner_id"`
	Code      string    `json:"code"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"  `
	Priority  int       `json:"priority"`
}
