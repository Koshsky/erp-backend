package service

import (
	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/domain"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/dto"
	"github.com/Koshsky/erp-backend/pkg/date"
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
		OwnerID:   project.OwnerID,
		Code:      project.Code,
		Color:     project.Color,
		StartDate: date.From(project.StartDate),
		EndDate:   date.From(project.EndDate),
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

// ToCreateDTO builds the create response: the project plus what the
// auto-create template added to it (counts fetched after the insert).
func (m *ProjectMapper) ToCreateDTO(
	project *domain.Project,
	counts domain.AutoCreatedCounts,
) *dto.CreateProjectResponse {
	if project == nil {
		return nil
	}
	return &dto.CreateProjectResponse{
		ID:        project.ID,
		OwnerID:   project.OwnerID,
		Code:      project.Code,
		Color:     project.Color,
		StartDate: date.From(project.StartDate),
		EndDate:   date.From(project.EndDate),
		Priority:  project.Priority,
		AutoCreated: dto.AutoCreatedCounts{
			Processes:   counts.Processes,
			Tasks:       counts.Tasks,
			Assignments: counts.Assignments,
		},
	}
}

func (m *ProjectMapper) ToDomainFromCreate(req dto.CreateProjectRequest) domain.Project {
	return domain.Project{
		OwnerID:   req.OwnerID,
		Code:      req.Code,
		Color:     req.Color,
		StartDate: req.StartDate.Time(),
		EndDate:   req.EndDate.Time(),
		Priority:  req.Priority,
	}
}

func (m *ProjectMapper) ApplyUpdateToDomain(project *domain.Project, req dto.UpdateProjectRequest) {
	if project == nil {
		return
	}

	if req.OwnerID != nil {
		project.OwnerID = req.OwnerID
	}
	if req.Code != nil {
		project.Code = *req.Code
	}
	if req.Color != nil {
		if *req.Color == "" {
			project.Color = nil
		} else {
			project.Color = req.Color
		}
	}
	if req.StartDate != nil {
		project.StartDate = req.StartDate.Time()
	}
	if req.EndDate != nil {
		project.EndDate = req.EndDate.Time()
	}
	if req.Priority != nil {
		project.Priority = *req.Priority
	}
}
