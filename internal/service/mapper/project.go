// internal/service/mapper/project_mapper.go
package mapper

import (
	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/dto"
)

type ProjectMapper struct{}

func NewProjectMapper() *ProjectMapper {
	return &ProjectMapper{}
}

func (m *ProjectMapper) ToDTO(project *domain.Project) *dto.ProjectResponse {
	if project == nil {
		return nil
	}
	return &dto.ProjectResponse{
		ID:        project.ID,
		Code:      project.Code,
		StartDate: project.StartDate,
		EndDate:   project.EndDate,
		Priority:  project.Priority,
	}
}

func (m *ProjectMapper) ToDTOs(projects []domain.Project) []dto.ProjectResponse {
	if projects == nil {
		return []dto.ProjectResponse{}
	}

	responses := make([]dto.ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = *m.ToDTO(&project)
	}
	return responses
}

func (m *ProjectMapper) ToDomainFromCreate(req dto.CreateProjectRequest) domain.Project {
	return domain.Project{
		Code:      req.Code,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Priority:  req.Priority,
	}
}

func (m *ProjectMapper) ApplyUpdateToDomain(project *domain.Project, req dto.UpdateProjectRequest) {
	if project == nil {
		return
	}

	if req.Code != nil {
		project.Code = *req.Code
	}
	if req.StartDate != nil {
		project.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		project.EndDate = *req.EndDate
	}
	if req.Priority != nil {
		project.Priority = *req.Priority
	}
}

func (m *ProjectMapper) ToDetailedDTO(project *domain.DetailedProject) *dto.ProjectDetailResponse {
	if project == nil {
		return nil
	}

	processes := make([]dto.ProcessDetailResponse, len(project.Processes))
	for i, p := range project.Processes {
		processes[i] = m.toDetailedProcessDTO(p)
	}

	return &dto.ProjectDetailResponse{
		ProjectResponse: dto.ProjectResponse{
			ID:        project.ID,
			Code:      project.Code,
			OwnerID:   project.OwnerID,
			StartDate: project.StartDate,
			EndDate:   project.EndDate,
			Priority:  project.Priority,
		},
		Processes: processes,
	}
}

func (m *ProjectMapper) toDetailedProcessDTO(p domain.DetailedProcess) dto.ProcessDetailResponse {
	// Маппим задачи
	tasks := make([]dto.TaskDetailResponse, len(p.Tasks))
	for i, t := range p.Tasks {
		tasks[i] = m.toDetailedTaskDTO(t)
	}

	// Маппим милейстоуны
	milestones := make([]dto.MilestoneResponse, len(p.Milestones))
	for i, ml := range p.Milestones {
		milestones[i] = dto.MilestoneResponse{
			ID:        ml.ID,
			ProcessID: ml.ProcessID,
			Content:   ml.Content,
			Title:     ml.Title,
			Date:      ml.Date,
		}
	}

	return dto.ProcessDetailResponse{
		ProcessResponse: dto.ProcessResponse{
			ID:        p.ID,
			OwnerID:   p.OwnerID,
			ProjectID: p.ProjectID,
			Title:     p.Title,
			StartDate: p.StartDate,
			EndDate:   p.EndDate,
		},
		Tasks:      tasks,
		Milestones: milestones,
	}
}

// toDetailedTaskDTO - маппинг domain.DetailedTask → dto.DetailedTaskResponse
func (m *ProjectMapper) toDetailedTaskDTO(t domain.DetailedTask) dto.TaskDetailResponse {
	// Маппим ассайнменты
	assignments := make([]dto.AssignmentResponse, len(t.Assignments))
	for i, a := range t.Assignments {
		assignments[i] = dto.AssignmentResponse{
			ID:         a.ID,
			TaskID:     a.TaskID,
			ResourceID: a.ResourceID,
			Quantity:   a.Quantity,
		}
	}

	return dto.TaskDetailResponse{
		TaskResponse: dto.TaskResponse{
			ID:        t.ID,
			ProcessID: t.ProcessID,
			Title:     t.Title,
			StartDate: t.StartDate,
			EndDate:   t.EndDate,
		},
		Assignments: assignments,
	}
}
