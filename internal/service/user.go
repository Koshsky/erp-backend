package service

import (
	"context"
	"log/slog"
	"strings"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/dto"
	"github.com/Koshsky/erp/api/internal/security/password"
	"github.com/Koshsky/erp/api/internal/service/mapper"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) (*domain.User, error)
	GetUser(ctx context.Context, id int64) (*domain.User, error)
	UpdateUser(ctx context.Context, new domain.User) (*domain.User, error)
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context) ([]domain.User, error)
}

type UserService struct {
	logger     *slog.Logger
	repository UserRepository
	mapper     *mapper.UserMapper
	validator  *Validator
	hasher     password.Hasher
}

func NewUserService(logger *slog.Logger, repository UserRepository, validator *Validator, hasher password.Hasher) *UserService {
	return &UserService{
		logger:     logger,
		repository: repository,
		mapper:     mapper.NewUserMapper(),
		validator:  validator,
		hasher:     hasher,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*dto.UserResponse, error) {
	if err := s.validator.ValidateUserCreate(req.Name, req.Username, domain.UserRole(req.Role), req.Password); err != nil {
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
		return nil, ErrUserNotFound
	}
	return s.mapper.ToDTO(user), nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, req dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.repository.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
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
