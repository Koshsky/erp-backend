package service

import (
	"context"
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/project_mgmt/project/repository"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/dto"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type ProjectService struct {
	logger     *slog.Logger
	repository ProjectRepository
	mapper     *ProjectMapper
	validator  *ProjectValidator
}

// NewProjectService builds the ProjectService service.
func NewProjectService(logger *slog.Logger, r *repo.ProjectRepository) *ProjectService {
	return &ProjectService{
		logger:     logger,
		repository: r,
		mapper:     NewProjectMapper(),
		validator:  &ProjectValidator{},
	}
}

// CreateProject creates a project. The middleware checked permissions; here
// only owner normalization: rp always becomes the owner.
func (s *ProjectService) CreateProject(
	ctx context.Context,
	req dto.CreateProjectRequest,
	userID int64,
	role string,
) (*dto.ProjectResponse, error) {
	if role == userdomain.ProjectManager {
		// the project manager immediately becomes the owner; a foreign owner from the request is ignored
		req.OwnerID = &userID
	}

	project := s.mapper.ToDomainFromCreate(req)
	if err := s.validator.ValidateProject(&project); err != nil {
		return nil, err
	}

	created, err := s.repository.CreateProject(ctx, project)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(created), nil
}

func (s *ProjectService) FindProject(ctx context.Context, id int64) (*dto.ProjectResponse, error) {
	project, err := s.repository.FindProject(ctx, id)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.ErrProjectNotFound
	}
	return s.mapper.ToDTO(project), nil
}

// UpdateProject changes code, dates, priority and owner. Permissions (priority
// vs other fields, owner change admin-only) are checked by the middleware
// against the request body.
func (s *ProjectService) UpdateProject(
	ctx context.Context,
	id int64,
	req dto.UpdateProjectRequest,
) (*dto.ProjectResponse, error) {
	project, err := s.repository.FindProject(ctx, id)
	if err != nil || project == nil {
		return nil, errors.ErrProjectNotFound
	}

	s.mapper.ApplyUpdateToDomain(project, req)
	if err = s.validator.ValidateProject(project); err != nil {
		return nil, err
	}

	updated, err := s.repository.UpdateProject(ctx, *project)
	if err != nil {
		return nil, err
	}

	return s.mapper.ToDTO(updated), nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, id int64) error {
	project, err := s.repository.FindProject(ctx, id)
	if err != nil || project == nil {
		return errors.ErrProjectNotFound
	}

	return s.repository.DeleteProject(ctx, id)
}

func (s *ProjectService) ListProjects(ctx context.Context, limit, offset int) ([]dto.ProjectResponse, int64, error) {
	rows, err := s.repository.ListProjects(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repository.CountProjects(ctx)
	if err != nil {
		return nil, 0, err
	}
	return s.mapper.ToDTOs(rows), total, nil
}
