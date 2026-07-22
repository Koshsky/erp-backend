package domain

type UserRole string

const (
	UserRoleProjectDirector UserRole = "ДП"
	UserRoleProjectManager  UserRole = "РП"
	UserRoleProcessOwner    UserRole = "ВП"
)

type User struct {
	ID           int64    `db:"id" json:"id"`
	Name         string   `db:"name" json:"name"`
	Role         UserRole `db:"role" json:"role"`
	Username     string   `db:"username" json:"username"`
	PasswordHash string   `db:"password_hash" json:"-"`
}
