package service

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/dto"
	"github.com/Koshsky/erp-backend/internal/validator"
)

type ProcessService struct {
	logger     *slog.Logger
	repository ProcessRepository
	mapper     *ProcessMapper
	validator  *ProcessValidator
}

func NewProcessService(logger *slog.Logger, repository ProcessRepository) *ProcessService {
	return &ProcessService{
		logger:     logger,
		repository: repository,
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
		if validator.IsNotFoundError(err) {
			return nil, validator.ErrProcessNotFound
		}
		return nil, err
	}
	if process == nil {
		return nil, validator.ErrProcessNotFound
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
		return nil, validator.ErrProcessNotFound
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
		return validator.ErrProcessNotFound
	}

	return s.repository.DeleteProcess(ctx, id)
}

func (s *ProcessService) ListProcesses(ctx context.Context) ([]dto.ProcessResponse, error) {
	rows, err := s.repository.ListProcesss(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(rows), nil
}
