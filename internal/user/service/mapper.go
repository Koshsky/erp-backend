package service

import (
	"github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/internal/user/dto"
)

type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

// Обычные мапперы
func (m *UserMapper) ToDTO(user *domain.User) *dto.UserResponse {
	if user == nil {
		return nil
	}
	return &dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Role:     string(user.Role),
	}
}

func (m *UserMapper) ToDTOs(users []domain.User) []dto.UserResponse {
	if users == nil {
		return []dto.UserResponse{}
	}
	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = *m.ToDTO(&user)
	}
	return responses
}

// ToDomainFromCreate maps request fields except password hashing.
func (m *UserMapper) ToDomainFromCreate(req dto.CreateUserRequest) domain.User {
	return domain.User{
		Name:     req.Name,
		Username: req.Username,
		Role:     domain.UserRole(req.Role),
	}
}

// ApplyUpdateToDomain applies mutable profile fields; password is handled by service.
func (m *UserMapper) ApplyUpdateToDomain(user *domain.User, req dto.UpdateUserRequest) {
	if user == nil {
		return
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Role != nil {
		user.Role = domain.UserRole(*req.Role)
	}
}
