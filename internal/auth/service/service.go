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

// activeSessionSweepWindow — how long expired sessions are kept before background cleanup.
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
		// Reusing a revoked token indicates theft: revoke
		// all of the user's active sessions.
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

	// Rotation: the old session is revoked and a new pair is issued.
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

// Logout revokes the session by refresh token (idempotently).
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	ctx, end := s.tracer.Start(ctx, "auth.Logout")
	defer end(nil)

	if refreshToken == "" {
		return nil
	}
	session, err := s.findSession(ctx, refreshToken)
	if err != nil {
		// Unknown/already revoked token — logout is idempotent, this is not an error.
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

// issueSession creates an opaque refresh token and stores its SHA-256 hash in the DB.
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

// sweepExpired — background cleanup of long-expired sessions (best-effort).
func (s *AuthService) sweepExpired(ctx context.Context) {
	_ = s.sessions.DeleteExpiredSessions(ctx, time.Now().Add(-activeSessionSweepWindow))
}

// generateRefreshToken — an opaque 256-bit token (crypto/rand, hex).
func generateRefreshToken() (string, error) {
	var buf [32]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf[:]), nil
}

// hashToken — SHA-256 of the token; only the hash is stored in the DB.
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
