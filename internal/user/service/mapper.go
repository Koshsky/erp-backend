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
		Name:            user.FullName(),
		LastName:        user.LastName,
		FirstName:       user.FirstName,
		MiddleName:      user.MiddleName,
		Username:        user.Username,
		Preset:          user.Preset,
		ManagerID:       user.ManagerID,
		Position:        user.Position,
		HireDate:        datePtr(user.HireDate),
		TerminationDate: datePtr(user.TerminationDate),
		PasswordHash:    user.PasswordHash,
		CreatedAt:       user.CreatedAt,
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
		LastName:        req.LastName,
		FirstName:       req.FirstName,
		MiddleName:      normalizeMiddle(req.MiddleName),
		Username:        req.Username,
		Preset:          req.Preset,
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

	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.MiddleName != nil {
		user.MiddleName = normalizeMiddle(req.MiddleName)
	}
	if req.Username != nil {
		user.Username = *req.Username
	}
	if req.Preset != nil {
		if *req.Preset == "" {
			user.Preset = nil
		} else {
			user.Preset = req.Preset
		}
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

// normalizeMiddle converts an empty middle-name string to nil (no middle name).
func normalizeMiddle(p *string) *string {
	if p != nil && *p == "" {
		return nil
	}
	return p
}
