package service

import (
	"context"
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/timesheet/state/repository"

	"github.com/Koshsky/erp-backend/internal/timesheet/state/dto"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type StateService struct {
	logger     *slog.Logger
	repository StateRepository
	mapper     *StateMapper
	validator  *StateValidator
	tracer     *tracingpkg.Tracer
}

// NewStateService builds the StateService service.
func NewStateService(logger *slog.Logger, tracer *tracingpkg.Tracer, r *repo.StateRepository) *StateService {
	return &StateService{
		logger:     logger,
		repository: r,
		mapper:     NewStateMapper(),
		validator:  &StateValidator{},
		tracer:     tracer,
	}
}

func (s *StateService) ListStates(ctx context.Context) ([]dto.StateResponse, error) {
	ctx, end := s.tracer.Start(ctx, "state.ListStates")
	defer end(nil)
	states, err := s.repository.ListStates(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(states), nil
}

func (s *StateService) CreateState(ctx context.Context, req dto.CreateStateRequest) (*dto.StateResponse, error) {
	ctx, end := s.tracer.Start(ctx, "state.CreateState")
	defer end(nil)

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
	ctx, end := s.tracer.Start(ctx, "state.FindState")
	defer end(nil)

	state, err := s.repository.FindState(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil, errors.ErrStateNotFound
		}
		return nil, err
	}
	if state == nil {
		return nil, errors.ErrStateNotFound
	}
	return s.mapper.ToDTO(state), nil
}

func (s *StateService) UpdateState(
	ctx context.Context,
	id int64,
	req dto.UpdateStateRequest,
) (*dto.StateResponse, error) {
	ctx, end := s.tracer.Start(ctx, "state.UpdateState")
	defer end(nil)

	state, err := s.repository.FindState(ctx, id)
	if err != nil || state == nil {
		return nil, errors.ErrStateNotFound
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
	ctx, end := s.tracer.Start(ctx, "state.DeleteState")
	defer end(nil)

	state, err := s.repository.FindState(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil // идемпотентный delete: уже удалено — не ошибка
		}
		return err
	}
	if state == nil {
		return nil // идемпотентный delete
	}

	return s.repository.DeleteState(ctx, id)
}
