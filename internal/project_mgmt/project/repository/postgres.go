//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/repository/sqlc"
	nullable "github.com/Koshsky/erp-backend/pkg/database"
	errapi "github.com/Koshsky/erp-backend/pkg/errors"
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

func (r *ProjectRepository) CreateProject(
	ctx context.Context,
	ownerID *int64,
	code string,
	color *string,
	startDate, endDate time.Time,
	priority int,
) (*sqlc.Project, error) {
	created, err := r.db.CreateProject(ctx, sqlc.CreateProjectParams{
		Code:      code,
		OwnerID:   nullable.ToInt8(ownerID),
		Color:     nullable.ToString(color),
		StartDate: startDate,
		EndDate:   endDate,
		Priority:  int64(priority),
	})
	if err != nil {
		// Idempotent create: the active code already exists (ON CONFLICT inserted
		// nothing) — this is not an internal error but a key conflict.
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errapi.Conflict("project with this code already exists")
		}
		return nil, mapProjectCreateError(err)
	}

	return &created, nil
}

func (r *ProjectRepository) FindProject(ctx context.Context, id int64) (*sqlc.Project, error) {
	project, err := r.db.FindProject(ctx, id)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *ProjectRepository) UpdateProject(
	ctx context.Context,
	project sqlc.Project,
	priority int,
) (*sqlc.Project, error) {
	updated, err := r.db.UpdateProject(ctx, sqlc.UpdateProjectParams{
		ProjectID: project.ID,
		OwnerID:   project.OwnerID,
		Code:      project.Code,
		Color:     project.Color,
		Priority:  int64(priority),
		StartDate: project.StartDate,
		EndDate:   project.EndDate,
	})
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, id int64) error {
	return r.db.DeleteProject(ctx, id)
}

func (r *ProjectRepository) ListProjects(
	ctx context.Context,
	userID int64,
	viewScope string,
	ownerID int64,
	limit, offset int,
) ([]sqlc.Project, error) {
	return r.db.ListProjects(ctx, sqlc.ListProjectsParams{
		ScopeView:  viewScope,
		UserID:     userID,
		OwnerID:    ownerID,
		PageLimit:  int64(limit),
		PageOffset: int64(offset),
	})
}

func (r *ProjectRepository) CountProjects(
	ctx context.Context,
	userID int64,
	viewScope string,
	ownerID int64,
) (int64, error) {
	return r.db.CountProjects(
		ctx,
		sqlc.CountProjectsParams{
			ScopeView: viewScope,
			UserID:    userID,
			OwnerID:   ownerID,
		},
	)
}

// AutoCreatedCounts returns what the auto-create trigger (V8) created for the
// project on insert (processes/tasks/assignments). All zero when the template
// is disabled or empty.
func (r *ProjectRepository) AutoCreatedCounts(
	ctx context.Context,
	projectID int64,
) (sqlc.CountAutoCreatedEntitiesRow, error) {
	return r.db.CountAutoCreatedEntities(ctx, projectID)
}

// OwnerChain returns the owner chain (for RBAC checks in the middleware).
func (r *ProjectRepository) OwnerChain(ctx context.Context, id int64) (rbac.Owners, error) {
	owner, err := r.db.OwnerChain(ctx, id)
	if err != nil {
		return rbac.Owners{}, err
	}
	return rbac.Owners{ProjectOwner: owner}, nil
}
