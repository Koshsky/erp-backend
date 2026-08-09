package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/security/hasher"
	"github.com/Koshsky/erp-backend/internal/user/dto"
)

type UserService struct {
	logger     *slog.Logger
	repository UserRepository
	mapper     *UserMapper
	validator  *UserValidator
}

func (s *UserService) FindUserByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	return s.FindUser(ctx, id)
}

func (s *UserService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	user, err := s.FindUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if err = hasher.Compare(user.PasswordHash, oldPassword); err != nil {
		return fmt.Errorf("invalid current password")
	}

	newHash, err := hasher.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	return s.repository.UpdatePassword(ctx, userID, newHash)
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	user := s.mapper.ToDomainFromCreate(req)
	if err := s.validator.ValidateUser(&user); err != nil {
		return nil, err
	}

	created, err := s.repository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(created), nil
}

func (s *UserService) FindUserByUsername(ctx context.Context, username string) (*dto.UserResponse, error) {
	user, err := s.repository.FindUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(user), nil
}

func (s *UserService) FindUser(ctx context.Context, id int64) (*dto.UserResponse, error) {
	user, err := s.repository.FindUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return s.mapper.ToDTO(user), nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repository.FindUser(ctx, id)
	if err != nil || user == nil {
		return nil, fmt.Errorf("user not found")
	}

	s.mapper.ApplyUpdateToDomain(user, req)
	if err = s.validator.ValidateUser(user); err != nil {
		return nil, err
	}

	updated, err := s.repository.UpdateUser(ctx, *user)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(updated), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	return s.repository.DeleteUser(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := s.repository.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(users), nil
}
