package service

import (
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/project_mgmt/project/repository"
)

// ProvideProjectService builds the ProjectService service.
func ProvideProjectService(logger *slog.Logger, r *repo.ProjectRepository) *ProjectService {
	return &ProjectService{
		logger:     logger,
		repository: r,
		mapper:     NewProjectMapper(),
		validator:  &ProjectValidator{},
	}
}
