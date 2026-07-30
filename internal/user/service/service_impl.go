package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/user/dto"
)

type UserService struct {
	logger     *slog.Logger
	repository UserRepository
	mapper     *UserMapper
	validator  *UserValidator
}

func NewUserService(logger *slog.Logger, repository UserRepository) *UserService {
	return &UserService{
		logger:     logger,
		repository: repository,
		mapper:     &UserMapper{},
		validator:  &UserValidator{},
	}
}

func (s *UserService) FindUserByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	return s.FindUser(ctx, id)
}

func (s *UserService) UpdatePassword(ctx context.Context, userID int64, passwordHash string) error {
	return s.repository.UpdatePassword(ctx, userID, passwordHash)
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	if err := s.validator.ValidateUserCreate(req.Name, req.Username, req.Role); err != nil {
		return nil, err
	}
	user := s.mapper.ToDomainFromCreate(req)
	createdUser, err := s.repository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(createdUser), nil
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
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	s.mapper.ApplyUpdateToDomain(user, req)

	if err = s.validator.ValidateUserUpdate(user.Name, user.Username, user.Role); err != nil {
		return nil, err
	}

	updatedUser, err := s.repository.UpdateUser(ctx, *user)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(updatedUser), nil
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
