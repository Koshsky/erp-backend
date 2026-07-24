package delivery

import (
	"context"

	"github.com/Koshsky/erp/api/internal/scheduling/domain"
)

type MilestoneService interface {
	GetProjectScheduling(ctx context.Context) (*domain.ProjectScheduling, error)
	GetProcessScheduling(ctx context.Context) (*domain.ProcessScheduling, error)
	GetTaskScheduling(ctx context.Context) (*domain.TaskScheduling, error)
}
