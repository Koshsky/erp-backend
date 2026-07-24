package domain

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Koshsky/erp/api/internal/security/password"
	"github.com/Koshsky/erp/api/internal/user/dto"
)

type UserService struct {
	logger     *slog.Logger
	repository RepositoryInterface
	mapper     *UserMapper
	validator  *UserValidator
	hasher     password.Hasher
}

func NewUserService(logger *slog.Logger, repository RepositoryInterface, hasher password.Hasher) *UserService {
	return &UserService{
		logger:     logger,
		repository: repository,
		mapper:     &UserMapper{},
		validator:  &UserValidator{},
		hasher:     hasher,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	if err := s.validator.ValidateUserCreate(req.Name, req.Username, UserRole(req.Role), req.Password); err != nil {
		return nil, err
	}
	user := s.mapper.ToDomainFromCreate(req)
	hash, err := s.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = hash
	createdUser, err := s.repository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(createdUser), nil
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*dto.UserResponse, error) {
	user, err := s.repository.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return s.mapper.ToDTO(user), nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repository.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	s.mapper.ApplyUpdateToDomain(user, req)

	if req.Password != nil && strings.TrimSpace(*req.Password) != "" {
		hash, err := s.hasher.Hash(*req.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hash
	}

	if err := s.validator.ValidateUserUpdate(user.Name, user.Username, user.Role); err != nil {
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
