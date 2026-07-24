package domain

import (
	"context"
)

type RepositoryInterface interface {
	GetProjectScheduling(ctx context.Context) (*ProjectScheduling, error)
	GetProcessScheduling(ctx context.Context) (*ProcessScheduling, error)
	GetTaskScheduling(ctx context.Context) (*TaskScheduling, error)
}
