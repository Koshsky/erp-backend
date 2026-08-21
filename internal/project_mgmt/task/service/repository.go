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
	ListTasks(ctx context.Context, userID int64, role string, ownerID int64, limit, offset int) ([]domain.Task, error)
	CountTasks(ctx context.Context, userID int64, role string, ownerID int64) (int64, error)
}
