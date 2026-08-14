package dto

import "github.com/Koshsky/erp-backend/pkg/date"

type ResourceResponse struct {
	ID             int64  `json:"id"              example:"1"`
	Code           string `json:"code"            example:"М"`
	Title          string `json:"title"           example:"Монтажник"`
	OwnerID        *int64 `json:"owner_id"        example:"3"`
	EmployeesCount int    `json:"employees_count" example:"4"`
}

// ResourceMemberResponse is a user attached to a resource.
type ResourceMemberResponse struct {
	ID              int64      `json:"id"               example:"7"`
	Name            string     `json:"name"             example:"Иванов Иван Иванович"`
	Role            string     `json:"role"             example:"worker"`
	Position        string     `json:"position"         example:"Инженер 2 категории"`
	ManagerID       *int64     `json:"manager_id"       example:"5"`
	HireDate        *date.Date `json:"hire_date"        example:"2024-01-15"           format:"date"`
	TerminationDate *date.Date `json:"termination_date" example:"2026-12-31"           format:"date"`
}
