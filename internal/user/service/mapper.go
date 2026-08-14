package service

import (
	"time"

	"github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/internal/user/dto"
	"github.com/Koshsky/erp-backend/pkg/date"
)

type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (m *UserMapper) ToDTO(user *domain.User) *dto.UserResponse {
	if user == nil {
		return nil
	}
	return &dto.UserResponse{
		ID:              user.ID,
		Name:            user.Name,
		Username:        user.Username,
		Role:            user.Role,
		ManagerID:       user.ManagerID,
		Position:        user.Position,
		HireDate:        datePtr(user.HireDate),
		TerminationDate: datePtr(user.TerminationDate),
		PasswordHash:    user.PasswordHash,
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

func (m *UserMapper) ToDomainFromCreate(req dto.CreateUserRequest) domain.User {
	return domain.User{
		Name:            req.Name,
		Username:        req.Username,
		Role:            req.Role,
		PasswordHash:    req.PasswordHash,
		ManagerID:       req.ManagerID,
		Position:        req.Position,
		HireDate:        timePtr(req.HireDate),
		TerminationDate: timePtr(req.TerminationDate),
	}
}

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
		user.Role = *req.Role
	}
	if req.ManagerID != nil {
		user.ManagerID = req.ManagerID
	}
	if req.Position != nil {
		user.Position = *req.Position
	}
	if req.HireDate != nil {
		t := req.HireDate.Time()
		user.HireDate = &t
	}
	if req.TerminationDate != nil {
		t := req.TerminationDate.Time()
		user.TerminationDate = &t
	}
}

func (m *UserMapper) ToStateDTOs(states []domain.UserState) []dto.UserStateResponse {
	if states == nil {
		return []dto.UserStateResponse{}
	}

	responses := make([]dto.UserStateResponse, len(states))
	for i, state := range states {
		responses[i] = dto.UserStateResponse{
			ID:          state.ID,
			StateID:     state.StateID,
			StateCode:   state.StateCode,
			StateName:   state.StateName,
			IsAvailable: state.IsAvailable,
			StartDate:   date.From(state.StartDate),
			EndDate:     date.From(state.EndDate),
		}
	}
	return responses
}

func datePtr(t *time.Time) *date.Date {
	if t == nil {
		return nil
	}
	d := date.From(*t)
	return &d
}

func timePtr(d *date.Date) *time.Time {
	if d == nil {
		return nil
	}
	t := d.Time()
	return &t
}
