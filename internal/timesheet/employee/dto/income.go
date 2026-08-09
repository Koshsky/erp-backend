package dto

import "github.com/Koshsky/erp-backend/pkg/date"

type CreateEmployeeRequest struct {
	Position        string     `json:"position"         example:"Ведущий инженер"`
	Name            string     `json:"name"             example:"Иванов Иван Иванович"`
	ManagerID       *int64     `json:"manager_id"       example:"5"`
	HireDate        *date.Date `json:"hire_date"        example:"2025-01-10"           format:"date"`
	TerminationDate *date.Date `json:"termination_date" example:"2026-12-31"           format:"date"`
}

type UpdateEmployeeRequest struct {
	ResourceID      *int64     `json:"resource_id"      example:"2"`
	Position        *string    `json:"position"         example:"Ведущий инженер"`
	Name            *string    `json:"name"             example:"Иванов Иван Иванович"`
	ManagerID       *int64     `json:"manager_id"       example:"5"`
	HireDate        *date.Date `json:"hire_date"        example:"2025-01-10"           format:"date"`
	TerminationDate *date.Date `json:"termination_date" example:"2026-12-31"           format:"date"`
}

// SetDaysRequest sets a state for a date range (expands into calendar cells).
type SetDaysRequest struct {
	StateID   int64     `json:"state_id"   example:"4"`
	StartDate date.Date `json:"start_date" example:"2026-07-20" format:"date"`
	EndDate   date.Date `json:"end_date"   example:"2026-08-02" format:"date"`
}
