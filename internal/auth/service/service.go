package service

import (
	"context"
	"fmt"
	"strings"

	userservice "github.com/Koshsky/erp-backend/internal/user/service"

	"github.com/Koshsky/erp-backend/internal/auth/dto"
	"github.com/Koshsky/erp-backend/internal/security/hasher"
	"github.com/Koshsky/erp-backend/internal/security/jwt"
	"github.com/Koshsky/erp-backend/internal/user/domain"
	userDTO "github.com/Koshsky/erp-backend/internal/user/dto"
)

type AuthService struct {
	users UserService
	jwt   *jwt.Service
}

// NewAuthService builds the auth service.
func NewAuthService(users *userservice.UserService, jwtService *jwt.Service) *AuthService {
	return &AuthService{
		users: users,
		jwt:   jwtService,
	}
}

func (s *AuthService) Register(ctx context.Context, name, username, password string) (*dto.AuthResponse, error) {
	username = strings.TrimSpace(username)

	hash, err := hasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}

	user, err := s.users.CreateUser(ctx, userDTO.CreateUserRequest{
		Name:         name,
		Username:     username,
		Role:         domain.Worker,
		PasswordHash: hash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	tokens, err := s.jwt.GenerateTokenPair(user.ID, user.Role, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens")
	}

	return newAuthResponse(user, tokens), nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*dto.AuthResponse, error) {
	username = strings.TrimSpace(username)

	user, err := s.users.FindUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err = hasher.Compare(user.PasswordHash, password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	tokens, err := s.jwt.GenerateTokenPair(user.ID, user.Role, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens")
	}

	return newAuthResponse(user, tokens), nil
}

// newAuthResponse builds the flattened login/register payload.
func newAuthResponse(user *userDTO.UserResponse, tokens *jwt.TokenPair) *dto.AuthResponse {
	return &dto.AuthResponse{
		AccessToken:  tokens.AccessToken,
		TokenType:    tokens.TokenType,
		ExpiresIn:    tokens.ExpiresIn,
		RefreshToken: tokens.RefreshToken,
		User: dto.UserInfo{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Role:     user.Role,
		},
	}
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	claims, err := s.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	var userID int64
	if _, err = fmt.Sscanf(claims.Subject, "%d", &userID); err != nil {
		return nil, fmt.Errorf("invalid token subject")
	}

	user, err := s.users.FindUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	tokens, err := s.jwt.GenerateTokenPair(userID, user.Role, claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens")
	}

	return &dto.RefreshResponse{
		Tokens:  tokens,
		Message: "Token refreshed successfully",
	}, nil
}
