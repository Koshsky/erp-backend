package service

import (
	"time"

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
		ProcessID: task.ProcessID,
		ParentID:  task.ParentID,
		OwnerID:   task.OwnerID,
		Title:     task.Title,
		Color:     task.Color,
		Status:    task.Status,
		StartDate: date.From(task.StartDate),
		EndDate:   date.From(task.EndDate),
		Order:     task.SortOrder,
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
	status := domain.StatusNotStarted
	if req.Status != nil && *req.Status != "" {
		status = *req.Status
	}
	return domain.Task{
		ProcessID: req.ProcessID,
		ParentID:  req.ParentID,
		OwnerID:   req.OwnerID,
		Title:     req.Title,
		Color:     req.Color,
		Status:    status,
		StartDate: dateFromPtr(req.StartDate),
		EndDate:   dateFromPtr(req.EndDate),
	}
}

// dateFromPtr converts a nullable request date; a nil pointer yields the zero
// time (subtask dates are inherited from the parent by the service).
func dateFromPtr(d *date.Date) time.Time {
	if d == nil {
		return time.Time{}
	}
	return d.Time()
}

func (m *TaskMapper) ApplyUpdateToDomain(task *domain.Task, req dto.UpdateTaskRequest) {
	if task == nil {
		return
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Color != nil {
		if *req.Color == "" {
			task.Color = nil
		} else {
			task.Color = req.Color
		}
	}
	if req.OwnerID != nil {
		task.OwnerID = req.OwnerID
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.StartDate != nil {
		task.StartDate = req.StartDate.Time()
	}
	if req.EndDate != nil {
		task.EndDate = req.EndDate.Time()
	}
}
