package domain

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/scheduling/domain"
)

type RepositoryInterface interface {
	ListProjects(ctx context.Context, userID int64, role string) ([]domain.Project, error)
	ListProcesses(ctx context.Context, userID int64, role string) ([]domain.Process, error)
	ListResources(ctx context.Context) ([]domain.Resource, error)

	ListProcessesByProjectIDs(ctx context.Context, projectIDs []int64) (map[int64][]domain.Process, error)
	ListMilestonesByProcessIDs(ctx context.Context, processIDs []int64) (map[int64][]domain.Milestone, error)
	ListTasksByProcessIDs(ctx context.Context, processIDs []int64) (map[int64][]domain.Task, error)
	ListAssignmentsByTaskIDs(ctx context.Context, taskIDs []int64) (map[int64][]domain.Assignment, error)
}
