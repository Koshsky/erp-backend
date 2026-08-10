//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/domain"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/repository/sqlc"
	nullable "github.com/Koshsky/erp-backend/pkg/database"
)

type ProjectRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

// NewProjectRepository builds the ProjectRepository repository.
func NewProjectRepository(logger *slog.Logger, pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *ProjectRepository) CreateProject(ctx context.Context, project domain.Project) (*domain.Project, error) {
	created, err := r.db.CreateProject(ctx, sqlc.CreateProjectParams{
		Code:      project.Code,
		OwnerID:   nullable.ToInt8(project.OwnerID),
		StartDate: project.StartDate,
		EndDate:   project.EndDate,
		Priority:  int64(project.Priority),
	})
	if err != nil {
		return nil, err
	}

	mapped := mapProject(created)
	return &mapped, nil
}

func (r *ProjectRepository) FindProject(ctx context.Context, id int64) (*domain.Project, error) {
	project, err := r.db.FindProject(ctx, id)
	if err != nil {
		return nil, err
	}

	mapped := mapProject(project)
	return &mapped, nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, project domain.Project) (*domain.Project, error) {
	updated, err := r.db.UpdateProject(ctx, sqlc.UpdateProjectParams{
		ProjectID: project.ID,
		OwnerID:   nullable.ToInt8(project.OwnerID),
		Code:      project.Code,
		Priority:  int64(project.Priority),
		StartDate: project.StartDate,
		EndDate:   project.EndDate,
	})
	if err != nil {
		return nil, err
	}

	mapped := mapProject(updated)
	return &mapped, nil
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, id int64) error {
	return r.db.DeleteProject(ctx, id)
}

func (r *ProjectRepository) ListProjects(ctx context.Context, limit, offset int) ([]domain.Project, error) {
	rows, err := r.db.ListProjects(ctx, sqlc.ListProjectsParams{
		PageLimit:  int64(limit),
		PageOffset: int64(offset),
	})
	if err != nil {
		return nil, err
	}
	projects := make([]domain.Project, 0, len(rows))
	for _, row := range rows {
		projects = append(projects, mapProject(row))
	}
	return projects, nil
}

func (r *ProjectRepository) CountProjects(ctx context.Context) (int64, error) {
	return r.db.CountProjects(ctx)
}

func mapProject(row sqlc.Project) domain.Project {
	return domain.Project{
		ID:        row.ID,
		OwnerID:   nullable.Int64Ptr(row.OwnerID),
		Code:      row.Code,
		StartDate: row.StartDate,
		EndDate:   row.EndDate,
		Priority:  int(row.Priority),
	}
}

// OwnerChain returns the owner chain (for RBAC checks in the middleware).
func (r *ProjectRepository) OwnerChain(ctx context.Context, id int64) (rbac.Owners, error) {
	owner, err := r.db.OwnerChain(ctx, id)
	if err != nil {
		return rbac.Owners{}, err
	}
	return rbac.Owners{ProjectOwner: owner}, nil
}
