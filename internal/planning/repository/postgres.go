//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/planning/repository/sqlc"
)

type PlanningRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

// NewPlanningRepository builds the PlanningRepository repository.
func NewPlanningRepository(logger *slog.Logger, pool *pgxpool.Pool) *PlanningRepository {
	return &PlanningRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *PlanningRepository) ListProjects(ctx context.Context, userID int64, viewScope string) ([]sqlc.Project, error) {
	return r.db.ListProjects(ctx, sqlc.ListProjectsParams{
		UserID:    userID,
		ScopeView: viewScope,
	})
}

// ListProjectsByIDs returns full project rows by ids (for attaching parent
// context to process-scoped aggregates: /planning/processes).
func (r *PlanningRepository) ListProjectsByIDs(ctx context.Context, ids []int64) ([]sqlc.Project, error) {
	return r.db.ListProjectsByIDs(ctx, ids)
}

// ListProcesses — process-scoped list (process.view matrix).
func (r *PlanningRepository) ListProcesses(
	ctx context.Context,
	userID int64,
	viewScope string,
) ([]sqlc.ListProcessesRow, error) {
	return r.db.ListProcesses(ctx, sqlc.ListProcessesParams{
		UserID:    userID,
		ScopeView: viewScope,
	})
}

// ListProcessesByTaskScope — process-scoped list for the TASK planning
// aggregate. The scope comes from the task matrix: 'parent' means "in my
// processes" (p.owner_id) rather than "in my projects" (pr.owner_id), so a
// process owner (vp) sees their processes in the task diagram.
func (r *PlanningRepository) ListProcessesByTaskScope(
	ctx context.Context,
	userID int64,
	viewScope string,
) ([]sqlc.ListProcessesByTaskScopeRow, error) {
	return r.db.ListProcessesByTaskScope(ctx, sqlc.ListProcessesByTaskScopeParams{
		UserID:    userID,
		ScopeView: viewScope,
	})
}

func (r *PlanningRepository) ListTasksByProcessIDs(
	ctx context.Context,
	processIDs []int64,
) ([]sqlc.Task, error) {
	return r.db.ListTasksByProcessIDs(ctx, processIDs)
}

func (r *PlanningRepository) ListMilestonesByProcessIDs(
	ctx context.Context,
	processIDs []int64,
) ([]sqlc.Milestone, error) {
	return r.db.ListMilestonesByProcessIDs(ctx, processIDs)
}

func (r *PlanningRepository) ListAssignmentsByTaskIDs(
	ctx context.Context,
	taskIDs []int64,
) ([]sqlc.Assignment, error) {
	return r.db.ListAssignmentsByTaskIDs(ctx, taskIDs)
}

// ListTaskCommentCountsByTaskIDs returns the number of active comments per task.
func (r *PlanningRepository) ListTaskCommentCountsByTaskIDs(
	ctx context.Context,
	taskIDs []int64,
) (map[int64]int64, error) {
	rows, err := r.db.ListTaskCommentCountsByTaskIDs(ctx, taskIDs)
	if err != nil {
		return nil, err
	}
	counts := make(map[int64]int64, len(rows))
	for _, row := range rows {
		counts[row.TaskID] = row.CommentsCount
	}
	return counts, nil
}

func (r *PlanningRepository) ListResources(ctx context.Context) ([]sqlc.Resource, error) {
	return r.db.ListResources(ctx)
}
