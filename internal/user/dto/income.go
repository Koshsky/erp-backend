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

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" example:"password"`
	NewPassword string `json:"new_password" example:"new_password"`
}
