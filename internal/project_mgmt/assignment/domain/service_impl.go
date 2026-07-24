package domain

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/dto"
)

type AssignmentService struct {
	logger     *slog.Logger
	repository RepositoryInterface
	mapper     *AssignmentMapper
	validator  *AssignmentValidator
}

func NewAssignmentService(logger *slog.Logger, repository RepositoryInterface) *AssignmentService {
	return &AssignmentService{
		logger:     logger,
		repository: repository,
		mapper:     NewAssignmentMapper(),
		validator:  &AssignmentValidator{},
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
		return nil, fmt.Errorf("assignment not found")
	}
	return s.mapper.ToDTO(assignment), nil
}

func (s *AssignmentService) UpdateAssignment(ctx context.Context, id int64, req dto.UpdateAssignmentRequest) (*dto.AssignmentResponse, error) {
	assignment, err := s.repository.GetAssignment(ctx, id)
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, fmt.Errorf("assignment not found")
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

func (s *AssignmentService) ListAssignments(ctx context.Context) ([]dto.AssignmentResponse, error) {
	rows, err := s.repository.ListAssignments(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(rows), nil
}
