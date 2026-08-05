package dto

// UserResponse — user data.
type UserResponse struct {
	ID       int64  `json:"id"       example:"1"`
	Name     string `json:"name"     example:"Ivan Ivanov"`
	Username string `json:"username" example:"ivanov"`
	Role     string `json:"role"     example:"dp"`
	// PasswordHash does not serialize to JSON.
	PasswordHash string `json:"-"`
}

type ChangePasswordResponse struct {
	Message string `json:"message" example:"password changed"`
}
