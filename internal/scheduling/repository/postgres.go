//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp-backend/internal/middleware/auth"
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

func (r *SchedulingRepository) ListProjects(ctx context.Context) ([]domain.Project, error) {
	userID := auth.GetUserID(ctx)
	role := auth.GetRole(ctx)
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

func (r *SchedulingRepository) ListProcessesByProjectID(ctx context.Context, projectIDs []int64) ([]domain.Process, error) {
	rows, err := r.db.ListProcessesByProjectID(ctx, projectIDs)
	if err != nil {
		return nil, err
	}
	processes := make([]domain.Process, len(rows))
	for i, row := range rows {
		processes[i] = domain.Process{
			ID:        row.ID,
			OwnerID:   row.OwnerID,
			ProjectID: row.ProjectID,
			StartDate: row.StartDate,
			EndDate:   row.EndDate,
		}
	}
	return processes, nil
}

func (r *SchedulingRepository) GetProcessScheduling(ctx context.Context) (*domain.ProcessScheduling, error) {
	rows, err := r.db.GetProcessScheduling(ctx, sqlc.GetProcessSchedulingParams{
		Role:   auth.GetRole(ctx),
		UserID: auth.GetUserID(ctx),
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
	taskRows, err := r.db.GetDescribedTasks(ctx, sqlc.GetDescribedTasksParams{
		Role:   auth.GetRole(ctx),
		UserID: auth.GetUserID(ctx),
	})
	if err != nil {
		return nil, err
	}

	scheduling := domain.TaskScheduling{
		Tasks:       make([]domain.Task, len(taskRows)),
		Processes:   make(map[int64]domain.Process),
		Projects:    make(map[int64]domain.Project),
		Assignments: make(map[int64][]domain.Assignment),
		Milestones:  make(map[int64][]domain.Milestone),
		Resources:   make(map[int64]domain.Resource),
	}

	taskIDs := make([]int64, len(taskRows))
	processIDs := []int64{}
	for i, r := range taskRows {
		task := domain.Task{
			ID:        r.Task.ID,
			ProcessID: r.Task.ProcessID,
			Title:     r.Task.Title,
			StartDate: r.Task.StartDate,
			EndDate:   r.Task.EndDate,
		}
		process := domain.Process{
			ID:        r.Process.ID,
			OwnerID:   r.Process.OwnerID,
			ProjectID: r.Process.ProjectID,
			Title:     r.Process.Title,
			StartDate: r.Process.StartDate,
			EndDate:   r.Process.EndDate,
		}
		project := domain.Project{
			ID:        r.Project.ID,
			Code:      r.Project.Code,
			StartDate: r.Project.StartDate,
			EndDate:   r.Project.EndDate,
			Priority:  int(r.Project.Priority),
			OwnerID:   r.Project.OwnerID,
		}
		scheduling.Tasks[i] = task
		scheduling.Processes[process.ID] = process
		scheduling.Projects[project.ID] = project

		if i == 0 {
			scheduling.Timeline.StartDate = process.StartDate
			scheduling.Timeline.EndDate = process.EndDate
		} else {
			if scheduling.Timeline.StartDate.After(process.StartDate) {
				scheduling.Timeline.StartDate = process.StartDate
			}
			if scheduling.Timeline.EndDate.Before(process.EndDate) {
				scheduling.Timeline.EndDate = process.EndDate
			}
		}

		taskIDs[i] = r.Task.ID
		if len(processIDs) == 0 || processIDs[len(processIDs)-1] != r.Task.ProcessID {
			processIDs = append(processIDs, r.Task.ProcessID)
		}
	}
	duration := scheduling.Timeline.EndDate.Sub(scheduling.Timeline.StartDate)
	scheduling.Timeline.TotalDays = int(duration.Hours() / 24)

	milestoneRows, err := r.db.ListMilestonesByProcessIDs(ctx, processIDs)
	if err != nil {
		return nil, err
	}
	for _, milestone := range milestoneRows {
		if scheduling.Milestones[milestone.ProcessID] == nil {
			scheduling.Milestones[milestone.ProcessID] = []domain.Milestone{}
		}
		scheduling.Milestones[milestone.ID] = append(
			scheduling.Milestones[milestone.ProcessID], domain.Milestone{
				ID:        milestone.ID,
				ProcessID: milestone.ProcessID,
				Title:     milestone.Title,
				Content:   milestone.Content,
				Date:      milestone.Date,
			},
		)
	}

	assignmentRows, err := r.db.ListAssignmentsByTaskIDs(ctx, taskIDs)
	if err != nil {
		return nil, err
	}
	for _, assignment := range assignmentRows {
		if scheduling.Assignments[assignment.TaskID] == nil {
			scheduling.Assignments[assignment.TaskID] = []domain.Assignment{}
		}
		scheduling.Assignments[assignment.TaskID] = append(
			scheduling.Assignments[assignment.TaskID], domain.Assignment{
				ID:         assignment.ID,
				TaskID:     assignment.TaskID,
				ResourceID: assignment.ResourceID,
				Quantity:   int(assignment.Quantity),
			},
		)
	}

	resourceRows, err := r.db.GetResources(ctx)
	if err != nil {
		return nil, err
	}
	scheduling.Resources = make(map[int64]domain.Resource)
	for _, resource := range resourceRows {
		scheduling.Resources[resource.ID] = domain.Resource{
			ID:       resource.ID,
			Code:     resource.Code,
			Quantity: int(resource.Quantity),
			Title:    resource.Title,
		}
	}

	return &scheduling, nil
}
