package service

import (
	"context"
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/repository"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/dto"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type MilestoneService struct {
	logger     *slog.Logger
	tracer     *tracingpkg.Tracer
	repository MilestoneRepository
	mapper     *MilestoneMapper
	validator  *MilestoneValidator
}

// NewMilestoneService builds the MilestoneService service.
func NewMilestoneService(
	logger *slog.Logger,
	tracer *tracingpkg.Tracer,
	r *repo.MilestoneRepository,
) *MilestoneService {
	return &MilestoneService{
		logger:     logger,
		tracer:     tracer,
		repository: r,
		mapper:     NewMilestoneMapper(),
		validator:  &MilestoneValidator{},
	}
}

func (s *MilestoneService) CreateMilestone(
	ctx context.Context,
	req dto.CreateMilestoneRequest,
) (*dto.MilestoneResponse, error) {
	ctx, end := s.tracer.Start(ctx, "milestone.CreateMilestone")
	defer end(nil)

	milestone := s.mapper.ToDomainFromCreate(req)
	if err := s.validator.ValidateMilestone(&milestone); err != nil {
		return nil, err
	}

	created, err := s.repository.CreateMilestone(ctx, milestone)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(created), nil
}

func (s *MilestoneService) FindMilestone(ctx context.Context, id int64) (*dto.MilestoneResponse, error) {
	ctx, end := s.tracer.Start(ctx, "milestone.FindMilestone")
	defer end(nil)

	milestone, err := s.repository.FindMilestone(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil, errors.ErrMilestoneNotFound
		}
		return nil, err
	}
	if milestone == nil {
		return nil, errors.ErrMilestoneNotFound
	}
	return s.mapper.ToDTO(milestone), nil
}

func (s *MilestoneService) UpdateMilestone(
	ctx context.Context,
	id int64,
	req dto.UpdateMilestoneRequest,
) (*dto.MilestoneResponse, error) {
	ctx, end := s.tracer.Start(ctx, "milestone.UpdateMilestone")
	defer end(nil)

	milestone, err := s.repository.FindMilestone(ctx, id)
	if err != nil || milestone == nil {
		return nil, errors.ErrMilestoneNotFound
	}

	s.mapper.ApplyUpdateToDomain(milestone, req)
	if err = s.validator.ValidateMilestone(milestone); err != nil {
		return nil, err
	}

	updated, err := s.repository.UpdateMilestone(ctx, *milestone)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(updated), nil
}

func (s *MilestoneService) DeleteMilestone(ctx context.Context, id int64) error {
	ctx, end := s.tracer.Start(ctx, "milestone.DeleteMilestone")
	defer end(nil)

	milestone, err := s.repository.FindMilestone(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil // идемпотентный delete: уже удалено — не ошибка
		}
		return err
	}
	if milestone == nil {
		return nil // идемпотентный delete
	}

	return s.repository.DeleteMilestone(ctx, id)
}

func (s *MilestoneService) ListMilestones(
	ctx context.Context,
	userID int64,
	role string,
	ownerID int64,
	limit, offset int,
) ([]dto.MilestoneResponse, int64, error) {
	ctx, end := s.tracer.Start(ctx, "milestone.ListMilestones")
	defer end(nil)

	rows, err := s.repository.ListMilestones(ctx, userID, role, ownerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repository.CountMilestones(ctx, userID, role, ownerID)
	if err != nil {
		return nil, 0, err
	}
	return s.mapper.ToDTOs(rows), total, nil
}
