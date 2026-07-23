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
