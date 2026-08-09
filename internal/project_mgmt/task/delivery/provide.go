package delivery

import (
	"log/slog"

	rbac "github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/project_mgmt/task/service"
)

// ProvideTaskHandler builds the TaskHandler handler.
func ProvideTaskHandler(logger *slog.Logger, svc *service.TaskService, mw *rbac.Middleware) *TaskHandler {
	return &TaskHandler{
		logger:  logger,
		service: svc,
		mw:      mw,
	}
}
