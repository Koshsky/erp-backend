package service

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/dto"
	"github.com/Koshsky/erp/api/internal/service/mapper"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, project domain.Project) (*domain.Project, error)
	GetProject(ctx context.Context, id int64) (*domain.Project, error)
	UpdateProject(ctx context.Context, project domain.Project) (*domain.Project, error)
	DeleteProject(ctx context.Context, id int64) error
	ListProjects(ctx context.Context) ([]domain.Project, error)
	GetDetailedProject(ctx context.Context, id int64) (*domain.DetailedProject, error)
}

type ProjectService struct {
	logger     *slog.Logger
	repository ProjectRepository
	mapper     *mapper.ProjectMapper
	validator  *Validator
}

func NewProjectService(logger *slog.Logger, repository ProjectRepository, validator *Validator) *ProjectService {
	return &ProjectService{
		logger:     logger,
		repository: repository,
		mapper:     mapper.NewProjectMapper(),
		validator:  validator,
	}
}

func (s *ProjectService) GetDetailedProject(ctx context.Context, id int64) (*dto.ProjectDetailResponse, error) {
	project, err := s.repository.GetDetailedProject(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDetailedDTO(project), nil
}

func (s *ProjectService) CreateProject(ctx context.Context, req dto.CreateProjectRequest) (*dto.ProjectResponse, error) {
	if err := s.validator.ValidateProject(req.Code, req.StartDate, req.EndDate, req.Priority); err != nil {
		return nil, err
	}
	project := s.mapper.ToDomainFromCreate(req)
	createdProject, err := s.repository.CreateProject(ctx, project)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(createdProject), nil
}

func (s *ProjectService) GetProject(ctx context.Context, id int64) (*dto.ProjectResponse, error) {
	project, err := s.repository.GetProject(ctx, id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}
	return s.mapper.ToDTO(project), nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, id int64, req dto.UpdateProjectRequest) (*dto.ProjectResponse, error) {
	project, err := s.repository.GetProject(ctx, id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}
	s.mapper.ApplyUpdateToDomain(project, req)
	if err := s.validator.ValidateProject(project.Code, project.StartDate, project.EndDate, project.Priority); err != nil {
		return nil, err
	}
	updatedProject, err := s.repository.UpdateProject(ctx, *project)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(updatedProject), nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, id int64) error {
	return s.repository.DeleteProject(ctx, id)
}

func (s *ProjectService) ListProjects(ctx context.Context) ([]dto.ProjectResponse, error) {
	projects, err := s.repository.ListProjects(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(projects), nil
}
