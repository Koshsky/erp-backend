package service

import (
	"context"
	"time"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/repository/sqlc"
)

type ProjectRepository interface {
	CreateProject(
		ctx context.Context,
		ownerID *int64,
		code string,
		color *string,
		startDate, endDate time.Time,
		priority int,
	) (*sqlc.Project, error)
	FindProject(ctx context.Context, id int64) (*sqlc.Project, error)
	UpdateProject(ctx context.Context, project sqlc.Project, priority int) (*sqlc.Project, error)
	DeleteProject(ctx context.Context, id int64) error
	ListProjects(
		ctx context.Context,
		userID int64,
		viewScope string,
		ownerID int64,
		limit, offset int,
	) ([]sqlc.Project, error)
	CountProjects(ctx context.Context, userID int64, viewScope string, ownerID int64) (int64, error)
	AutoCreatedCounts(ctx context.Context, projectID int64) (sqlc.CountAutoCreatedEntitiesRow, error)
}
