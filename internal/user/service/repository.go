package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/user/domain"
)

type RepositoryInterface interface {
	CreateUser(ctx context.Context, user domain.User) (*domain.User, error)
	GetUser(ctx context.Context, id int64) (*domain.User, error)
	FindUserByUsername(ctx context.Context, username string) (*domain.User, error)
	UpdateUser(ctx context.Context, new domain.User) (*domain.User, error)
	UpdatePassword(ctx context.Context, userID int64, userHash string) error
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context) ([]domain.User, error)
}

