package delivery

import (
	"log/slog"

	rbac "github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/project/service"
)

// ProvideProjectHandler builds the ProjectHandler handler.
func ProvideProjectHandler(logger *slog.Logger, svc *service.ProjectService, mw *rbac.Middleware) *ProjectHandler {
	return &ProjectHandler{
		logger:  logger,
		service: svc,
		mw:      mw,
	}
}
