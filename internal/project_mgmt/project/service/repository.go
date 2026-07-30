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
	ListProjects(ctx context.Context) ([]domain.Project, error)
}
