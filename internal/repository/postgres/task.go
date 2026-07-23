package postgres

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/repository/postgres/sqlc"
)

type TaskRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

func NewTaskRepository(logger *slog.Logger, db *sqlc.Queries) *TaskRepository {
	return &TaskRepository{logger: logger, db: db}
}

func (r *TaskRepository) CreateTask(ctx context.Context, task domain.Task) (*domain.Task, error) {
	row, err := r.db.CreateTask(ctx, sqlc.CreateTaskParams{
		ProcessID: task.ProcessID,
		Title:     task.Title,
		StartDate: toDate(task.StartDate),
		EndDate:   toDate(task.EndDate),
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
		StartDate: toDate(task.StartDate),
		EndDate:   toDate(task.EndDate),
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