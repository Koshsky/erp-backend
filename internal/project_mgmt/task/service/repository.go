package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/domain"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task domain.Task) (*domain.Task, error)
	FindTask(ctx context.Context, id int64) (*domain.Task, error)
	UpdateTask(ctx context.Context, task domain.Task) (*domain.Task, error)
	DeleteTask(ctx context.Context, id int64) error
	ListTasks(ctx context.Context) ([]domain.Task, error)
}
