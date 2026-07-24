package domain

import (
	"context"
)

type RepositoryInterface interface {
	CreateUser(ctx context.Context, user User) (*User, error)
	GetUser(ctx context.Context, id int64) (*User, error)
	UpdateUser(ctx context.Context, new User) (*User, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context) ([]User, error)
}
