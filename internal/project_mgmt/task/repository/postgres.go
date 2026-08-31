//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	errapi "github.com/Koshsky/erp-backend/pkg/errors"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/policies"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/domain"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/repository/sqlc"
	nullable "github.com/Koshsky/erp-backend/pkg/database"
)

type TaskRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

// NewTaskRepository builds the TaskRepository repository.
func NewTaskRepository(logger *slog.Logger, pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task domain.Task) (*domain.Task, error) {
	row, err := r.db.CreateTask(ctx, sqlc.CreateTaskParams{
		ProcessID: task.ProcessID,
		OwnerID:   nullable.ToInt8(task.OwnerID),
		Title:     task.Title,
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
	role string,
	ownerID int64,
	limit, offset int,
) ([]domain.Task, error) {
	rows, err := r.db.ListTasks(ctx, sqlc.ListTasksParams{
		ScopeView:  policies.ViewScopeCode(role, rbac.ResourceTask),
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

func (r *TaskRepository) CountTasks(ctx context.Context, userID int64, role string, ownerID int64) (int64, error) {
	return r.db.CountTasks(
		ctx,
		sqlc.CountTasksParams{
			ScopeView: policies.ViewScopeCode(role, rbac.ResourceTask),
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
		StartDate: row.StartDate,
		EndDate:   row.EndDate,
	}
}

// OwnerChain returns the owner chain (for RBAC checks in the middleware).
func (r *TaskRepository) OwnerChain(ctx context.Context, id int64) (rbac.Owners, error) {
	row, err := r.db.OwnerChain(ctx, id)
	if err != nil {
		return rbac.Owners{}, err
	}
	return rbac.Owners{ProjectOwner: row.ProjectOwner, ProcessOwner: row.ProcessOwner, Owner: row.OwnerID}, nil
}
