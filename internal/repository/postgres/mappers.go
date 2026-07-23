package postgres

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/repository/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func mapProject(row sqlc.Project) domain.Project {
	return domain.Project{
		ID:        row.ID,
		OwnerID:   row.OwnerID,
		Code:      row.Code,
		StartDate: fromDate(row.StartDate),
		EndDate:   fromDate(row.EndDate),
		Priority:  int(row.Priority),
	}
}

func mapProcess(row sqlc.Process) domain.Process {
	return domain.Process{
		ID:        row.ID,
		OwnerID:   row.OwnerID,
		ProjectID: row.ProjectID,
		Title:     row.Title,
		StartDate: fromDate(row.StartDate),
		EndDate:   fromDate(row.EndDate),
	}
}

func mapMilestone(row sqlc.Milestone) domain.Milestone {
	return domain.Milestone{
		ID:        row.ID,
		ProcessID: row.ProcessID,
		Title:     row.Title,
		Content:   row.Content,
		Date:      fromDate(row.Date),
	}
}

func mapResource(row sqlc.Resource) domain.Resource {
	return domain.Resource{
		ID:       row.ID,
		Title:    row.Title,
		Code:     row.Code,
		Quantity: int(row.Quantity),
	}
}

func mapTask(row sqlc.Task) domain.Task {
	return domain.Task{
		ID:        row.ID,
		ProcessID: row.ProcessID,
		Title:     row.Title,
		StartDate: fromDate(row.StartDate),
		EndDate:   fromDate(row.EndDate),
	}
}

func mapAssignment(row sqlc.Assignment) domain.Assignment {
	return domain.Assignment{
		ID:         row.ID,
		TaskID:     row.TaskID,
		ResourceID: row.ResourceID,
		Quantity:   int(row.Quantity),
	}
}

func mapUser(row sqlc.User) domain.User {
	return domain.User{
		ID:           row.ID,
		Name:         row.Name,
		Role:         domain.UserRole(row.Role),
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
	}
}

func toDate(value time.Time) pgtype.Date {
	return pgtype.Date{
		Time:  value,
		Valid: true,
	}
}

func fromDate(value pgtype.Date) time.Time {
	if !value.Valid {
		return time.Time{}
	}

	return value.Time
}

func mapDetailedTask(row sqlc.GetDetailedTaskRow) (domain.DetailedTask, error) {
	var assignments []domain.Assignment

	if len(row.Assignments) > 0 && string(row.Assignments) != "[]" {
		if err := json.Unmarshal(row.Assignments, &assignments); err != nil {
			return domain.DetailedTask{}, fmt.Errorf("parse assignments: %w", err)
		}
	}

	return domain.DetailedTask{
		Task: domain.Task{
			ID:        row.ID,
			ProcessID: row.ProcessID,
			Title:     row.Title,
			StartDate: fromDate(row.StartDate),
			EndDate:   fromDate(row.EndDate),
		},
		Assignments: assignments,
	}, nil
}

func mapDetailedProcess(row sqlc.GetDetailedProcessRow) (domain.DetailedProcess, error) {
	var tasks []domain.DetailedTask

	if row.Tasks != nil {
		var jsonData []byte
		switch v := row.Tasks.(type) {
		case []byte:
			jsonData = v
		case string:
			jsonData = []byte(v)
		default:
			data, err := json.Marshal(v)
			if err != nil {
				return domain.DetailedProcess{}, fmt.Errorf("marshal tasks: %w", err)
			}
			jsonData = data
		}

		if len(jsonData) > 0 && string(jsonData) != "[]" {
			var rawTasks []struct {
				ID          int64       `json:"id"`
				ProcessID   int64       `json:"process_id"`
				Title       string      `json:"title"`
				StartDate   pgtype.Date `json:"start_date"`
				EndDate     pgtype.Date `json:"end_date"`
				Assignments []struct {
					ID         int64 `json:"id"`
					TaskID     int64 `json:"task_id"`
					ResourceID int64 `json:"resource_id"`
					Quantity   int   `json:"quantity"`
				} `json:"assignments"`
			}

			if err := json.Unmarshal(jsonData, &rawTasks); err != nil {
				return domain.DetailedProcess{}, fmt.Errorf("parse tasks: %w", err)
			}

			tasks = make([]domain.DetailedTask, len(rawTasks))
			for i, t := range rawTasks {
				assignments := make([]domain.Assignment, len(t.Assignments))
				for j, a := range t.Assignments {
					assignments[j] = domain.Assignment{
						ID:         a.ID,
						TaskID:     a.TaskID,
						ResourceID: a.ResourceID,
						Quantity:   a.Quantity,
					}
				}

				tasks[i] = domain.DetailedTask{
					Task: domain.Task{
						ID:        t.ID,
						ProcessID: t.ProcessID,
						Title:     t.Title,
						StartDate: fromDate(t.StartDate),
						EndDate:   fromDate(t.EndDate),
					},
					Assignments: assignments,
				}
			}
		}
	}

	var milestones []domain.Milestone

	if row.Milestones != nil {
		var jsonData []byte
		switch v := row.Milestones.(type) {
		case []byte:
			jsonData = v
		case string:
			jsonData = []byte(v)
		default:
			data, err := json.Marshal(v)
			if err != nil {
				return domain.DetailedProcess{}, fmt.Errorf("marshal milestones: %w", err)
			}
			jsonData = data
		}

		if len(jsonData) > 0 && string(jsonData) != "[]" {
			var rawMilestones []struct {
				ID        int64       `json:"id"`
				ProcessID int64       `json:"process_id"`
				Content   string      `json:"content"`
				Title     string      `json:"title"`
				Date      pgtype.Date `json:"date"`
			}

			if err := json.Unmarshal(jsonData, &rawMilestones); err != nil {
				return domain.DetailedProcess{}, fmt.Errorf("parse milestones: %w", err)
			}

			milestones = make([]domain.Milestone, len(rawMilestones))
			for i, m := range rawMilestones {
				milestones[i] = domain.Milestone{
					ID:        m.ID,
					Content:   m.Content,
					ProcessID: m.ProcessID,
					Title:     m.Title,
					Date:      fromDate(m.Date),
				}
			}
		}
	}

	return domain.DetailedProcess{
		Process: domain.Process{
			ID:        row.ID,
			OwnerID:   row.OwnerID,
			ProjectID: row.ProjectID,
			Title:     row.Title,
			StartDate: fromDate(row.StartDate),
			EndDate:   fromDate(row.EndDate),
		},
		Tasks:      tasks,
		Milestones: milestones,
	}, nil
}

func mapDetailedProject(
	project *domain.Project,
	processRows []sqlc.Process,
	taskRows []sqlc.ListTasksWithAssignmentsByProcessIDsRow,
	milestoneRows []sqlc.Milestone,
) *domain.DetailedProject {
	processIDs := make([]int64, 0, len(processRows))
	processMap := make(map[int64]domain.Process, len(processRows))
	for _, row := range processRows {
		processIDs = append(processIDs, row.ID)
		processMap[row.ID] = mapProcess(row)
	}

	if len(processIDs) == 0 {
		return &domain.DetailedProject{
			Project:   *project,
			Processes: []domain.DetailedProcess{},
		}
	}

	tasksMap := mapDetailedProjectTasks(taskRows)

	milestonesMap := make(map[int64][]domain.Milestone, len(milestoneRows))
	for _, row := range milestoneRows {
		milestonesMap[row.ProcessID] = append(milestonesMap[row.ProcessID], mapMilestone(row))
	}

	processes := make([]domain.DetailedProcess, 0, len(processIDs))
	for _, pid := range processIDs {
		p := processMap[pid]
		processes = append(processes, domain.DetailedProcess{
			Process:    p,
			Tasks:      tasksMap[pid],
			Milestones: milestonesMap[pid],
		})
	}

	return &domain.DetailedProject{
		Project:   *project,
		Processes: processes,
	}
}

func mapDetailedProjectTasks(rows []sqlc.ListTasksWithAssignmentsByProcessIDsRow) map[int64][]domain.DetailedTask {
	tasksMap := make(map[int64][]domain.DetailedTask)

	for _, row := range rows {
		processID := row.ProcessID

		var found bool
		for i, task := range tasksMap[processID] {
			if task.ID == row.TaskID {
				if row.AssignmentID.Valid {
					tasksMap[processID][i].Assignments = append(
						tasksMap[processID][i].Assignments,
						domain.Assignment{
							ID:         row.AssignmentID.Int64,
							TaskID:     row.AssignmentTaskID.Int64,
							ResourceID: row.ResourceID.Int64,
							Quantity:   int(row.AssignmentQuantity.Int32),
						},
					)
				}
				found = true
				break
			}
		}

		if !found {
			newTask := domain.DetailedTask{
				Task: domain.Task{
					ID:        row.TaskID,
					ProcessID: row.ProcessID,
					Title:     row.TaskTitle,
					StartDate: fromDate(row.TaskStartDate),
					EndDate:   fromDate(row.TaskEndDate),
				},
				Assignments: []domain.Assignment{},
			}

			if row.AssignmentID.Valid {
				newTask.Assignments = append(newTask.Assignments, domain.Assignment{
					ID:         row.AssignmentID.Int64,
					TaskID:     row.AssignmentTaskID.Int64,
					ResourceID: row.ResourceID.Int64,
					Quantity:   int(row.AssignmentQuantity.Int32),
				})
			}

			tasksMap[processID] = append(tasksMap[processID], newTask)
		}
	}

	return tasksMap
}
