package dto

import "github.com/Koshsky/erp-backend/pkg/date"

// CreateUserRequest creates a user (worker with auto-generated credentials when
// username/password_hash are empty). Креды не возвращаются наружу.
type CreateUserRequest struct {
	LastName        string     `json:"last_name"        example:"Иванов"`
	FirstName       string     `json:"first_name"       example:"Иван"`
	MiddleName      *string    `json:"middle_name"      example:"Иванович"`
	Role            string     `json:"role"             example:"worker"`
	Username        string     `json:"username"         example:"worker_1"`
	PasswordHash    string     `json:"password_hash"    example:""`
	ManagerID       *int64     `json:"manager_id"       example:"5"`
	Position        string     `json:"position"         example:"Инженер 2 категории"`
	HireDate        *date.Date `json:"hire_date"        example:"2025-01-10"          format:"date"`
	TerminationDate *date.Date `json:"termination_date" example:"2026-12-31"          format:"date"`
}

type UpdateUserRequest struct {
	LastName        *string    `json:"last_name"        example:"Иванов"`
	FirstName       *string    `json:"first_name"       example:"Иван"`
	MiddleName      *string    `json:"middle_name"      example:"Иванович"`
	Username        *string    `json:"username"         example:"ivanov"`
	Role            *string    `json:"role"             example:"worker"`
	ManagerID       *int64     `json:"manager_id"       example:"5"`
	Position        *string    `json:"position"         example:"Инженер 2 категории"`
	HireDate        *date.Date `json:"hire_date"        example:"2025-01-10"          format:"date"`
	TerminationDate *date.Date `json:"termination_date" example:"2026-12-31"          format:"date"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" example:"password"`
	NewPassword string `json:"new_password" example:"new_password"`
}

// UpdateManagerRequest explicitly sets/clears the manager (null — без руководителя).
type UpdateManagerRequest struct {
	ManagerID *int64 `json:"manager_id" example:"5"`
}

// SetDaysRequest sets a state for a date range (expands into calendar cells).
type SetDaysRequest struct {
	StateID   int64     `json:"state_id"   example:"4"`
	StartDate date.Date `json:"start_date" example:"2026-07-20" format:"date"`
	EndDate   date.Date `json:"end_date"   example:"2026-08-02" format:"date"`
}
