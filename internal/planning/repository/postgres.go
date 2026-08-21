//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/planning/dto"
	"github.com/Koshsky/erp-backend/internal/planning/repository/sqlc"
	nullable "github.com/Koshsky/erp-backend/pkg/database"
	"github.com/Koshsky/erp-backend/pkg/date"
)

type PlanningRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

// NewPlanningRepository builds the PlanningRepository repository.
func NewPlanningRepository(logger *slog.Logger, pool *pgxpool.Pool) *PlanningRepository {
	return &PlanningRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *PlanningRepository) ListProjects(ctx context.Context, userID int64, role string) ([]dto.Project, error) {
	rows, err := r.db.ListProjects(ctx, sqlc.ListProjectsParams{
		UserID: userID,
		Role:   role,
	})
	if err != nil {
		return nil, err
	}
	projetcs := make([]dto.Project, len(rows))
	for i, row := range rows {
		projetcs[i] = dto.Project{
			ID:        row.ID,
			OwnerID:   nullable.Int64Ptr(row.OwnerID),
			Code:      row.Code,
			StartDate: date.From(row.StartDate),
			EndDate:   date.From(row.EndDate),
			Priority:  int(row.Priority),
		}
	}

	return projetcs, nil
}

func (r *PlanningRepository) ListProcesses(ctx context.Context, userID int64, role string) ([]dto.Process, error) {
	rows, err := r.db.ListProcesses(ctx, sqlc.ListProcessesParams{
		UserID: userID,
		Role:   role,
	})
	if err != nil {
		return nil, err
	}
	processes := make([]dto.Process, len(rows))
	for i, row := range rows {
		processes[i] = dto.Process{
			ID:          row.Process.ID,
			Title:       row.Process.Title,
			OwnerID:     nullable.Int64Ptr(row.Process.OwnerID),
			ProjectID:   row.Process.ProjectID,
			ProjectCode: row.ProjectCode,
			StartDate:   date.From(row.Process.StartDate),
			EndDate:     date.From(row.Process.EndDate),
		}
	}
	return processes, nil
}

func (r *PlanningRepository) ListProcessesByProjectIDs(
	ctx context.Context,
	projectIDs []int64,
) (map[int64][]dto.Process, error) {
	rows, err := r.db.ListProcessesByProjectIDs(ctx, projectIDs)
	if err != nil {
		return nil, err
	}
	processes := make([]dto.Process, len(rows))
	for i, row := range rows {
		processes[i] = dto.Process{
			ID:        row.ID,
			Title:     row.Title,
			OwnerID:   nullable.Int64Ptr(row.OwnerID),
			ProjectID: row.ProjectID,
			StartDate: date.From(row.StartDate),
			EndDate:   date.From(row.EndDate),
		}
	}
	return groupByKey(processes, func(p dto.Process) int64 { return p.ProjectID }), nil
}

func (r *PlanningRepository) ListTasksByProcessIDs(
	ctx context.Context,
	processIDs []int64,
) (map[int64][]dto.Task, error) {
	rows, err := r.db.ListTasksByProcessIDs(ctx, processIDs)
	if err != nil {
		return nil, err
	}
	tasks := make([]dto.Task, len(rows))
	for i, row := range rows {
		tasks[i] = dto.Task{
			ID:        row.ID,
			ProcessID: row.ProcessID,
			OwnerID:   nullable.Int64Ptr(row.OwnerID),
			Title:     row.Title,
			StartDate: date.From(row.StartDate),
			EndDate:   date.From(row.EndDate),
		}
	}
	return groupByKey(tasks, func(t dto.Task) int64 { return t.ProcessID }), nil
}

func (r *PlanningRepository) ListMilestonesByProcessIDs(
	ctx context.Context,
	processIDs []int64,
) (map[int64][]dto.Milestone, error) {
	rows, err := r.db.ListMilestonesByProcessIDs(ctx, processIDs)
	if err != nil {
		return nil, err
	}
	milestones := make([]dto.Milestone, len(rows))
	for i, row := range rows {
		milestones[i] = dto.Milestone{
			ID:        row.ID,
			ProcessID: row.ProcessID,
			Title:     row.Title,
			Content:   row.Content,
			Date:      date.From(row.Date),
		}
	}
	return groupByKey(milestones, func(m dto.Milestone) int64 { return m.ProcessID }), nil
}

func (r *PlanningRepository) ListAssignmentsByTaskIDs(
	ctx context.Context,
	taskIDs []int64,
) (map[int64][]dto.Assignment, error) {
	rows, err := r.db.ListAssignmentsByTaskIDs(ctx, taskIDs)
	if err != nil {
		return nil, err
	}
	assignments := make([]dto.Assignment, len(rows))
	for i, row := range rows {
		assignments[i] = dto.Assignment{
			ID:         row.ID,
			TaskID:     row.TaskID,
			ResourceID: row.ResourceID,
			Quantity:   int(row.Quantity),
		}
	}
	return groupByKey(assignments, func(a dto.Assignment) int64 { return a.TaskID }), nil
}

// groupByKey groups items by the key returned by the key function.
func groupByKey[T any](items []T, key func(T) int64) map[int64][]T {
	result := make(map[int64][]T)
	for _, item := range items {
		k := key(item)
		if _, ok := result[k]; !ok {
			result[k] = make([]T, 0)
		}
		result[k] = append(result[k], item)
	}
	return result
}

func (r *PlanningRepository) ListResources(ctx context.Context) ([]dto.Resource, error) {
	rows, err := r.db.ListResources(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.Resource, 0, len(rows))
	for _, row := range rows {
		r := dto.Resource{
			ID:    row.ID,
			Title: row.Title,
			Code:  row.Code,
		}
		result = append(result, r)
	}
	return result, nil
}
