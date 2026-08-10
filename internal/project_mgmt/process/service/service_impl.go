package service

import (
	"context"
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/project_mgmt/process/repository"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/dto"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type ProcessService struct {
	logger     *slog.Logger
	repository ProcessRepository
	mapper     *ProcessMapper
	validator  *ProcessValidator
}

// NewProcessService builds the ProcessService service.
func NewProcessService(logger *slog.Logger, r *repo.ProcessRepository) *ProcessService {
	return &ProcessService{
		logger:     logger,
		repository: r,
		mapper:     NewProcessMapper(),
		validator:  &ProcessValidator{},
	}
}

func (s *ProcessService) CreateProcess(
	ctx context.Context,
	req dto.CreateProcessRequest,
) (*dto.ProcessResponse, error) {
	process := s.mapper.ToDomainFromCreate(req)
	if err := s.validator.ValidateProcess(&process); err != nil {
		return nil, err
	}

	created, err := s.repository.CreateProcess(ctx, process)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(created), nil
}

func (s *ProcessService) FindProcess(ctx context.Context, id int64) (*dto.ProcessResponse, error) {
	process, err := s.repository.FindProcess(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil, errors.ErrProcessNotFound
		}
		return nil, err
	}
	if process == nil {
		return nil, errors.ErrProcessNotFound
	}
	return s.mapper.ToDTO(process), nil
}

func (s *ProcessService) UpdateProcess(
	ctx context.Context,
	id int64,
	req dto.UpdateProcessRequest,
) (*dto.ProcessResponse, error) {
	process, err := s.repository.FindProcess(ctx, id)
	if err != nil || process == nil {
		return nil, errors.ErrProcessNotFound
	}

	s.mapper.ApplyUpdateToDomain(process, req)
	if err = s.validator.ValidateProcess(process); err != nil {
		return nil, err
	}

	updated, err := s.repository.UpdateProcess(ctx, *process)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(updated), nil
}

func (s *ProcessService) DeleteProcess(ctx context.Context, id int64) error {
	process, err := s.repository.FindProcess(ctx, id)
	if err != nil || process == nil {
		return errors.ErrProcessNotFound
	}

	return s.repository.DeleteProcess(ctx, id)
}

func (s *ProcessService) ListProcesses(
	ctx context.Context,
	userID int64,
	role string,
	ownerID int64,
	limit, offset int,
) ([]dto.ProcessResponse, int64, error) {
	rows, err := s.repository.ListProcesss(ctx, userID, role, ownerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repository.CountProcesses(ctx, userID, role, ownerID)
	if err != nil {
		return nil, 0, err
	}
	return s.mapper.ToDTOs(rows), total, nil
}
