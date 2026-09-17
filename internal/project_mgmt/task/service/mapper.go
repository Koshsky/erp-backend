package service

import (
	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/dto"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/repository/sqlc"
	nullable "github.com/Koshsky/erp-backend/pkg/database"
	"github.com/Koshsky/erp-backend/pkg/date"
)

type TaskMapper struct{}

func NewTaskMapper() *TaskMapper {
	return &TaskMapper{}
}

func (m *TaskMapper) ToDTO(task *sqlc.Task) *dto.TaskResponse {
	if task == nil {
		return nil
	}
	return &dto.TaskResponse{
		ID:        task.ID,
		ProcessID: task.ProcessID,
		ParentID:  nullable.Int64Ptr(task.ParentID),
		OwnerID:   nullable.Int64Ptr(task.OwnerID),
		Title:     task.Title,
		Color:     nullable.StringPtr(task.Color),
		Status:    task.Status,
		StartDate: date.From(task.StartDate),
		EndDate:   date.From(task.EndDate),
		Order:     int(task.SortOrder),
	}
}

func (m *TaskMapper) ToDTOs(tasks []sqlc.Task) []dto.TaskResponse {
	if tasks == nil {
		return []dto.TaskResponse{}
	}

	responses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = *m.ToDTO(&task)
	}
	return responses
}
