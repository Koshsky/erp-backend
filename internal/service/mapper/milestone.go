// internal/mapper/milestone_mapper.go
package mapper

import (
	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/dto"
)

type MilestoneMapper struct{}

func NewMilestoneMapper() *MilestoneMapper {
	return &MilestoneMapper{}
}

func (m *MilestoneMapper) ToDTO(milestone *domain.Milestone) *dto.MilestoneResponse {
	if milestone == nil {
		return nil
	}
	return &dto.MilestoneResponse{
		ID:        milestone.ID,
		Title:     milestone.Title,
		Content:   milestone.Content,
		Date:      milestone.Date,
		ProcessID: milestone.ProcessID,
	}
}

func (m *MilestoneMapper) ToDTOs(milestones []domain.Milestone) []dto.MilestoneResponse {
	if milestones == nil {
		return []dto.MilestoneResponse{}
	}

	responses := make([]dto.MilestoneResponse, len(milestones))
	for i, milestone := range milestones {
		responses[i] = *m.ToDTO(&milestone)
	}
	return responses
}

func (m *MilestoneMapper) ToDomainFromCreate(req dto.CreateMilestoneRequest) domain.Milestone {
	return domain.Milestone{
		Title:     req.Title,
		Content:   req.Content,
		Date:      req.Date,
		ProcessID: req.ProcessID,
	}
}

func (m *MilestoneMapper) ApplyUpdateToDomain(milestone *domain.Milestone, req dto.UpdateMilestoneRequest) {
	if req.Title != nil {
		milestone.Title = *req.Title
	}
	if req.Content != nil {
		milestone.Content = *req.Content
	}
	if req.Date != nil {
		milestone.Date = *req.Date
	}
}
