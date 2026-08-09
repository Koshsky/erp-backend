package delivery

import (
	"log/slog"

	rbac "github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/process/service"
)

// ProvideProcessHandler builds the ProcessHandler handler.
func ProvideProcessHandler(logger *slog.Logger, svc *service.ProcessService, mw *rbac.Middleware) *ProcessHandler {
	return &ProcessHandler{
		logger:  logger,
		service: svc,
		mw:      mw,
	}
}
