//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/task/domain"
	"github.com/Koshsky/erp-backend/internal/task/repository/sqlc"
)

type TaskRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

func NewTaskRepository(logger *slog.Logger, pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task domain.Task) (*domain.Task, error) {
	row, err := r.db.CreateTask(ctx, sqlc.CreateTaskParams{
		ProcessID: task.ProcessID,
		Title:     task.Title,
		StartDate: task.StartDate,
		EndDate:   task.EndDate,
	})
	if err != nil {
		return nil, err
	}

	mapped := mapTask(row)
	return &mapped, nil
}

func (r *TaskRepository) GetTask(ctx context.Context, id int64) (*domain.Task, error) {
	row, err := r.db.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	mapped := mapTask(row)
	return &mapped, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, task domain.Task) (*domain.Task, error) {
	row, err := r.db.UpdateTask(ctx, sqlc.UpdateTaskParams{
		TaskID:    task.ID,
		Title:     task.Title,
		StartDate: task.StartDate,
		EndDate:   task.EndDate,
	})
	if err != nil {
		return nil, err
	}

	mapped := mapTask(row)
	return &mapped, nil
}

func (r *TaskRepository) DeleteTask(ctx context.Context, id int64) error {
	return r.db.DeleteTask(ctx, id)
}

func (r *TaskRepository) ListTasks(ctx context.Context) ([]domain.Task, error) {
	rows, err := r.db.ListTasks(ctx)
	if err != nil {
		return nil, err
	}
	tasks := make([]domain.Task, 0, len(rows))
	for _, row := range rows {
		tasks = append(tasks, mapTask(row))
	}
	return tasks, nil
}

func mapTask(row sqlc.Task) domain.Task {
	return domain.Task{
		ID:        row.ID,
		ProcessID: row.ProcessID,
		Title:     row.Title,
		StartDate: row.StartDate,
		EndDate:   row.EndDate,
	}
}
