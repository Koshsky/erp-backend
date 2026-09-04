package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/domain"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, Project domain.Project) (*domain.Project, error)
	FindProject(ctx context.Context, id int64) (*domain.Project, error)
	UpdateProject(ctx context.Context, project domain.Project) (*domain.Project, error)
	DeleteProject(ctx context.Context, id int64) error
	ListProjects(
		ctx context.Context,
		userID int64,
		viewScope string,
		ownerID int64,
		limit, offset int,
	) ([]domain.Project, error)
	CountProjects(ctx context.Context, userID int64, viewScope string, ownerID int64) (int64, error)
	AutoCreatedCounts(ctx context.Context, projectID int64) (domain.AutoCreatedCounts, error)
}
