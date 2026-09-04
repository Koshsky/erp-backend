//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	errapi "github.com/Koshsky/erp-backend/pkg/errors"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/domain"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/repository/sqlc"
	nullable "github.com/Koshsky/erp-backend/pkg/database"
)

type TaskRepository struct {
	logger *slog.Logger
	pool   *pgxpool.Pool
	db     *sqlc.Queries
}

// NewTaskRepository builds the TaskRepository repository.
func NewTaskRepository(logger *slog.Logger, pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		logger: logger,
		pool:   pool,
		db:     sqlc.New(pool),
	}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task domain.Task) (*domain.Task, error) {
	row, err := r.db.CreateTask(ctx, sqlc.CreateTaskParams{
		ProcessID: task.ProcessID,
		OwnerID:   nullable.ToInt8(task.OwnerID),
		Title:     task.Title,
		Color:     nullable.ToString(task.Color),
		StartDate: task.StartDate,
		EndDate:   task.EndDate,
	})
	if err != nil {
		return nil, errapi.MapPgConstraint(errapi.FromPgInvalidParam(err))
	}

	mapped := mapTask(row)
	return &mapped, nil
}

func (r *TaskRepository) FindTask(ctx context.Context, id int64) (*domain.Task, error) {
	row, err := r.db.FindTask(ctx, id)
	if err != nil {
		return nil, err
	}

	mapped := mapTask(row)
	return &mapped, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, task domain.Task) (*domain.Task, error) {
	row, err := r.db.UpdateTask(ctx, sqlc.UpdateTaskParams{
		TaskID:    task.ID,
		ProcessID: task.ProcessID,
		OwnerID:   nullable.ToInt8(task.OwnerID),
		Title:     task.Title,
		Color:     nullable.ToString(task.Color),
		StartDate: task.StartDate,
		EndDate:   task.EndDate,
	})
	if err != nil {
		return nil, errapi.MapPgConstraint(errapi.FromPgInvalidParam(err))
	}

	mapped := mapTask(row)
	return &mapped, nil
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id int64) error {
	return r.db.DeleteTask(ctx, id)
}

func (r *TaskRepository) ListTasks(
	ctx context.Context,
	userID int64,
	viewScope string,
	ownerID int64,
	limit, offset int,
) ([]domain.Task, error) {
	rows, err := r.db.ListTasks(ctx, sqlc.ListTasksParams{
		ScopeView:  viewScope,
		UserID:     userID,
		OwnerID:    ownerID,
		PageLimit:  int64(limit),
		PageOffset: int64(offset),
	})
	if err != nil {
		return nil, err
	}
	tasks := make([]domain.Task, 0, len(rows))
	for _, row := range rows {
		tasks = append(tasks, mapTask(row))
	}
	return tasks, nil
}

func (r *TaskRepository) CountTasks(ctx context.Context, userID int64, viewScope string, ownerID int64) (int64, error) {
	return r.db.CountTasks(
		ctx,
		sqlc.CountTasksParams{
			ScopeView: viewScope,
			UserID:    userID,
			OwnerID:   ownerID,
		},
	)
}

func mapTask(row sqlc.Task) domain.Task {
	return domain.Task{
		ID:        row.ID,
		ProcessID: row.ProcessID,
		OwnerID:   nullable.Int64Ptr(row.OwnerID),
		Title:     row.Title,
		Color:     nullable.StringPtr(row.Color),
		StartDate: row.StartDate,
		EndDate:   row.EndDate,
		SortOrder: int(row.SortOrder),
	}
}

// ListTaskIDsByProcess returns the active task ids of a process in their
// display order — to validate a reorder request covers the whole group.
func (r *TaskRepository) ListTaskIDsByProcess(ctx context.Context, processID int64) ([]int64, error) {
	return r.db.ListTaskIdsByProcess(ctx, processID)
}

// ReorderTasks rewrites the sort_order of the given task ids by list position
// (1-based) in one transaction. The caller validates that the ids cover the
// whole group. The two-phase UPDATE parks the rows on offset slots first,
// because a single-statement value swap would transiently violate the partial
// unique index (process_id, sort_order).
func (r *TaskRepository) ReorderTasks(ctx context.Context, ids []int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := r.db.WithTx(tx)
	if err = q.ReorderTasksMark(ctx, ids); err != nil {
		return err
	}
	if err = q.ReorderTasksApply(ctx, ids); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// OwnerChain returns the owner chain (for RBAC checks in the middleware).
func (r *TaskRepository) OwnerChain(ctx context.Context, id int64) (rbac.Owners, error) {
	row, err := r.db.OwnerChain(ctx, id)
	if err != nil {
		return rbac.Owners{}, err
	}
	return rbac.Owners{ProjectOwner: row.ProjectOwner, ProcessOwner: row.ProcessOwner, Owner: row.OwnerID}, nil
}
