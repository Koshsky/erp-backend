package service

import (
	"context"
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/project_mgmt/process/repository"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/dto"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type ProcessService struct {
	logger     *slog.Logger
	tracer     *tracingpkg.Tracer
	repository ProcessRepository
	mapper     *ProcessMapper
	validator  *ProcessValidator
}

// NewProcessService builds the ProcessService service.
func NewProcessService(logger *slog.Logger, tracer *tracingpkg.Tracer, r *repo.ProcessRepository) *ProcessService {
	return &ProcessService{
		logger:     logger,
		tracer:     tracer,
		repository: r,
		mapper:     NewProcessMapper(),
		validator:  &ProcessValidator{},
	}
}

func (s *ProcessService) CreateProcess(
	ctx context.Context,
	req dto.CreateProcessRequest,
) (*dto.ProcessResponse, error) {
	ctx, end := s.tracer.Start(ctx, "process.CreateProcess")
	defer end(nil)

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
	ctx, end := s.tracer.Start(ctx, "process.FindProcess")
	defer end(nil)

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
	ctx, end := s.tracer.Start(ctx, "process.UpdateProcess")
	defer end(nil)

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
	ctx, end := s.tracer.Start(ctx, "process.DeleteProcess")
	defer end(nil)

	process, err := s.repository.FindProcess(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil // идемпотентный delete: уже удалено — не ошибка
		}
		return err
	}
	if process == nil {
		return nil // идемпотентный delete
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
	ctx, end := s.tracer.Start(ctx, "process.ListProcesses")
	defer end(nil)

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
