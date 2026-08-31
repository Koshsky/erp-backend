package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/user/dto"
	"github.com/Koshsky/erp-backend/pkg/date"
)

type UserService interface {
	ListAllUsers(ctx context.Context) ([]dto.UserResponse, error)
	ListUsers(
		ctx context.Context,
		userID int64,
		role string,
		roleFilter string,
		managerID int64,
		limit, offset int,
	) ([]dto.UserResponse, int64, error)
	FindUser(ctx context.Context, id int64) (*dto.UserResponse, error)
	CreateUserWithCreds(
		ctx context.Context,
		req dto.CreateUserRequest,
		callerRole string,
	) (*dto.CreateUserResult, error)
	ResetPassword(ctx context.Context, id int64) (*dto.ResetPasswordResponse, error)
	DeleteUser(ctx context.Context, id int64) error
	ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
	UpdateUser(
		ctx context.Context,
		id int64,
		user dto.UpdateUserRequest,
		callerRole string,
		callerID int64,
	) (*dto.UserResponse, error)
	UpdateManager(
		ctx context.Context,
		id int64,
		managerID *int64,
	) (*dto.UserResponse, error)

	ListStates(ctx context.Context, userID int64, start, end date.Date) ([]dto.UserStateResponse, error)
	SetDays(ctx context.Context, userID int64, req dto.SetDaysRequest) error
	DeleteDays(ctx context.Context, userID int64, start, end date.Date, stateID *int64) error
}
