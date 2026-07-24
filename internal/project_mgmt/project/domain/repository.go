package domain

import (
	"context"
)

type RepositoryInterface interface {
	CreateProject(ctx context.Context, Project Project) (*Project, error)
	GetProject(ctx context.Context, id int64) (*Project, error)
	UpdateProject(ctx context.Context, new Project) (*Project, error)
	DeleteProject(ctx context.Context, id int64) error
	ListProjects(ctx context.Context) ([]Project, error)
}
