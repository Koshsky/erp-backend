package service

import (
	"github.com/Koshsky/erp-backend/internal/planning/dto"
	"github.com/Koshsky/erp-backend/internal/planning/repository/sqlc"
	nullable "github.com/Koshsky/erp-backend/pkg/database"
	"github.com/Koshsky/erp-backend/pkg/date"
)

// toProject converts a project row (planning read model).
func toProject(p sqlc.Project) dto.Project {
	return dto.Project{
		ID:        p.ID,
		OwnerID:   nullable.Int64Ptr(p.OwnerID),
		Code:      p.Code,
		Color:     nullable.StringPtr(p.Color),
		StartDate: date.From(p.StartDate),
		EndDate:   date.From(p.EndDate),
		Priority:  int(p.Priority),
	}
}

// toProcess converts a process row with its parent project code (from an
// embed query that also selects the project code).
func toProcess(p sqlc.ListProcessesRow) dto.Process {
	return dto.Process{
		ID:          p.Process.ID,
		Title:       p.Process.Title,
		Color:       nullable.StringPtr(p.Process.Color),
		OwnerID:     nullable.Int64Ptr(p.Process.OwnerID),
		ProjectID:   p.Process.ProjectID,
		ProjectCode: p.ProjectCode,
		StartDate:   date.From(p.Process.StartDate),
		EndDate:     date.From(p.Process.EndDate),
		Order:       int(p.Process.SortOrder),
	}
}

// toTaskScopeProcess converts a process row from the TASK-scope embed query.
func toTaskScopeProcess(p sqlc.ListProcessesByTaskScopeRow) dto.Process {
	return dto.Process{
		ID:          p.Process.ID,
		Title:       p.Process.Title,
		Color:       nullable.StringPtr(p.Process.Color),
		OwnerID:     nullable.Int64Ptr(p.Process.OwnerID),
		ProjectID:   p.Process.ProjectID,
		ProjectCode: p.ProjectCode,
		StartDate:   date.From(p.Process.StartDate),
		EndDate:     date.From(p.Process.EndDate),
		Order:       int(p.Process.SortOrder),
	}
}

// toTask converts a task row (planning read model).
func toTask(t sqlc.Task) dto.Task {
	return dto.Task{
		ID:        t.ID,
		ProcessID: t.ProcessID,
		ParentID:  nullable.Int64Ptr(t.ParentID),
		OwnerID:   nullable.Int64Ptr(t.OwnerID),
		Title:     t.Title,
		Color:     nullable.StringPtr(t.Color),
		Status:    t.Status,
		StartDate: date.From(t.StartDate),
		EndDate:   date.From(t.EndDate),
		Order:     int(t.SortOrder),
	}
}

// toMilestone converts a milestone row (planning read model).
func toMilestone(m sqlc.Milestone) dto.Milestone {
	return dto.Milestone{
		ID:        m.ID,
		ProcessID: m.ProcessID,
		Title:     m.Title,
		Content:   m.Content,
		Color:     nullable.StringPtr(m.Color),
		Date:      date.From(m.Date),
	}
}

// toAssignment converts an assignment row (planning read model).
func toAssignment(a sqlc.Assignment) dto.Assignment {
	return dto.Assignment{
		ID:         a.ID,
		TaskID:     a.TaskID,
		ResourceID: a.ResourceID,
		Quantity:   int(a.Quantity),
	}
}

// toResource converts a resource row (planning read model).
func toResource(r sqlc.Resource) dto.Resource {
	return dto.Resource{
		ID:    r.ID,
		Title: r.Title,
		Code:  r.Code,
		Color: nullable.StringPtr(r.Color),
	}
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
