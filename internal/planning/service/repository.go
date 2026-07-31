package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/planning/dto"
)

type PlanningRepository interface {
	ListProjects(ctx context.Context, userID int64, role string) ([]dto.Project, error)
	ListProcesses(ctx context.Context, userID int64, role string) ([]dto.Process, error)
	ListResources(ctx context.Context) ([]dto.Resource, error)

	ListProcessesByProjectIDs(ctx context.Context, projectIDs []int64) (map[int64][]dto.Process, error)
	ListMilestonesByProcessIDs(ctx context.Context, processIDs []int64) (map[int64][]dto.Milestone, error)
	ListTasksByProcessIDs(ctx context.Context, processIDs []int64) (map[int64][]dto.Task, error)
	ListAssignmentsByTaskIDs(ctx context.Context, taskIDs []int64) (map[int64][]dto.Assignment, error)
}
