package service

import (
	"context"
	"fmt"

	"github.com/Koshsky/erp-backend/internal/auth/dto"
	"github.com/Koshsky/erp-backend/internal/security/jwt"
	userDTO "github.com/Koshsky/erp-backend/internal/user/dto"
)

type AuthService struct {
	users  UserService
	hasher PasswordHasher
	jwt    *jwt.JWTService
}

func NewAuthService(users UserService, hasher PasswordHasher, jwtService *jwt.JWTService) *AuthService {
	return &AuthService{
		users:  users,
		hasher: hasher,
		jwt:    jwtService,
	}
}

func (s *AuthService) Register(ctx context.Context, name, username, password string) (*dto.AuthResponse, error) {
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}

	user, err := s.users.CreateUser(ctx, userDTO.CreateUserRequest{
		Name:         name,
		Username:     username,
		Role:         "ВП",
		PasswordHash: hash,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	tokens, err := s.jwt.GenerateTokenPair(user.ID, string(user.Role), user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens")
	}

	return &dto.AuthResponse{
		User: dto.UserInfo{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Role:     string(user.Role),
		},
		Tokens: tokens,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*dto.AuthResponse, error) {
	user, err := s.users.FindUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := s.hasher.Compare(user.PasswordHash, password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	tokens, err := s.jwt.GenerateTokenPair(user.ID, string(user.Role), user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens")
	}

	return &dto.AuthResponse{
		User: dto.UserInfo{
			ID:       user.ID,
			Name:     user.Name,
			Username: user.Username,
			Role:     string(user.Role),
		},
		Tokens: tokens,
	}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	user, err := s.users.FindUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if err := s.hasher.Compare(user.PasswordHash, oldPassword); err != nil {
		return fmt.Errorf("invalid current password")
	}

	newHash, err := s.hasher.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	return s.users.UpdatePassword(ctx, userID, newHash)
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	claims, err := s.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}

	var userID int64
	if _, err := fmt.Sscanf(claims.Subject, "%d", &userID); err != nil {
		return nil, fmt.Errorf("invalid token subject")
	}

	user, err := s.users.FindUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	tokens, err := s.jwt.GenerateTokenPair(userID, string(user.Role), claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens")
	}

	return &dto.RefreshResponse{
		Tokens:  tokens,
		Message: "Token refreshed successfully",
	}, nil
}
