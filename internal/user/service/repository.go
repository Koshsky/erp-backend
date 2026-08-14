package service

import (
	"context"
	"time"

	"github.com/Koshsky/erp-backend/internal/user/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) (*domain.User, error)
	FindUser(ctx context.Context, id int64) (*domain.User, error)
	FindUserByUsername(ctx context.Context, username string) (*domain.User, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	UpdateUser(ctx context.Context, user domain.User) (*domain.User, error)
	UpdatePassword(ctx context.Context, userID int64, userHash string) error
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(
		ctx context.Context,
		userID int64,
		role string,
		roleFilter string,
		managerID int64,
		limit, offset int,
	) ([]domain.User, error)
	CountUsers(ctx context.Context, userID int64, role string, roleFilter string, managerID int64) (int64, error)
	ListAllUsers(ctx context.Context) ([]domain.User, error)

	ListStates(ctx context.Context, userID int64, start, end time.Time) ([]domain.UserState, error)
	SetStateRange(ctx context.Context, userID, stateID int64, start, end time.Time) error
	DeleteStateRange(ctx context.Context, userID int64, start, end time.Time, stateID *int64) error
}
