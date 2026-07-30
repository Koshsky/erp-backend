package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/dto"
)

type MilestoneService struct {
	logger     *slog.Logger
	repository MilestoneRepository
	mapper     *MilestoneMapper
	validator  *MilestoneValidator
}

func NewMilestoneService(logger *slog.Logger, repository MilestoneRepository) *MilestoneService {
	return &MilestoneService{
		logger:     logger,
		repository: repository,
		mapper:     NewMilestoneMapper(),
		validator:  &MilestoneValidator{},
	}
}

func (s *MilestoneService) CreateMilestone(
	ctx context.Context,
	req dto.CreateMilestoneRequest,
) (*dto.MilestoneResponse, error) {
	if err := s.validator.ValidateMilestone(req.ProcessID, req.Title, req.Content, req.Date); err != nil {
		return nil, err
	}
	created, err := s.repository.CreateMilestone(ctx, s.mapper.ToDomainFromCreate(req))
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(created), nil
}

func (s *MilestoneService) FindMilestone(ctx context.Context, id int64) (*dto.MilestoneResponse, error) {
	milestone, err := s.repository.FindMilestone(ctx, id)
	if err != nil {
		return nil, err
	}
	if milestone == nil {
		return nil, fmt.Errorf("milestone not found")
	}

	return s.mapper.ToDTO(milestone), nil
}

func (s *MilestoneService) UpdateMilestone(
	ctx context.Context,
	id int64,
	req dto.UpdateMilestoneRequest,
) (*dto.MilestoneResponse, error) {
	milestone, err := s.repository.FindMilestone(ctx, id)
	if err != nil {
		return nil, err
	}
	if milestone == nil {
		return nil, fmt.Errorf("milestone not found")
	}

	s.mapper.ApplyUpdateToDomain(milestone, req)
	if err := s.validator.ValidateMilestone(milestone.ProcessID, milestone.Title, milestone.Content, milestone.Date); err != nil {
		return nil, err
	}
	updated, err := s.repository.UpdateMilestone(ctx, *milestone)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(updated), nil
}

func (s *MilestoneService) DeleteMilestone(ctx context.Context, id int64) error {
	return s.repository.DeleteMilestone(ctx, id)
}

func (s *MilestoneService) ListMilestones(ctx context.Context) ([]dto.MilestoneResponse, error) {
	rows, err := s.repository.ListMilestones(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(rows), nil
}
