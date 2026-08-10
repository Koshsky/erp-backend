package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/dto"
)

type ProjectService interface {
	ListProjects(ctx context.Context) ([]dto.ProjectResponse, error)
	FindProject(ctx context.Context, id int64) (*dto.ProjectResponse, error)
	CreateProject(
		ctx context.Context,
		project dto.CreateProjectRequest,
		userID int64,
		role string,
	) (*dto.ProjectResponse, error)
	UpdateProject(ctx context.Context, id int64, project dto.UpdateProjectRequest) (*dto.ProjectResponse, error)
	DeleteProject(ctx context.Context, id int64) error
}
