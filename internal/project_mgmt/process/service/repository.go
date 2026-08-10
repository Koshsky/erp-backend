package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/domain"
)

type ProcessRepository interface {
	CreateProcess(ctx context.Context, Process domain.Process) (*domain.Process, error)
	FindProcess(ctx context.Context, id int64) (*domain.Process, error)
	UpdateProcess(ctx context.Context, process domain.Process) (*domain.Process, error)
	DeleteProcess(ctx context.Context, id int64) error
	ListProcesss(
		ctx context.Context,
		userID int64,
		role string,
		ownerID int64,
		limit, offset int,
	) ([]domain.Process, error)
	CountProcesses(ctx context.Context, userID int64, role string, ownerID int64) (int64, error)
}
