package delivery

import (
	"log/slog"

	planningservice "github.com/Koshsky/erp-backend/internal/planning/service"
)

// ProvidePlanningHandler builds the planning handler.
func ProvidePlanningHandler(logger *slog.Logger, svc *planningservice.PlanningService) *PlanningHandler {
	return &PlanningHandler{
		logger:  logger,
		service: svc,
	}
}
