package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/dto"
)

type ProcessService interface {
	ListProcesses(ctx context.Context) ([]dto.ProcessResponse, error)
	GetProcess(ctx context.Context, id int64) (*dto.ProcessResponse, error)
	CreateProcess(ctx context.Context, process dto.CreateProcessRequest) (*dto.ProcessResponse, error)
	DeleteProcess(ctx context.Context, id int64) error
	UpdateProcess(ctx context.Context, id int64, process dto.UpdateProcessRequest) (*dto.ProcessResponse, error)
}
