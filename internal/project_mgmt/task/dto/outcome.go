package dto

import "github.com/Koshsky/erp-backend/pkg/date"

type TaskResponse struct {
	ID        int64     `json:"id"`
	ProcessID int64     `json:"process_id"`
	OwnerID   *int64    `json:"owner_id"`
	Title     string    `json:"title"`
	Color     *string   `json:"color"`
	StartDate date.Date `json:"start_date"`
	EndDate   date.Date `json:"end_date"`
	// Order of the task within its process (ascending display order).
	Order int `json:"order" example:"1"`
}

// ReorderTaskRequest — the complete ordered list of the process's active
// tasks; the server rewrites their order values from the list positions.
type ReorderTaskRequest struct {
	ProcessID int64   `json:"process_id" example:"1"`
	IDs       []int64 `json:"ids"`
}
