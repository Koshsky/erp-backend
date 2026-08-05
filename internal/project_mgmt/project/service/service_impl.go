package service

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/domain"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/dto"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
)

type ProjectService struct {
	logger     *slog.Logger
	repository ProjectRepository
	mapper     *ProjectMapper
	validator  *ProjectValidator
}

func NewProjectService(logger *slog.Logger, repository ProjectRepository) *ProjectService {
	return &ProjectService{
		logger:     logger,
		repository: repository,
		mapper:     NewProjectMapper(),
		validator:  &ProjectValidator{},
	}
}

// canViewAllProjects возвращает true, если роль видит все проекты без привязки к owner_id.
func canViewAllProjects(role string) bool {
	return role == userdomain.Admin || role == userdomain.ProjectDirector
}

// isOwner проверяет, что текущий пользователь является владельцем проекта.
func isOwner(project *domain.Project, userID int64) bool {
	return project.OwnerID != nil && *project.OwnerID == userID
}

func (s *ProjectService) CreateProject(
	ctx context.Context,
	req dto.CreateProjectRequest,
	userID int64,
	role string,
) (*dto.ProjectResponse, error) {
	switch role {
	case userdomain.Admin:
		// admin может указать любого владельца (или оставить без owner)
	case userdomain.ProjectManager:
		// руководитель проекта сразу становится его owner, чужой owner из запроса игнорируется
		req.OwnerID = &userID
	default:
		return nil, ErrForbidden
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

func (s *ProjectService) FindProject(ctx context.Context, id int64, userID int64, role string) (*dto.ProjectResponse, error) {
	project, err := s.repository.FindProject(ctx, id)
	if err != nil {
		return nil, err
	}
	if project == nil || !s.canView(project, userID, role) {
		return nil, ErrNotFound
	}
	return s.mapper.ToDTO(project), nil
}

func (s *ProjectService) UpdateProject(
	ctx context.Context,
	id int64,
	req dto.UpdateProjectRequest,
	userID int64,
	role string,
) (*dto.ProjectResponse, error) {
	project, err := s.repository.FindProject(ctx, id)
	if err != nil || project == nil {
		return nil, ErrNotFound
	}

	switch role {
	case userdomain.Admin:
		// admin может менять все поля включая owner
	case userdomain.ProjectDirector:
		// директор проектов может менять только приоритет
		if req.OwnerID != nil || req.Code != nil || req.StartDate != nil || req.EndDate != nil {
			return nil, ErrForbidden
		}
	case userdomain.ProjectManager:
		// руководитель проекта редактирует только свои проекты и не может менять owner
		if !isOwner(project, userID) {
			return nil, ErrForbidden
		}
		if req.OwnerID != nil {
			return nil, ErrForbidden
		}
	default:
		return nil, ErrForbidden
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

func (s *ProjectService) DeleteProject(ctx context.Context, id int64, userID int64, role string) error {
	switch role {
	case userdomain.Admin:
		// admin может удалять любой проект
	case userdomain.ProjectManager:
		project, err := s.repository.FindProject(ctx, id)
		if err != nil || project == nil {
			return ErrNotFound
		}
		if !isOwner(project, userID) {
			return ErrForbidden
		}
	default:
		// dp и остальные роли удалять проекты не могут
		return ErrForbidden
	}

	return s.repository.DeleteProject(ctx, id)
}

func (s *ProjectService) ListProjects(ctx context.Context, userID int64, role string) ([]dto.ProjectResponse, error) {
	rows, err := s.repository.ListProjects(ctx)
	if err != nil {
		return nil, err
	}

	if canViewAllProjects(role) {
		return s.mapper.ToDTOs(rows), nil
	}

	filtered := make([]domain.Project, 0, len(rows))
	for _, project := range rows {
		if isOwner(&project, userID) {
			filtered = append(filtered, project)
		}
	}
	return s.mapper.ToDTOs(filtered), nil
}

// canView определяет, виден ли проект пользователю (dp/admin — все, остальные — только свои).
func (s *ProjectService) canView(project *domain.Project, userID int64, role string) bool {
	if canViewAllProjects(role) {
		return true
	}
	return isOwner(project, userID)
}
