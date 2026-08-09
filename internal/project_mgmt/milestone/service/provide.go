package service

import (
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/repository"
)

// ProvideMilestoneService builds the MilestoneService service.
func ProvideMilestoneService(logger *slog.Logger, r *repo.MilestoneRepository) *MilestoneService {
	return &MilestoneService{
		logger:     logger,
		repository: r,
		mapper:     NewMilestoneMapper(),
		validator:  &MilestoneValidator{},
	}
}
