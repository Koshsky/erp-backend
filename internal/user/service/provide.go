package service

import (
	"log/slog"

	repo "github.com/Koshsky/erp-backend/internal/user/repository"
)

// ProvideUserService builds the UserService service.
func ProvideUserService(logger *slog.Logger, r *repo.UserRepository) *UserService {
	return &UserService{
		logger:     logger,
		repository: r,
		mapper:     &UserMapper{},
		validator:  &UserValidator{},
	}
}
