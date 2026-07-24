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

func (r *SchedulingRepository) GetProjectScheduling(ctx context.Context) (*domain.ProjectScheduling, error) {
	role := ctx.Value("role").(string)

	rows, err := r.db.GetProjectScheduling(ctx, role)
	if err != nil {
		return nil, err
	}
	scheduling := domain.ProjectScheduling{
		Projects: make([]domain.Project, len(rows)),
	}
	for i, row := range rows {
		scheduling.Projects[i] = domain.Project{
			ID:        row.ID,
			OwnerID:   row.OwnerID,
			Code:      row.Code,
			StartDate: row.StartDate,
			EndDate:   row.EndDate,
			Priority:  int(row.Priority),
		}
	}

	return &scheduling, nil
}

func (r *SchedulingRepository) GetProcessScheduling(ctx context.Context) (*domain.ProcessScheduling, error) {
	role := ctx.Value("role").(string)
	userID := ctx.Value("user_id").(int64)

	rows, err := r.db.GetProcessScheduling(ctx, sqlc.GetProcessSchedulingParams{
		Role:   role,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}
	scheduling := domain.ProcessScheduling{
		Processes: make([]domain.Process, len(rows)),
		Projects:  make(map[int64]domain.Project),
	}
	for i, row := range rows {
		process := row.Process
		project := row.Project
		scheduling.Processes[i] = domain.Process{
			ID:        process.ID,
			OwnerID:   process.OwnerID,
			Title:     process.Title,
			ProjectID: process.ProjectID,
			StartDate: process.StartDate,
			EndDate:   process.EndDate,
		}
		scheduling.Projects[project.ID] = domain.Project{
			ID:        project.ID,
			OwnerID:   project.OwnerID,
			Code:      project.Code,
			StartDate: project.StartDate,
			EndDate:   project.EndDate,
			Priority:  int(project.Priority),
		}
	}

	return &scheduling, nil
}

func (r *SchedulingRepository) GetTaskScheduling(ctx context.Context) (*domain.TaskScheduling, error) {
	role := ctx.Value("role").(string)
	userID := ctx.Value("user_id").(int64)

	rows, err := r.db.GetTaskScheduling(ctx, sqlc.GetTaskSchedulingParams{
		Role:   role,
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}
	resources, err := r.db.GetResources(ctx)
	if err != nil {
		return nil, err
	}

	scheduling := domain.TaskScheduling{
		Tasks:       make([]domain.Task, len(rows)),
		Projects:    make(map[int64]domain.Project),
		Processes:   make(map[int64]domain.Process),
		Assignments: make(map[int64][]domain.Assignment),
		Resources:   make(map[int64]domain.Resource),
	}

	for i, row := range rows {
		task := domain.Task{
			ID:        row.Task.ID,
			ProcessID: row.Task.ProcessID,
			Title:     row.Task.Title,
			StartDate: row.Task.StartDate,
			EndDate:   row.Task.EndDate,
		}
		project := domain.Project{
			ID:        row.Project.ID,
			OwnerID:   row.Project.OwnerID,
			Code:      row.Project.Code,
			StartDate: row.Project.StartDate,
			EndDate:   row.Project.EndDate,
			Priority:  int(row.Project.Priority),
		}
		process := domain.Process{
			ID:        row.Process.ID,
			OwnerID:   row.Process.OwnerID,
			ProjectID: row.Process.ProjectID,
			Title:     row.Process.Title,
			StartDate: row.Process.StartDate,
			EndDate:   row.Process.EndDate,
		}
		assignment := domain.Assignment{
			ID:         row.Assignment.ID,
			TaskID:     row.Assignment.TaskID,
			ResourceID: row.Assignment.ResourceID,
			Quantity:   int(row.Assignment.Quantity),
		}

		scheduling.Tasks[i] = task
		scheduling.Processes[process.ID] = process
		scheduling.Projects[project.ID] = project
		if scheduling.Assignments[task.ID] == nil {
			scheduling.Assignments[task.ID] = []domain.Assignment{}
		}
		scheduling.Assignments[task.ID] = append(scheduling.Assignments[task.ID], assignment)
	}

	for _, row := range resources {
		resource := domain.Resource{
			ID:       row.ID,
			Code:     row.Code,
			Title:    row.Title,
			Quantity: int(row.Quantity),
		}
		scheduling.Resources[resource.ID] = resource
	}
	return &scheduling, nil
}
