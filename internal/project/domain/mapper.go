package domain

import (
	"github.com/Koshsky/erp/api/internal/dto"
)

type ProjectMapper struct{}

func NewProjectMapper() *ProjectMapper {
	return &ProjectMapper{}
}

func (m *ProjectMapper) ToDTO(project *Project) *dto.ProjectResponse {
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

func (m *ProjectMapper) ToDTOs(projects []Project) []dto.ProjectResponse {
	if projects == nil {
		return []dto.ProjectResponse{}
	}

	responses := make([]dto.ProjectResponse, len(projects))
	for i, project := range projects {
		responses[i] = *m.ToDTO(&project)
	}
	return responses
}

func (m *ProjectMapper) ToDomainFromCreate(req dto.CreateProjectRequest) Project {
	return Project{
		Code:      req.Code,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Priority:  req.Priority,
	}
}

func (m *ProjectMapper) ApplyUpdateToDomain(project *Project, req dto.UpdateProjectRequest) {
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