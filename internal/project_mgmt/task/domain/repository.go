package domain

import (
	"context"
)

type RepositoryInterface interface {
	CreateTask(ctx context.Context, task Task) (*Task, error)
	GetTask(ctx context.Context, id int64) (*Task, error)
	UpdateTask(ctx context.Context, task Task) (*Task, error)
	DeleteTask(ctx context.Context, id int64) error
	ListTasks(ctx context.Context) ([]Task, error)
}
