package service

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/dto"
	"github.com/Koshsky/erp/api/internal/service/mapper"
)

type AssignmentRepository interface {
	CreateAssignment(ctx context.Context, assignment domain.Assignment) (*domain.Assignment, error)
	GetAssignment(ctx context.Context, id int64) (*domain.Assignment, error)
	UpdateAssignment(ctx context.Context, assignment domain.Assignment) (*domain.Assignment, error)
	DeleteAssignment(ctx context.Context, id int64) error
}

type AssignmentService struct {
	logger     *slog.Logger
	repository AssignmentRepository
	mapper     *mapper.AssignmentMapper
	validator  *Validator
}

func NewAssignmentService(logger *slog.Logger, repository AssignmentRepository, validator *Validator) *AssignmentService {
	return &AssignmentService{
		logger:     logger,
		repository: repository,
		mapper:     mapper.NewAssignmentMapper(),
		validator:  validator,
	}
}

func (s *AssignmentService) CreateAssignment(ctx context.Context, req dto.CreateAssignmentRequest) (*dto.AssignmentResponse, error) {
	if err := s.validator.ValidateAssignment(req.TaskID, req.ResourceID, req.Quantity); err != nil {
		return nil, err
	}
	created, err := s.repository.CreateAssignment(ctx, s.mapper.ToDomainFromCreate(req))
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(created), nil
}

func (s *AssignmentService) GetAssignment(ctx context.Context, id int64) (*dto.AssignmentResponse, error) {
	assignment, err := s.repository.GetAssignment(ctx, id)
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, ErrAssignmentNotFound
	}
	return s.mapper.ToDTO(assignment), nil
}

func (s *AssignmentService) UpdateAssignment(ctx context.Context, id int64, req dto.UpdateAssignmentRequest) (*dto.AssignmentResponse, error) {
	assignment, err := s.repository.GetAssignment(ctx, id)
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, ErrAssignmentNotFound
	}

	s.mapper.ApplyUpdateToDomain(assignment, req)
	if err := s.validator.ValidateAssignment(assignment.TaskID, assignment.ResourceID, assignment.Quantity); err != nil {
		return nil, err
	}

	updated, err := s.repository.UpdateAssignment(ctx, *assignment)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(updated), nil
}

func (s *AssignmentService) DeleteAssignment(ctx context.Context, id int64) error {
	return s.repository.DeleteAssignment(ctx, id)
}
