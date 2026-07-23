package postgres

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/repository/postgres/sqlc"
)

type ProjectRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

func NewProjectRepository(logger *slog.Logger, db *sqlc.Queries) *ProjectRepository {
	return &ProjectRepository{logger: logger, db: db}
}

func (r *ProjectRepository) CreateProject(ctx context.Context, project domain.Project) (*domain.Project, error) {
	created, err := r.db.CreateProject(ctx, sqlc.CreateProjectParams{
		Code:      project.Code,
		OwnerID:   project.OwnerID,
		StartDate: toDate(project.StartDate),
		EndDate:   toDate(project.EndDate),
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
		StartDate: toDate(project.StartDate),
		EndDate:   toDate(project.EndDate),
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
	rows, err := r.db.ListProjects(ctx, sqlc.ListProjectsParams{
		Role:   ctx.Value("role").(string),
		UserID: ctx.Value("user_id").(int64),
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
