//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/scheduling/domain"
	"github.com/Koshsky/erp-backend/internal/scheduling/repository/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SchedulingRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

func NewSchedulingRepository(logger *slog.Logger, pool *pgxpool.Pool) *SchedulingRepository {
	return &SchedulingRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *SchedulingRepository) ListProjects(ctx context.Context, userID int64, role string) ([]domain.Project, error) {
	rows, err := r.db.ListProjects(ctx, sqlc.ListProjectsParams{
		UserID: userID,
		Role:   role,
	})

	if err != nil {
		return nil, err
	}
	projetcs := make([]domain.Project, len(rows))
	for i, row := range rows {
		projetcs[i] = domain.Project{
			ID:        row.ID,
			OwnerID:   row.OwnerID,
			Code:      row.Code,
			StartDate: row.StartDate,
			EndDate:   row.EndDate,
			Priority:  int(row.Priority),
		}
	}

	return projetcs, nil
}

func (r *SchedulingRepository) ListProcesses(ctx context.Context, userID int64, role string) ([]domain.Process, error) {
	rows, err := r.db.ListProcesses(ctx, sqlc.ListProcessesParams{
		UserID: userID,
		Role:   role,
	})
	if err != nil {
		return nil, err
	}
	processes := make([]domain.Process, len(rows))
	for i, row := range rows {
		processes[i] = domain.Process{
			ID:        row.Process.ID,
			OwnerID:   row.Process.OwnerID,
			ProjectID: row.Process.ProjectID,
			StartDate: row.Process.StartDate,
			EndDate:   row.Process.EndDate,
		}
	}
	return processes, nil
}

func (r *SchedulingRepository) ListProcessesByProjectIDs(ctx context.Context, projectIDs []int64) (map[int64][]domain.Process, error) {
	rows, err := r.db.ListProcessesByProjectIDs(ctx, projectIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[int64][]domain.Process)
	for _, row := range rows {
		p := domain.Process{
			ID:        row.ID,
			OwnerID:   row.OwnerID,
			ProjectID: row.ProjectID,
			StartDate: row.StartDate,
			EndDate:   row.EndDate,
		}
		if _, ok := result[row.ProjectID]; !ok {
			result[row.ProjectID] = make([]domain.Process, 0)
		}
		result[row.ProjectID] = append(result[row.ProjectID], p)
	}
	return result, nil

}

func (r *SchedulingRepository) ListTasksByProcessIDs(ctx context.Context, processIDs []int64) (map[int64][]domain.Task, error) {
	rows, err := r.db.ListTasksByProcessIDs(ctx, processIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[int64][]domain.Task)
	for _, row := range rows {
		t := domain.Task{
			ID:        row.ID,
			ProcessID: row.ProcessID,
			Title:     row.Title,
			StartDate: row.StartDate,
			EndDate:   row.EndDate,
		}
		if _, ok := result[row.ProcessID]; !ok {
			result[row.ProcessID] = make([]domain.Task, 0)
		}
		result[row.ProcessID] = append(result[row.ProcessID], t)
	}
	return result, nil
}

func (r *SchedulingRepository) ListMilestonesByProcessIDs(ctx context.Context, processIDs []int64) (map[int64][]domain.Milestone, error) {
	rows, err := r.db.ListMilestonesByProcessIDs(ctx, processIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[int64][]domain.Milestone)
	for _, row := range rows {
		m := domain.Milestone{
			ID:        row.ID,
			ProcessID: row.ProcessID,
			Title:     row.Title,
			Content:   row.Content,
			Date:      row.Date,
		}
		if _, ok := result[row.ProcessID]; !ok {
			result[row.ProcessID] = make([]domain.Milestone, 0)
		}
		result[row.ProcessID] = append(result[row.ProcessID], m)
	}
	return result, nil
}

func (r *SchedulingRepository) ListAssignmentsByTaskIDs(ctx context.Context, taskIDs []int64) (map[int64][]domain.Assignment, error) {
	rows, err := r.db.ListAssignmentsByTaskIDs(ctx, taskIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[int64][]domain.Assignment)
	for _, row := range rows {
		a := domain.Assignment{
			ID:         row.ID,
			TaskID:     row.TaskID,
			ResourceID: row.ResourceID,
			Quantity:   int(row.Quantity),
		}
		if _, ok := result[row.TaskID]; !ok {
			result[row.TaskID] = make([]domain.Assignment, 0)
		}
		result[row.TaskID] = append(result[row.TaskID], a)
	}
	return result, nil
}

func (r *SchedulingRepository) ListResources(ctx context.Context) ([]domain.Resource, error) {
	rows, err := r.db.ListResources(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]domain.Resource, 0, len(rows))
	for _, row := range rows {
		r := domain.Resource{
			ID:       row.ID,
			Title:    row.Title,
			Code:     row.Code,
			Quantity: int(row.Quantity),
		}
		result = append(result, r)
	}
	return result, nil
}
