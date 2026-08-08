package service

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/dto"
	"github.com/Koshsky/erp-backend/internal/validator"
)

type AssignmentService struct {
	logger     *slog.Logger
	repository AssignmentRepository
	mapper     *AssignmentMapper
	validator  *AssignmentValidator
}

func NewAssignmentService(logger *slog.Logger, repository AssignmentRepository) *AssignmentService {
	return &AssignmentService{
		logger:     logger,
		repository: repository,
		mapper:     NewAssignmentMapper(),
		validator:  &AssignmentValidator{},
	}
}

func (s *AssignmentService) CreateAssignment(
	ctx context.Context,
	req dto.CreateAssignmentRequest,
) (*dto.AssignmentResponse, error) {
	assignment := s.mapper.ToDomainFromCreate(req)
	if err := s.validator.ValidateAssignment(&assignment); err != nil {
		return nil, err
	}

	created, err := s.repository.CreateAssignment(ctx, assignment)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(created), nil
}

func (s *AssignmentService) FindAssignment(ctx context.Context, id int64) (*dto.AssignmentResponse, error) {
	assignment, err := s.repository.FindAssignment(ctx, id)
	if err != nil {
		if validator.IsNotFoundError(err) {
			return nil, validator.ErrAssignmentNotFound
		}
		return nil, err
	}
	if assignment == nil {
		return nil, validator.ErrAssignmentNotFound
	}
	return s.mapper.ToDTO(assignment), nil
}

func (s *AssignmentService) UpdateAssignment(
	ctx context.Context,
	id int64,
	req dto.UpdateAssignmentRequest,
) (*dto.AssignmentResponse, error) {
	assignment, err := s.repository.FindAssignment(ctx, id)
	if err != nil || assignment == nil {
		return nil, validator.ErrAssignmentNotFound
	}

	s.mapper.ApplyUpdateToDomain(assignment, req)
	if err = s.validator.ValidateAssignment(assignment); err != nil {
		return nil, err
	}

	updated, err := s.repository.UpdateAssignment(ctx, *assignment)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(updated), nil
}

func (s *AssignmentService) DeleteAssignment(ctx context.Context, id int64) error {
	assignment, err := s.repository.FindAssignment(ctx, id)
	if err != nil || assignment == nil {
		return validator.ErrAssignmentNotFound
	}

	return s.repository.DeleteAssignment(ctx, id)
}

func (s *AssignmentService) ListAssignments(ctx context.Context) ([]dto.AssignmentResponse, error) {
	rows, err := s.repository.ListAssignments(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(rows), nil
}
