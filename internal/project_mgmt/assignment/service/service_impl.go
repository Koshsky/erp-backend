package service

import (
	"context"
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/repository"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/dto"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type AssignmentService struct {
	logger     *slog.Logger
	repository AssignmentRepository
	mapper     *AssignmentMapper
	validator  *AssignmentValidator
}

// NewAssignmentService builds the AssignmentService service.
func NewAssignmentService(logger *slog.Logger, r *repo.AssignmentRepository) *AssignmentService {
	return &AssignmentService{
		logger:     logger,
		repository: r,
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
		if errors.IsNotFoundError(err) {
			return nil, errors.ErrAssignmentNotFound
		}
		return nil, err
	}
	if assignment == nil {
		return nil, errors.ErrAssignmentNotFound
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
		return nil, errors.ErrAssignmentNotFound
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
		return errors.ErrAssignmentNotFound
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
