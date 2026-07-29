package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/scheduling/dto"
)

type MilestoneService interface {
	GetProjectScheduling(ctx context.Context, userID int64, role string) (*dto.ProjectScheduling, error)
	GetProcessScheduling(ctx context.Context, userID int64, role string) (*dto.ProcessScheduling, error)
	GetTaskScheduling(ctx context.Context, userID int64, role string) (*dto.TaskScheduling, error)
}
