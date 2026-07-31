package dto

type CreateUserRequest struct {
	Name         string `json:"name"`
	Username     string `json:"username"`
	Role         string `json:"role"`
	PasswordHash string `json:"password_hash"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name"     example:"Ivan Ivanov"`
	Username *string `json:"username" example:"ivanov"`
	Role     *string `json:"role"     example:"ДП"`
}
