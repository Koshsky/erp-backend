package service

import (
	"context"

	"github.com/Koshsky/erp-backend/internal/planning/repository/sqlc"
)

type PlanningRepository interface {
	ListProjects(ctx context.Context, userID int64, viewScope string) ([]sqlc.Project, error)
	ListProjectsByIDs(ctx context.Context, ids []int64) ([]sqlc.Project, error)
	ListProcesses(ctx context.Context, userID int64, viewScope string) ([]sqlc.ListProcessesRow, error)
	ListProcessesByTaskScope(
		ctx context.Context,
		userID int64,
		viewScope string,
	) ([]sqlc.ListProcessesByTaskScopeRow, error)
	ListResources(ctx context.Context) ([]sqlc.Resource, error)

	ListMilestonesByProcessIDs(ctx context.Context, processIDs []int64) ([]sqlc.Milestone, error)
	ListTasksByProcessIDs(ctx context.Context, processIDs []int64) ([]sqlc.Task, error)
	ListAssignmentsByTaskIDs(ctx context.Context, taskIDs []int64) ([]sqlc.Assignment, error)
	ListTaskCommentCountsByTaskIDs(ctx context.Context, taskIDs []int64) (map[int64]int64, error)
}
