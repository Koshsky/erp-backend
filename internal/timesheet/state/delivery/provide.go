package delivery

import (
	"log/slog"

	rbac "github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/timesheet/state/service"
)

// ProvideStateHandler builds the StateHandler handler.
func ProvideStateHandler(logger *slog.Logger, svc *service.StateService, mw *rbac.Middleware) *StateHandler {
	return &StateHandler{
		logger:  logger,
		service: svc,
		mw:      mw,
	}
}
