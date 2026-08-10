package service

import (
	"context"
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/repository"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/dto"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type MilestoneService struct {
	logger     *slog.Logger
	repository MilestoneRepository
	mapper     *MilestoneMapper
	validator  *MilestoneValidator
}

// NewMilestoneService builds the MilestoneService service.
func NewMilestoneService(logger *slog.Logger, r *repo.MilestoneRepository) *MilestoneService {
	return &MilestoneService{
		logger:     logger,
		repository: r,
		mapper:     NewMilestoneMapper(),
		validator:  &MilestoneValidator{},
	}
}

func (s *MilestoneService) CreateMilestone(
	ctx context.Context,
	req dto.CreateMilestoneRequest,
) (*dto.MilestoneResponse, error) {
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
	milestone, err := s.repository.FindMilestone(ctx, id)
	if err != nil || milestone == nil {
		return errors.ErrMilestoneNotFound
	}

	return s.repository.DeleteMilestone(ctx, id)
}

func (s *MilestoneService) ListMilestones(ctx context.Context) ([]dto.MilestoneResponse, error) {
	rows, err := s.repository.ListMilestones(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(rows), nil
}
