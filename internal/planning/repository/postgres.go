//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/common/date"
	"github.com/Koshsky/erp-backend/internal/planning/dto"
	"github.com/Koshsky/erp-backend/internal/planning/repository/sqlc"
)

type PlanningRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

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
			OwnerID:   row.OwnerID,
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
			ID:        row.Process.ID,
			Title:     row.Process.Title,
			OwnerID:   row.Process.OwnerID,
			ProjectID: row.Process.ProjectID,
			StartDate: date.From(row.Process.StartDate),
			EndDate:   date.From(row.Process.EndDate),
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
	result := make(map[int64][]dto.Process)
	for _, row := range rows {
		p := dto.Process{
			ID:        row.ID,
			Title:     row.Title,
			OwnerID:   row.OwnerID,
			ProjectID: row.ProjectID,
			StartDate: date.From(row.StartDate),
			EndDate:   date.From(row.EndDate),
		}
		if _, ok := result[row.ProjectID]; !ok {
			result[row.ProjectID] = make([]dto.Process, 0)
		}
		result[row.ProjectID] = append(result[row.ProjectID], p)
	}
	return result, nil
}

func (r *PlanningRepository) ListTasksByProcessIDs(
	ctx context.Context,
	processIDs []int64,
) (map[int64][]dto.Task, error) {
	rows, err := r.db.ListTasksByProcessIDs(ctx, processIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[int64][]dto.Task)
	for _, row := range rows {
		t := dto.Task{
			ID:        row.ID,
			ProcessID: row.ProcessID,
			Title:     row.Title,
			StartDate: date.From(row.StartDate),
			EndDate:   date.From(row.EndDate),
		}
		if _, ok := result[row.ProcessID]; !ok {
			result[row.ProcessID] = make([]dto.Task, 0)
		}
		result[row.ProcessID] = append(result[row.ProcessID], t)
	}
	return result, nil
}

func (r *PlanningRepository) ListMilestonesByProcessIDs(
	ctx context.Context,
	processIDs []int64,
) (map[int64][]dto.Milestone, error) {
	rows, err := r.db.ListMilestonesByProcessIDs(ctx, processIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[int64][]dto.Milestone)
	for _, row := range rows {
		m := dto.Milestone{
			ID:        row.ID,
			ProcessID: row.ProcessID,
			Title:     row.Title,
			Content:   row.Content,
			Date:      date.From(row.Date),
		}
		if _, ok := result[row.ProcessID]; !ok {
			result[row.ProcessID] = make([]dto.Milestone, 0)
		}
		result[row.ProcessID] = append(result[row.ProcessID], m)
	}
	return result, nil
}

func (r *PlanningRepository) ListAssignmentsByTaskIDs(
	ctx context.Context,
	taskIDs []int64,
) (map[int64][]dto.Assignment, error) {
	rows, err := r.db.ListAssignmentsByTaskIDs(ctx, taskIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[int64][]dto.Assignment)
	for _, row := range rows {
		a := dto.Assignment{
			ID:         row.ID,
			TaskID:     row.TaskID,
			ResourceID: row.ResourceID,
			Quantity:   int(row.Quantity),
		}
		if _, ok := result[row.TaskID]; !ok {
			result[row.TaskID] = make([]dto.Assignment, 0)
		}
		result[row.TaskID] = append(result[row.TaskID], a)
	}
	return result, nil
}

func (r *PlanningRepository) ListResources(ctx context.Context) ([]dto.Resource, error) {
	rows, err := r.db.ListResources(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]dto.Resource, 0, len(rows))
	for _, row := range rows {
		r := dto.Resource{
			ID:       row.ID,
			Title:    row.Title,
			Code:     row.Code,
			Quantity: int(row.Quantity),
		}
		result = append(result, r)
	}
	return result, nil
}
