package service

import (
	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/domain"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/dto"
	"github.com/Koshsky/erp-backend/pkg/date"
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
		StartDate: date.From(process.StartDate),
		EndDate:   date.From(process.EndDate),
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
		StartDate: req.StartDate.Time(),
		EndDate:   req.EndDate.Time(),
		OwnerID:   req.OwnerID,
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
		process.StartDate = req.StartDate.Time()
	}
	if req.EndDate != nil {
		process.EndDate = req.EndDate.Time()
	}
	if req.OwnerID != nil {
		process.OwnerID = req.OwnerID
	}
	if req.ProjectID != nil {
		process.ProjectID = *req.ProjectID
	}
}
