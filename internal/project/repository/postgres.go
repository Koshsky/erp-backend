//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp/api/internal/project/domain"
	"github.com/Koshsky/erp/api/internal/project/repository/sqlc"
)

type ProjectRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

func NewProjectRepository(logger *slog.Logger, pool *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *ProjectRepository) CreateProject(ctx context.Context, project domain.Project) (*domain.Project, error) {
	created, err := r.db.CreateProject(ctx, sqlc.CreateProjectParams{
		Code:      project.Code,
		OwnerID:   project.OwnerID,
		StartDate: project.StartDate,
		EndDate:   project.EndDate,
		Priority:  int32(project.Priority),
	})
	if err != nil {
		return nil, err
	}

	mapped := mapProject(created)
	return &mapped, nil
}

func (r *ProjectRepository) GetProject(ctx context.Context, id int64) (*domain.Project, error) {
	project, err := r.db.GetProject(ctx, id)
	if err != nil {
		return nil, err
	}

	mapped := mapProject(project)
	return &mapped, nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, project domain.Project) (*domain.Project, error) {
	updated, err := r.db.UpdateProject(ctx, sqlc.UpdateProjectParams{
		ProjectID: project.ID,
		OwnerID:   project.OwnerID,
		Code:      project.Code,
		Priority:  int32(project.Priority),
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

func (r *ProjectRepository) ListProjects(ctx context.Context) ([]domain.Project, error) {
	rows, err := r.db.ListProjects(ctx)
	if err != nil {
		return nil, err
	}
	projects := make([]domain.Project, 0, len(rows))
	for _, row := range rows {
		projects = append(projects, mapProject(row))
	}
	return projects, nil
}

func mapProject(row sqlc.Project) domain.Project {
	return domain.Project{
		ID:        row.ID,
		OwnerID:   row.OwnerID,
		Code:      row.Code,
		StartDate: row.StartDate,
		EndDate:   row.EndDate,
		Priority:  int(row.Priority),
	}
}
