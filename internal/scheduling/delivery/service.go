package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/scheduling/dto"
)

type MilestoneService interface {
	GetProjectScheduling(ctx context.Context) (*dto.ProjectScheduling, error)
	GetProcessScheduling(ctx context.Context) (*dto.ProcessScheduling, error)
	GetTaskScheduling(ctx context.Context) (*dto.TaskScheduling, error)
}
