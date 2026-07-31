package dto

type LoginRequest struct {
	Username string `json:"username" example:"ivanov"`
	Password string `json:"password" example:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"     example:"Ivan Ivanov"`
	Username string `json:"username" example:"ivanov"`
	Password string `json:"password" example:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6I..."`
}
