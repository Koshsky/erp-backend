package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	userservice "github.com/Koshsky/erp-backend/internal/user/service"

	"github.com/Koshsky/erp-backend/internal/auth/dto"
	"github.com/Koshsky/erp-backend/internal/auth/repository"
	"github.com/Koshsky/erp-backend/internal/security/hasher"
	"github.com/Koshsky/erp-backend/internal/security/jwt"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"
	userDTO "github.com/Koshsky/erp-backend/internal/user/dto"
)

// activeSessionSweepWindow — насколько давно истёкшие сессии вычищаем фоново.
const activeSessionSweepWindow = 30 * 24 * time.Hour

type AuthService struct {
	users    UserService
	jwt      *jwt.Service
	sessions *repository.AuthRepository
	tracer   *tracingpkg.Tracer
}

// NewAuthService builds the auth service.
func NewAuthService(
	users *userservice.UserService,
	jwtService *jwt.Service,
	sessions *repository.AuthRepository,
	tracer *tracingpkg.Tracer,
) *AuthService {
	return &AuthService{
		users:    users,
		jwt:      jwtService,
		sessions: sessions,
		tracer:   tracer,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*dto.SessionResult, error) {
	ctx, end := s.tracer.Start(ctx, "auth.Login")
	defer end(nil)

	username = strings.TrimSpace(username)

	user, err := s.users.FindUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if err = hasher.Compare(user.PasswordHash, password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	refresh, err := s.issueSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	access, err := s.jwt.GenerateAccessToken(user.ID, user.Role, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token")
	}

	return &dto.SessionResult{
		Auth:         s.newAuthResponse(user, access),
		RefreshToken: refresh,
	}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.SessionResult, error) {
	ctx, end := s.tracer.Start(ctx, "auth.RefreshToken")
	defer end(nil)

	session, err := s.findSession(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if session.RevokedAt != nil {
		// Повторное использование отозванного токена — признак кражи: сбрасываем
		// все активные сессии пользователя.
		_ = s.sessions.RevokeAllUserSessions(ctx, session.UserID)
		return nil, fmt.Errorf("invalid refresh token")
	}
	if !session.ExpiresAt.After(time.Now()) {
		return nil, fmt.Errorf("invalid refresh token")
	}

	user, err := s.users.FindUserByID(ctx, session.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Ротация: старая сессия отзывается, выписывается новая пара.
	if err = s.sessions.RevokeSession(ctx, session.ID); err != nil {
		return nil, fmt.Errorf("failed to rotate session")
	}
	refresh, err := s.issueSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	access, err := s.jwt.GenerateAccessToken(user.ID, user.Role, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token")
	}
	s.sweepExpired(ctx)

	return &dto.SessionResult{
		Auth:         s.newAuthResponse(user, access),
		RefreshToken: refresh,
	}, nil
}

// Logout отзывает сессию по refresh-токену (идемпотентно).
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	ctx, end := s.tracer.Start(ctx, "auth.Logout")
	defer end(nil)

	if refreshToken == "" {
		return nil
	}
	session, err := s.findSession(ctx, refreshToken)
	if err != nil {
		// Неизвестный/уже отозванный токен — logout идемпотентен, это не ошибка.
		return nil
	}
	if err = s.sessions.RevokeSession(ctx, session.ID); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) findSession(ctx context.Context, refreshToken string) (repository.Session, error) {
	return s.sessions.FindSessionByHash(ctx, hashToken(refreshToken))
}

// issueSession создаёт opaque refresh-токен и сохраняет его SHA-256 хэш в БД.
func (s *AuthService) issueSession(ctx context.Context, userID int64) (string, error) {
	refresh, err := generateRefreshToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token")
	}
	_, err = s.sessions.CreateSession(
		ctx,
		userID,
		hashToken(refresh),
		time.Now().Add(s.jwt.RefreshExpiry()),
	)
	if err != nil {
		return "", fmt.Errorf("failed to create session")
	}
	return refresh, nil
}

// sweepExpired — фоновая очистка давно истёкших сессий (best-effort).
func (s *AuthService) sweepExpired(ctx context.Context) {
	_ = s.sessions.DeleteExpiredSessions(ctx, time.Now().Add(-activeSessionSweepWindow))
}

// generateRefreshToken — opaque 256-битный токен (crypto/rand, hex).
func generateRefreshToken() (string, error) {
	var buf [32]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf[:]), nil
}

// hashToken — SHA-256 токена; в БД хранится только хэш.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *AuthService) newAuthResponse(user *userDTO.UserResponse, access string) *dto.AuthResponse {
	return &dto.AuthResponse{
		AccessToken: access,
		TokenType:   "Bearer",
		ExpiresIn:   int(s.jwt.AccessExpiry().Seconds()),
		User: dto.UserInfo{
			ID:         user.ID,
			Name:       user.Name,
			LastName:   user.LastName,
			FirstName:  user.FirstName,
			MiddleName: user.MiddleName,
			Username:   user.Username,
			Role:       user.Role,
		},
	}
}
