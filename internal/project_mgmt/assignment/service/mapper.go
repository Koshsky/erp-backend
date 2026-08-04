package service

import (
	"github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/domain"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/dto"
)

type AssignmentMapper struct{}

func NewAssignmentMapper() *AssignmentMapper {
	return &AssignmentMapper{}
}

func (m *AssignmentMapper) ToDTO(assignment *domain.Assignment) *dto.AssignmentResponse {
	if assignment == nil {
		return nil
	}
	return &dto.AssignmentResponse{
		ID:         assignment.ID,
		TaskID:     assignment.TaskID,
		ResourceID: assignment.ResourceID,
		Quantity:   assignment.Quantity,
	}
}

func (m *AssignmentMapper) ToDTOs(assignments []domain.Assignment) []dto.AssignmentResponse {
	if assignments == nil {
		return []dto.AssignmentResponse{}
	}
	responses := make([]dto.AssignmentResponse, len(assignments))
	for i, assignment := range assignments {
		responses[i] = *m.ToDTO(&assignment)
	}
	return responses
}

func (m *AssignmentMapper) ToDomainFromCreate(req dto.CreateAssignmentRequest) domain.Assignment {
	return domain.Assignment{
		TaskID:     req.TaskID,
		ResourceID: req.ResourceID,
		Quantity:   req.Quantity,
	}
}

func (m *AssignmentMapper) ApplyUpdateToDomain(assignment *domain.Assignment, req dto.UpdateAssignmentRequest) {
	if req.TaskID != nil {
		assignment.TaskID = *req.TaskID
	}
	if req.ResourceID != nil {
		assignment.ResourceID = *req.ResourceID
	}
	if req.Quantity != nil {
		assignment.Quantity = *req.Quantity
	}
}
