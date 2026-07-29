package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/domain"
)

type ProcessRepository interface {
	CreateProcess(ctx context.Context, Process domain.Process) (*domain.Process, error)
	GetProcess(ctx context.Context, id int64) (*domain.Process, error)
	UpdateProcess(ctx context.Context, new domain.Process) (*domain.Process, error)
	DeleteProcess(ctx context.Context, id int64) error
	ListProcesss(ctx context.Context) ([]domain.Process, error)
}
