package dto

import "github.com/Koshsky/erp-backend/pkg/date"

// UserResponse — user data (incl. worker profile fields).
type UserResponse struct {
	ID              int64      `json:"id"               example:"1"`
	Name            string     `json:"name"             example:"Иванов Иван Иванович"`
	Username        string     `json:"username"         example:"worker_1"`
	Role            string     `json:"role"             example:"worker"`
	ManagerID       *int64     `json:"manager_id"       example:"5"`
	Position        string     `json:"position"         example:"Инженер 2 категории"`
	HireDate        *date.Date `json:"hire_date"        example:"2024-01-15"           format:"date"`
	TerminationDate *date.Date `json:"termination_date" example:"2026-12-31"           format:"date"`
	// PasswordHash never serializes in the regular response (только AdminUserResponse).
	PasswordHash string `json:"-"`
}

// CreateUserResult — созданный пользователь; password возвращается один раз,
// если креды были сгенерированы автоматически.
type CreateUserResult struct {
	User     UserResponse `json:"user"`
	Password string       `json:"password,omitempty"`
}

// AdminUserResponse — пользователь для админ-страницы (включая хеш пароля).
type AdminUserResponse struct {
	ID              int64      `json:"id"                      example:"1"`
	Name            string     `json:"name"                    example:"Иванов Иван Иванович"`
	Username        string     `json:"username"                example:"worker_1"`
	Role            string     `json:"role"                    example:"worker"`
	ManagerID       *int64     `json:"manager_id"              example:"5"`
	Position        string     `json:"position"                example:"Инженер 2 категории"`
	HireDate        *date.Date `json:"hire_date"               example:"2024-01-15"           format:"date"`
	TerminationDate *date.Date `json:"termination_date"        example:"2026-12-31"           format:"date"`
	PasswordHash    string     `json:"password_hash,omitempty" example:"$2a$10$..."`
}

// ResetPasswordResponse — новый сгенерированный пароль (показывается один раз).
type ResetPasswordResponse struct {
	Password string `json:"password" example:"Xy9kLm2QrT8wAb3z"`
}

type UserStateResponse struct {
	ID          int64     `json:"id"           example:"1"`
	StateID     int64     `json:"state_id"     example:"4"`
	StateCode   string    `json:"state_code"   example:"ОТП"`
	StateName   string    `json:"state_name"   example:"Отпуск"`
	IsAvailable bool      `json:"is_available" example:"false"`
	StartDate   date.Date `json:"start_date"   example:"2026-07-20" format:"date"`
	EndDate     date.Date `json:"end_date"     example:"2026-08-02" format:"date"`
}

type ChangePasswordResponse struct {
	Message string `json:"message" example:"password changed"`
}
