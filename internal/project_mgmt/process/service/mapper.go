package service

import (
	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/domain"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/dto"
)

type ProcessMapper struct{}

func NewProcessMapper() *ProcessMapper {
	return &ProcessMapper{}
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
