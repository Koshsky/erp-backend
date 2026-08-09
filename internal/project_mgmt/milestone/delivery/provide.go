package delivery

import (
	"log/slog"

	rbac "github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/milestone/service"
)

// ProvideMilestoneHandler builds the MilestoneHandler handler.
func ProvideMilestoneHandler(logger *slog.Logger, svc *service.MilestoneService, mw *rbac.Middleware) *MilestoneHandler {
	return &MilestoneHandler{
		logger:  logger,
		service: svc,
		mw:      mw,
	}
}
