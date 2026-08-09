package service

import (
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/planning/repository"
)

// ProvidePlanningService builds the PlanningService service.
func ProvidePlanningService(logger *slog.Logger, r *repo.PlanningRepository) *PlanningService {
	return &PlanningService{
		logger:     logger,
		repository: r,
	}
}
