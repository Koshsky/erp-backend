package service

import (
	"context"
	"fmt"

	"github.com/Koshsky/erp-backend/internal/auth/dto"
	"github.com/Koshsky/erp-backend/internal/security/jwt"
)

type AuthService struct {
	repo   AuthRepository
	hasher PasswordHasher
	jwt    *jwt.JWTService
}

func NewAuthService(repo AuthRepository, hasher PasswordHasher, jwtService *jwt.JWTService) *AuthService {
	return &AuthService{
		repo:   repo,
		hasher: hasher,
		jwt:    jwtService,
	}
}

func (s *AuthService) Register(ctx context.Context, name, username, password, role string) (*dto.AuthResponse, error) {
	if role == "" {
		role = "user"
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}

	userID, err := s.repo.CreateUser(ctx, name, username, role, hash)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	tokens, err := s.jwt.GenerateTokenPair(userID, role, username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens")
	}

	return &dto.AuthResponse{
		User: dto.UserInfo{
			ID:       userID,
			Name:     name,
			Username: username,
			Role:     role,
		},
		Tokens: tokens,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*dto.AuthResponse, error) {
	userID, name, role, passwordHash, err := s.repo.FindUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err := s.hasher.Compare(passwordHash, password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	tokens, err := s.jwt.GenerateTokenPair(userID, role, username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens")
	}

	return &dto.AuthResponse{
		User: dto.UserInfo{
			ID:       userID,
			Name:     name,
			Username: username,
			Role:     role,
		},
		Tokens: tokens,
	}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	_, _, _, passwordHash, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if err := s.hasher.Compare(passwordHash, oldPassword); err != nil {
		return fmt.Errorf("invalid current password")
	}

	newHash, err := s.hasher.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	return s.repo.UpdatePassword(ctx, userID, newHash)
}

func (s *AuthService) Logout(ctx context.Context, userID int64) error {
	return s.repo.DeleteRefreshToken(ctx, userID)
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

	_, _, role, _, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	tokens, err := s.jwt.GenerateTokenPair(userID, role, claims.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens")
	}

	return &dto.RefreshResponse{
		Tokens:  tokens,
		Message: "Token refreshed successfully",
	}, nil
}
