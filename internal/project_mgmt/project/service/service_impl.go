package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/dto"
)

type ProjectService struct {
	logger     *slog.Logger
	repository RepositoryInterface
	mapper     *ProjectMapper
	validator  *ProjectValidator
}

func NewProjectService(logger *slog.Logger, repository RepositoryInterface) *ProjectService {
	return &ProjectService{
		logger:     logger,
		repository: repository,
		mapper:     NewProjectMapper(),
		validator:  &ProjectValidator{},
	}
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
		return nil, fmt.Errorf("project not found")
	}
	return s.mapper.ToDTO(project), nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, id int64, req dto.UpdateProjectRequest) (*dto.ProjectResponse, error) {
	project, err := s.repository.GetProject(ctx, id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
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
	rows, err := s.repository.ListProjects(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(rows), nil
}
