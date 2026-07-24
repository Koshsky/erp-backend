package domain

import (
	"github.com/Koshsky/erp/api/internal/dto"
)

type TaskMapper struct{}

func NewTaskMapper() *TaskMapper {
	return &TaskMapper{}
}

func (m *TaskMapper) ToDTO(task *Task) *dto.TaskResponse {
	if task == nil {
		return nil
	}
	return &dto.TaskResponse{
		ID:        task.ID,
		Title:     task.Title,
		StartDate: task.StartDate,
		EndDate:   task.EndDate,
		ProcessID: task.ProcessID,
	}
}

func (m *TaskMapper) ToDTOs(tasks []Task) []dto.TaskResponse {
	if tasks == nil {
		return []dto.TaskResponse{}
	}

	responses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		responses[i] = *m.ToDTO(&task)
	}
	return responses
}

func (m *TaskMapper) ToDomainFromCreate(req dto.CreateTaskRequest) Task {
	return Task{
		Title:     req.Title,
		ProcessID: req.ProcessID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}
}

func (m *TaskMapper) ApplyUpdateToDomain(task *Task, req dto.UpdateTaskRequest) {
	if task == nil {
		return
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.StartDate != nil {
		task.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		task.EndDate = *req.EndDate
	}
}
