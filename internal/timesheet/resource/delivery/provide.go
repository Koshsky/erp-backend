package delivery

import (
	"log/slog"

	rbac "github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/timesheet/resource/service"
)

// ProvideResourceHandler builds the ResourceHandler handler.
func ProvideResourceHandler(logger *slog.Logger, svc *service.ResourceService, mw *rbac.Middleware) *ResourceHandler {
	return &ResourceHandler{
		logger:  logger,
		service: svc,
		mw:      mw,
	}
}
