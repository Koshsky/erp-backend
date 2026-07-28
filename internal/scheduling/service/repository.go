package domain

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/scheduling/domain"
	"github.com/Koshsky/erp-backend/internal/scheduling/dto"
)

type RepositoryInterface interface {
	ListProjects(ctx context.Context) ([]domain.Project, error)
	ListProcessesByProjectID(ctx context.Context, projectIDs []int64) ([]domain.Process, error)
	GetProcessScheduling(ctx context.Context) (*dto.ProcessScheduling, error)
	GetTaskScheduling(ctx context.Context) (*dto.TaskScheduling, error)
}
