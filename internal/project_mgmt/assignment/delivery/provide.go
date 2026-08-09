package delivery

import (
	"log/slog"

	rbac "github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/assignment/service"
)

// ProvideAssignmentHandler builds the AssignmentHandler handler.
func ProvideAssignmentHandler(logger *slog.Logger, svc *service.AssignmentService, mw *rbac.Middleware) *AssignmentHandler {
	return &AssignmentHandler{
		logger:  logger,
		service: svc,
		mw:      mw,
	}
}
