package delivery

import (
	"log/slog"

	userservice "github.com/Koshsky/erp-backend/internal/user/service"
)

// ProvideUserHandler builds the user handler.
func ProvideUserHandler(logger *slog.Logger, svc *userservice.UserService) *UserHandler {
	return &UserHandler{
		logger:  logger,
		service: svc,
	}
}
