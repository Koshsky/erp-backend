package dto

type LoginRequest struct {
	Username string `json:"username" example:"ivanov"`
	Password string `json:"password" example:"password"`
}

type RegisterRequest struct {
	LastName   string `json:"last_name"   example:"Иванов"`
	FirstName  string `json:"first_name"  example:"Иван"`
	MiddleName string `json:"middle_name" example:"Иванович"`
	Username   string `json:"username"    example:"ivanov"`
	Password   string `json:"password"    example:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6I..."`
}
