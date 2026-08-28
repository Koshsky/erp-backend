package service

import (
	"context"
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/project_mgmt/project/repository"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/dto"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type ProjectService struct {
	logger     *slog.Logger
	tracer     *tracingpkg.Tracer
	repository ProjectRepository
	mapper     *ProjectMapper
	validator  *ProjectValidator
}

// NewProjectService builds the ProjectService service.
func NewProjectService(logger *slog.Logger, tracer *tracingpkg.Tracer, r *repo.ProjectRepository) *ProjectService {
	return &ProjectService{
		logger:     logger,
		tracer:     tracer,
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
	ctx, end := s.tracer.Start(ctx, "project.CreateProject")
	defer end(nil)

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
	ctx, end := s.tracer.Start(ctx, "project.FindProject")
	defer end(nil)

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
	ctx, end := s.tracer.Start(ctx, "project.UpdateProject")
	defer end(nil)

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
	ctx, end := s.tracer.Start(ctx, "project.DeleteProject")
	defer end(nil)

	project, err := s.repository.FindProject(ctx, id)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return nil // idempotent delete: already deleted — not an error
		}
		return err
	}
	if project == nil {
		return nil // idempotent delete
	}

	return s.repository.DeleteProject(ctx, id)
}

func (s *ProjectService) ListProjects(
	ctx context.Context,
	userID int64,
	role string,
	ownerID int64,
	limit, offset int,
) ([]dto.ProjectResponse, int64, error) {
	ctx, end := s.tracer.Start(ctx, "project.ListProjects")
	defer end(nil)

	rows, err := s.repository.ListProjects(ctx, userID, role, ownerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repository.CountProjects(ctx, userID, role, ownerID)
	if err != nil {
		return nil, 0, err
	}
	return s.mapper.ToDTOs(rows), total, nil
}
