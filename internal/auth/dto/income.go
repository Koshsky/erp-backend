package dto

type LoginRequest struct {
	Username string `json:"username" example:"ivanov"`
	Password string `json:"password" example:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name" example:"Ivan Ivanov"`
	Username string `json:"username" example:"ivanov"`
	Password string `json:"password" example:"password"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" example:"password"`
	NewPassword string `json:"new_password" example:"new_password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6I..."`
}
