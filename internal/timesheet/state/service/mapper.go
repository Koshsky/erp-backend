package service

import (
	"github.com/Koshsky/erp-backend/internal/timesheet/state/domain"
	"github.com/Koshsky/erp-backend/internal/timesheet/state/dto"
)

type StateMapper struct{}

func NewStateMapper() *StateMapper {
	return &StateMapper{}
}

func (m *StateMapper) ToDTO(state *domain.State) *dto.StateResponse {
	if state == nil {
		return nil
	}
	return &dto.StateResponse{
		ID:          state.ID,
		Code:        state.Code,
		Name:        state.Name,
		IsAvailable: state.IsAvailable,
	}
}

func (m *StateMapper) ToDTOs(states []domain.State) []dto.StateResponse {
	if states == nil {
		return []dto.StateResponse{}
	}

	responses := make([]dto.StateResponse, len(states))
	for i, state := range states {
		responses[i] = *m.ToDTO(&state)
	}
	return responses
}

func (m *StateMapper) ToDomainFromCreate(req dto.CreateStateRequest) domain.State {
	return domain.State{
		Code:        req.Code,
		Name:        req.Name,
		IsAvailable: req.IsAvailable,
	}
}

func (m *StateMapper) ApplyUpdateToDomain(state *domain.State, req dto.UpdateStateRequest) {
	if state == nil {
		return
	}

	if req.Code != nil {
		state.Code = *req.Code
	}
	if req.Name != nil {
		state.Name = *req.Name
	}
	if req.IsAvailable != nil {
		state.IsAvailable = *req.IsAvailable
	}
}
