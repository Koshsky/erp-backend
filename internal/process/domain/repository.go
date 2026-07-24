package domain

import (
	"context"
)

type RepositoryInterface interface {
	CreateProcess(ctx context.Context, Process Process) (*Process, error)
	GetProcess(ctx context.Context, id int64) (*Process, error)
	UpdateProcess(ctx context.Context, new Process) (*Process, error)
	DeleteProcess(ctx context.Context, id int64) error
	ListProcesss(ctx context.Context) ([]Process, error)
}
