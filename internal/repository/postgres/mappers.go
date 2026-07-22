package postgres

import (
	"time"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/repository/postgres/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func mapProject(row any) domain.Project {
	switch r := row.(type) {
	case sqlc.Project:
		return domain.Project{
			ID:        r.ID,
			Code:      r.Code,
			StartDate: fromDate(r.StartDate),
			EndDate:   fromDate(r.EndDate),
			Priority:  int(r.Priority),
		}
	case sqlc.CreateProjectRow:
		return domain.Project{
			ID:        r.ID,
			Code:      r.Code,
			StartDate: fromDate(r.StartDate),
			EndDate:   fromDate(r.EndDate),
			Priority:  int(r.Priority),
		}
	case sqlc.GetProjectRow:
		return domain.Project{
			ID:        r.ID,
			Code:      r.Code,
			StartDate: fromDate(r.StartDate),
			EndDate:   fromDate(r.EndDate),
			Priority:  int(r.Priority),
		}
	case sqlc.ListProjectsRow:
		return domain.Project{
			ID:        r.ID,
			Code:      r.Code,
			StartDate: fromDate(r.StartDate),
			EndDate:   fromDate(r.EndDate),
			Priority:  int(r.Priority),
		}
	case sqlc.UpdateProjectRow:
		return domain.Project{
			ID:        r.ID,
			Code:      r.Code,
			StartDate: fromDate(r.StartDate),
			EndDate:   fromDate(r.EndDate),
			Priority:  int(r.Priority),
		}
	default:
		return domain.Project{}
	}
}

func mapProcess(row any) domain.Process {
	switch r := row.(type) {
	case sqlc.Process:
		return domain.Process{
			ID:        r.ID,
			ProjectID: r.ProjectID,
			Title:     r.Title,
			StartDate: fromDate(r.StartDate),
			EndDate:   fromDate(r.EndDate),
		}
	case sqlc.CreateProcessRow:
		return domain.Process{
			ID:        r.ID,
			ProjectID: r.ProjectID,
			Title:     r.Title,
			StartDate: fromDate(r.StartDate),
			EndDate:   fromDate(r.EndDate),
		}
	case sqlc.GetProcessRow:
		return domain.Process{
			ID:        r.ID,
			ProjectID: r.ProjectID,
			Title:     r.Title,
			StartDate: fromDate(r.StartDate),
			EndDate:   fromDate(r.EndDate),
		}
	case sqlc.ListProcessesByProjectIDRow:
		return domain.Process{
			ID:        r.ID,
			ProjectID: r.ProjectID,
			Title:     r.Title,
			StartDate: fromDate(r.StartDate),
			EndDate:   fromDate(r.EndDate),
		}
	case sqlc.UpdateProcessRow:
		return domain.Process{
			ID:        r.ID,
			ProjectID: r.ProjectID,
			Title:     r.Title,
			StartDate: fromDate(r.StartDate),
			EndDate:   fromDate(r.EndDate),
		}
	default:
		return domain.Process{}
	}
}

func mapMilestone(row any) domain.Milestone {
	switch r := row.(type) {
	case sqlc.Milestone:
		return domain.Milestone{
			ID:        r.ID,
			ProcessID: r.ProcessID,
			Title:     r.Title,
			Content:   r.Content,
			Date:      fromDate(r.Date),
		}
	case sqlc.CreateMilestoneRow:
		return domain.Milestone{
			ID:        r.ID,
			ProcessID: r.ProcessID,
			Title:     r.Title,
			Content:   r.Content,
			Date:      fromDate(r.Date),
		}
	case sqlc.GetMilestoneRow:
		return domain.Milestone{
			ID:        r.ID,
			ProcessID: r.ProcessID,
			Title:     r.Title,
			Content:   r.Content,
			Date:      fromDate(r.Date),
		}
	case sqlc.ListMilestonesByProcessIDRow:
		return domain.Milestone{
			ID:        r.ID,
			ProcessID: r.ProcessID,
			Title:     r.Title,
			Content:   r.Content,
			Date:      fromDate(r.Date),
		}
	case sqlc.UpdateMilestoneRow:
		return domain.Milestone{
			ID:        r.ID,
			ProcessID: r.ProcessID,
			Title:     r.Title,
			Content:   r.Content,
			Date:      fromDate(r.Date),
		}
	default:
		return domain.Milestone{}
	}
}

func mapResource(row any) domain.Resource {
	switch r := row.(type) {
	case sqlc.Resource:
		return domain.Resource{
			ID:       r.ID,
			Title:    r.Title,
			Code:     r.Code,
			Quantity: int(r.Quantity),
		}
	case sqlc.CreateResourceRow:
		return domain.Resource{
			ID:       r.ID,
			Title:    r.Title,
			Code:     r.Code,
			Quantity: int(r.Quantity),
		}
	case sqlc.GetResourceRow:
		return domain.Resource{
			ID:       r.ID,
			Title:    r.Title,
			Code:     r.Code,
			Quantity: int(r.Quantity),
		}
	case sqlc.ListResourcesRow:
		return domain.Resource{
			ID:       r.ID,
			Title:    r.Title,
			Code:     r.Code,
			Quantity: int(r.Quantity),
		}
	case sqlc.UpdateResourceRow:
		return domain.Resource{
			ID:       r.ID,
			Title:    r.Title,
			Code:     r.Code,
			Quantity: int(r.Quantity),
		}
	default:
		return domain.Resource{}
	}
}

func mapTask(row any) domain.Task {
	switch r := row.(type) {
	case sqlc.Task:
		return domain.Task{
			ID:        r.ID,
			ProcessID: r.ProcessID,
			Title:     r.Title,
			StartDate: fromDate(r.StartDate),
			EndDate:   fromDate(r.EndDate),
		}
	case sqlc.CreateTaskRow:
		return domain.Task{
			ID:        r.ID,
			ProcessID: r.ProcessID,
			Title:     r.Title,
			StartDate: fromDate(r.StartDate),
			EndDate:   fromDate(r.EndDate),
		}
	case sqlc.GetTaskRow:
		return domain.Task{
			ID:        r.ID,
			ProcessID: r.ProcessID,
			Title:     r.Title,
			StartDate: fromDate(r.StartDate),
			EndDate:   fromDate(r.EndDate),
		}
	case sqlc.ListTasksByProcessIDRow:
		return domain.Task{
			ID:        r.ID,
			ProcessID: r.ProcessID,
			Title:     r.Title,
			StartDate: fromDate(r.StartDate),
			EndDate:   fromDate(r.EndDate),
		}
	case sqlc.UpdateTaskRow:
		return domain.Task{
			ID:        r.ID,
			ProcessID: r.ProcessID,
			Title:     r.Title,
			StartDate: fromDate(r.StartDate),
			EndDate:   fromDate(r.EndDate),
		}
	default:
		return domain.Task{}
	}
}

func mapAssignment(row any) domain.Assignment {
	switch r := row.(type) {
	case sqlc.Assignment:
		return domain.Assignment{
			ID:         r.ID,
			TaskID:     r.TaskID,
			ResourceID: r.ResourceID,
			Quantity:   int(r.Quantity),
		}
	case sqlc.CreateAssignmentRow:
		return domain.Assignment{
			ID:         r.ID,
			TaskID:     r.TaskID,
			ResourceID: r.ResourceID,
			Quantity:   int(r.Quantity),
		}
	case sqlc.GetAssignmentRow:
		return domain.Assignment{
			ID:         r.ID,
			TaskID:     r.TaskID,
			ResourceID: r.ResourceID,
			Quantity:   int(r.Quantity),
		}
	case sqlc.ListAssignmentsByTaskIDRow:
		return domain.Assignment{
			ID:         r.ID,
			TaskID:     r.TaskID,
			ResourceID: r.ResourceID,
			Quantity:   int(r.Quantity),
		}
	case sqlc.UpdateAssignmentRow:
		return domain.Assignment{
			ID:         r.ID,
			TaskID:     r.TaskID,
			ResourceID: r.ResourceID,
			Quantity:   int(r.Quantity),
		}
	default:
		return domain.Assignment{}
	}
}

func mapUser(row any) domain.User {
	switch r := row.(type) {
	case sqlc.User:
		return domain.User{
			ID:           r.ID,
			Name:         r.Name,
			Role:         domain.UserRole(r.Role),
			Username:     r.Username,
			PasswordHash: r.PasswordHash,
		}
	case sqlc.CreateUserRow:
		return domain.User{
			ID:           r.ID,
			Role:         domain.UserRole(r.Role),
			Username:     r.Username,
			PasswordHash: r.PasswordHash,
		}
	case sqlc.GetUserRow:
		return domain.User{
			ID:           r.ID,
			Role:         domain.UserRole(r.Role),
			Username:     r.Username,
			PasswordHash: r.PasswordHash,
		}
	case sqlc.ListUsersRow:
		return domain.User{
			ID:           r.ID,
			Role:         domain.UserRole(r.Role),
			Username:     r.Username,
			PasswordHash: r.PasswordHash,
		}
	default:
		return domain.User{}
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
