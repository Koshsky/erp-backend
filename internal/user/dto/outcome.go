package dto

// UserResponse — user data.
type UserResponse struct {
	ID       int64  `json:"id" example:"1"`
	Name     string `json:"name" example:"Ivan Ivanov"`
	Username string `json:"username" example:"ivanov"`
	Role     string `json:"role" example:"ДП"`
	// PasswordHash does not serialize to JSON.
	PasswordHash string `json:"-"`
}
