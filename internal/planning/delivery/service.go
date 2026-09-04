package delivery

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/planning/dto"
)

type MilestoneService interface {
	GetProjectPlanning(ctx context.Context, userID int64, viewScope string) (*dto.ProjectPlanning, error)
	GetProcessPlanning(ctx context.Context, userID int64, viewScope string) (*dto.ProcessPlanning, error)
	GetTaskPlanning(ctx context.Context, userID int64, viewScope string) (*dto.TaskPlanning, error)
}
