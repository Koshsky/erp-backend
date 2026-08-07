package dto

import "github.com/Koshsky/erp-backend/internal/common/date"

type CreateEmployeeRequest struct {
	Name            string     `json:"name"             example:"Иванов Иван Иванович"`
	ManagerID       *int64     `json:"manager_id"       example:"5"`
	HireDate        *date.Date `json:"hire_date"        example:"2025-01-10"           format:"date"`
	TerminationDate *date.Date `json:"termination_date" example:"2026-12-31"           format:"date"`
}

type UpdateEmployeeRequest struct {
	Name            *string    `json:"name"             example:"Иванов Иван Иванович"`
	ManagerID       *int64     `json:"manager_id"       example:"5"`
	HireDate        *date.Date `json:"hire_date"        example:"2025-01-10"           format:"date"`
	TerminationDate *date.Date `json:"termination_date" example:"2026-12-31"           format:"date"`
}

// SetDaysRequest задаёт состояние на диапазон дат (разворачивается в ячейки календаря).
type SetDaysRequest struct {
	StateID   int64     `json:"state_id"   example:"4"`
	StartDate date.Date `json:"start_date" example:"2026-07-20" format:"date"`
	EndDate   date.Date `json:"end_date"   example:"2026-08-02" format:"date"`
}
