package service

import (
	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/domain"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/dto"
	"github.com/Koshsky/erp-backend/pkg/date"
)

type TaskMapper struct{}

func NewTaskMapper() *TaskMapper {
	return &TaskMapper{}
}

func (m *TaskMapper) ToDTO(task *domain.Task) *dto.TaskResponse {
	if task == nil {
		return nil
	}
	return &dto.TaskResponse{
		ID:        task.ID,
		OwnerID:   task.OwnerID,
		Title:     task.Title,
		StartDate: date.From(task.StartDate),
		EndDate:   date.From(task.EndDate),
		ProcessID: task.ProcessID,
	}
}

func (m *TaskMapper) ToDTOs(tasks []domain.Task) []dto.TaskResponse {
	if tasks == nil {
		return []dto.TaskResponse{}
	}

	responses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = *m.ToDTO(&task)
	}
	return responses
}

func (m *TaskMapper) ToDomainFromCreate(req dto.CreateTaskRequest) domain.Task {
	return domain.Task{
		OwnerID:   req.OwnerID,
		Title:     req.Title,
		ProcessID: req.ProcessID,
		StartDate: req.StartDate.Time(),
		EndDate:   req.EndDate.Time(),
	}
}

func (m *TaskMapper) ApplyUpdateToDomain(task *domain.Task, req dto.UpdateTaskRequest) {
	if task == nil {
		return
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.OwnerID != nil {
		task.OwnerID = req.OwnerID
	}
	if req.ProcessID != nil {
		task.ProcessID = *req.ProcessID
	}
	if req.StartDate != nil {
		task.StartDate = req.StartDate.Time()
	}
	if req.EndDate != nil {
		task.EndDate = req.EndDate.Time()
	}
}
