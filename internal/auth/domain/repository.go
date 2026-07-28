package domain

import "context"

type AuthRepository interface {
	FindUserByUsername(ctx context.Context, username string) (id int64, name, role, passwordHash string, err error)
	FindUserByID(ctx context.Context, userID int64) (id int64, name, role, passwordHash string, err error)
	CreateUser(ctx context.Context, name, username, role, passwordHash string) (int64, error)
	UpdatePassword(ctx context.Context, userID int64, newHash string) error
	SaveRefreshToken(ctx context.Context, userID int64, token string) error
	DeleteRefreshToken(ctx context.Context, userID int64) error
}
type PasswordHasher interface {
	Hash(raw string) (string, error)
	Compare(hashed, raw string) error
}
