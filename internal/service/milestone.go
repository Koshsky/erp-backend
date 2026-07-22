package service

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/dto"
	"github.com/Koshsky/erp/api/internal/service/mapper"
)

type MilestoneRepository interface {
	CreateMilestone(ctx context.Context, milestone domain.Milestone) (*domain.Milestone, error)
	GetMilestone(ctx context.Context, id int64) (*domain.Milestone, error)
	UpdateMilestone(ctx context.Context, milestone domain.Milestone) (*domain.Milestone, error)
	DeleteMilestone(ctx context.Context, id int64) error
	ListMilestonesByProcessID(ctx context.Context, processID int64) ([]domain.Milestone, error)
}

type MilestoneService struct {
	logger     *slog.Logger
	repository MilestoneRepository
	mapper     *mapper.MilestoneMapper
	validator  *Validator
}

func NewMilestoneService(logger *slog.Logger, repository MilestoneRepository, validator *Validator) *MilestoneService {
	return &MilestoneService{
		logger:     logger,
		repository: repository,
		mapper:     mapper.NewMilestoneMapper(),
		validator:  validator,
	}
}

func (s *MilestoneService) CreateMilestone(ctx context.Context, req dto.CreateMilestoneRequest) (*dto.MilestoneResponse, error) {
	if err := s.validator.ValidateMilestone(req.ProcessID, req.Title, req.Content, req.Date); err != nil {
		return nil, err
	}
	created, err := s.repository.CreateMilestone(ctx, s.mapper.ToDomainFromCreate(req))
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(created), nil
}

func (s *MilestoneService) GetMilestone(ctx context.Context, id int64) (*dto.MilestoneResponse, error) {
	milestone, err := s.repository.GetMilestone(ctx, id)
	if err != nil {
		return nil, err
	}
	if milestone == nil {
		return nil, ErrMilestoneNotFound
	}

	return s.mapper.ToDTO(milestone), nil
}

func (s *MilestoneService) UpdateMilestone(ctx context.Context, id int64, req dto.UpdateMilestoneRequest) (*dto.MilestoneResponse, error) {
	milestone, err := s.repository.GetMilestone(ctx, id)
	if err != nil {
		return nil, err
	}
	if milestone == nil {
		return nil, ErrMilestoneNotFound
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

func (s *MilestoneService) ListMilestonesByProcessID(ctx context.Context, processID int64) ([]dto.MilestoneResponse, error) {
	milestones, err := s.repository.ListMilestonesByProcessID(ctx, processID)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTOs(milestones), nil
}
