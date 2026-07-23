package service

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/dto"
	"github.com/Koshsky/erp/api/internal/service/mapper"
)

type ProcessRepository interface {
	CreateProcess(ctx context.Context, process domain.Process) (*domain.Process, error)
	GetProcess(ctx context.Context, id int64) (*domain.Process, error)
	UpdateProcess(ctx context.Context, process domain.Process) (*domain.Process, error)
	DeleteProcess(ctx context.Context, id int64) error
	ListProcesses(ctx context.Context) ([]domain.Process, error)
}

type ProcessService struct {
	logger     *slog.Logger
	repository ProcessRepository
	mapper     *mapper.ProcessMapper
	validator  *Validator
}

func NewProcessService(logger *slog.Logger, repository ProcessRepository, validator *Validator) *ProcessService {
	return &ProcessService{
		logger:     logger,
		repository: repository,
		mapper:     mapper.NewProcessMapper(),
		validator:  validator,
	}
}

func (s *ProcessService) CreateProcess(ctx context.Context, req dto.CreateProcessRequest) (*dto.ProcessResponse, error) {
	if err := s.validator.ValidateProcess(req.ProjectID, req.Title, req.StartDate, req.EndDate); err != nil {
		return nil, err
	}
	created, err := s.repository.CreateProcess(ctx, s.mapper.ToDomainFromCreate(req))
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(created), nil
}

func (s *ProcessService) GetProcess(ctx context.Context, id int64) (*dto.ProcessResponse, error) {
	process, err := s.repository.GetProcess(ctx, id)
	if err != nil {
		return nil, err
	}
	if process == nil {
		return nil, ErrProcessNotFound
	}
	return s.mapper.ToDTO(process), nil
}

func (s *ProcessService) UpdateProcess(ctx context.Context, id int64, req dto.UpdateProcessRequest) (*dto.ProcessResponse, error) {
	process, err := s.repository.GetProcess(ctx, id)
	if err != nil {
		return nil, err
	}
	if process == nil {
		return nil, ErrProcessNotFound
	}
	s.mapper.ApplyUpdateToDomain(process, req)
	if err := s.validator.ValidateProcess(process.ProjectID, process.Title, process.StartDate, process.EndDate); err != nil {
		return nil, err
	}
	updated, err := s.repository.UpdateProcess(ctx, *process)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(updated), nil
}

func (s *ProcessService) DeleteProcess(ctx context.Context, id int64) error {
	return s.repository.DeleteProcess(ctx, id)
}

func (s *ProcessService) ListProcesses(ctx context.Context, projectID int64) ([]dto.ProcessResponse, error) {
	domainProcesses, err := s.repository.ListProcesses(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(domainProcesses), nil
}
