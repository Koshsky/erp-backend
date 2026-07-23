// internal/service/mapper/process_mapper.go
package mapper

import (
	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/dto"
)

type ProcessMapper struct{}

func NewProcessMapper() *ProcessMapper {
	return &ProcessMapper{}
}

func (p *ProcessMapper) ToDetailedDTO(process *domain.DetailedProcess) *dto.ProcessDetailResponse {
	tasks := make([]dto.TaskDetailResponse, len(process.Tasks))
	for i, task := range process.Tasks {
		assignments := make([]dto.AssignmentResponse, len(task.Assignments))
		for j, a := range task.Assignments {
			assignments[j] = dto.AssignmentResponse{
				ID:         a.ID,
				TaskID:     a.TaskID,
				ResourceID: a.ResourceID,
				Quantity:   a.Quantity,
			}
		}

		tasks[i] = dto.TaskDetailResponse{
			TaskResponse: dto.TaskResponse{
				ID:        task.ID,
				ProcessID: task.ProcessID,
				Title:     task.Title,
				StartDate: task.StartDate,
				EndDate:   task.EndDate,
			},
			Assignments: assignments,
		}
	}

	milestones := make([]dto.MilestoneResponse, len(process.Milestones))
	for i, m := range process.Milestones {
		milestones[i] = dto.MilestoneResponse{
			ID:        m.ID,
			Content:   m.Content,
			ProcessID: m.ProcessID,
			Title:     m.Title,
			Date:      m.Date,
		}
	}

	return &dto.ProcessDetailResponse{
		ProcessResponse: dto.ProcessResponse{
			ID:        process.ID,
			OwnerID:   process.OwnerID,
			ProjectID: process.ProjectID,
			Title:     process.Title,
			StartDate: process.StartDate,
			EndDate:   process.EndDate,
		},
		Tasks:      tasks,
		Milestones: milestones,
	}
}

func (m *ProcessMapper) ToDTO(process *domain.Process) *dto.ProcessResponse {
	if process == nil {
		return nil
	}
	return &dto.ProcessResponse{
		ID:        process.ID,
		OwnerID:   process.OwnerID,
		ProjectID: process.ProjectID,
		Title:     process.Title,
		StartDate: process.StartDate,
		EndDate:   process.EndDate,
	}
}

func (m *ProcessMapper) ToDTOs(processes []domain.Process) []dto.ProcessResponse {
	if processes == nil {
		return []dto.ProcessResponse{}
	}

	responses := make([]dto.ProcessResponse, len(processes))
	for i, process := range processes {
		responses[i] = *m.ToDTO(&process)
	}
	return responses
}

func (m *ProcessMapper) ToDomainFromCreate(req dto.CreateProcessRequest) domain.Process {
	return domain.Process{
		ProjectID: req.ProjectID,
		Title:     req.Title,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}
}

func (m *ProcessMapper) ApplyUpdateToDomain(process *domain.Process, req dto.UpdateProcessRequest) {
	if process == nil {
		return
	}

	if req.Title != nil {
		process.Title = *req.Title
	}
	if req.StartDate != nil {
		process.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		process.EndDate = *req.EndDate
	}
}
