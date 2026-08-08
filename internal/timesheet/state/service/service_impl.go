package service

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/timesheet/state/dto"
	"github.com/Koshsky/erp-backend/internal/validator"
)

type StateService struct {
	logger     *slog.Logger
	repository StateRepository
	mapper     *StateMapper
	validator  *StateValidator
}

func NewStateService(logger *slog.Logger, repository StateRepository) *StateService {
	return &StateService{
		logger:     logger,
		repository: repository,
		mapper:     NewStateMapper(),
		validator:  &StateValidator{},
	}
}

func (s *StateService) ListStates(ctx context.Context) ([]dto.StateResponse, error) {
	states, err := s.repository.ListStates(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(states), nil
}

func (s *StateService) CreateState(ctx context.Context, req dto.CreateStateRequest) (*dto.StateResponse, error) {
	state := s.mapper.ToDomainFromCreate(req)
	if err := s.validator.ValidateState(&state); err != nil {
		return nil, err
	}

	created, err := s.repository.CreateState(ctx, state)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(created), nil
}

func (s *StateService) FindState(ctx context.Context, id int64) (*dto.StateResponse, error) {
	state, err := s.repository.FindState(ctx, id)
	if err != nil {
		if validator.IsNotFoundError(err) {
			return nil, validator.ErrStateNotFound
		}
		return nil, err
	}
	if state == nil {
		return nil, validator.ErrStateNotFound
	}
	return s.mapper.ToDTO(state), nil
}

func (s *StateService) UpdateState(
	ctx context.Context,
	id int64,
	req dto.UpdateStateRequest,
) (*dto.StateResponse, error) {
	state, err := s.repository.FindState(ctx, id)
	if err != nil || state == nil {
		return nil, validator.ErrStateNotFound
	}

	s.mapper.ApplyUpdateToDomain(state, req)
	if err = s.validator.ValidateState(state); err != nil {
		return nil, err
	}

	updated, err := s.repository.UpdateState(ctx, *state)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(updated), nil
}

func (s *StateService) DeleteState(ctx context.Context, id int64) error {
	state, err := s.repository.FindState(ctx, id)
	if err != nil || state == nil {
		return validator.ErrStateNotFound
	}

	return s.repository.DeleteState(ctx, id)
}
