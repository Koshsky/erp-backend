package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/user/dto"
)

type UserService interface {
	ListUsers(ctx context.Context) ([]dto.UserResponse, error)
	FindUser(ctx context.Context, id int64) (*dto.UserResponse, error)
	ListSubordinates(ctx context.Context, managerID int64) ([]dto.UserResponse, error)
	DeleteUser(ctx context.Context, id int64) error
	ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
	UpdateUser(ctx context.Context, id int64, user dto.UpdateUserRequest) (*dto.UserResponse, error)
}
