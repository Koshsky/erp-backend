package service

import (
	"github.com/Koshsky/erp-backend/internal/security/jwt"
	userservice "github.com/Koshsky/erp-backend/internal/user/service"
)

// ProvideAuthService builds the auth service.
func ProvideAuthService(users *userservice.UserService, jwtService *jwt.Service) *AuthService {
	return &AuthService{
		users: users,
		jwt:   jwtService,
	}
}
