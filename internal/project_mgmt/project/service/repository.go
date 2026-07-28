package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/domain"
)

type RepositoryInterface interface {
	CreateProject(ctx context.Context, Project domain.Project) (*domain.Project, error)
	GetProject(ctx context.Context, id int64) (*domain.Project, error)
	UpdateProject(ctx context.Context, new domain.Project) (*domain.Project, error)
	DeleteProject(ctx context.Context, id int64) error
	ListProjects(ctx context.Context) ([]domain.Project, error)
}
