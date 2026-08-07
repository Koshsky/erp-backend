package dto

import "github.com/Koshsky/erp-backend/internal/common/date"

type EmployeeResponse struct {
	ID              int64      `json:"id"               example:"1"`
	ResourceID      int64      `json:"resource_id"      example:"3"`
	ResourceTitle   string     `json:"resource_title"   example:"Инженер"`
	Name            string     `json:"name"             example:"Иванов Иван Иванович"`
	Position        string     `json:"position"         example:"Ведущий инженер"`
	ManagerID       *int64     `json:"manager_id"       example:"5"`
	HireDate        *date.Date `json:"hire_date"        example:"2024-01-15"           format:"date"`
	TerminationDate *date.Date `json:"termination_date" example:"2026-12-31"           format:"date"`
}

type EmployeeStateResponse struct {
	ID          int64     `json:"id"           example:"1"`
	StateID     int64     `json:"state_id"     example:"4"`
	StateCode   string    `json:"state_code"   example:"ОТП"`
	StateName   string    `json:"state_name"   example:"Отпуск"`
	IsAvailable bool      `json:"is_available" example:"false"`
	StartDate   date.Date `json:"start_date"   example:"2026-07-20" format:"date"`
	EndDate     date.Date `json:"end_date"     example:"2026-08-02" format:"date"`
}
