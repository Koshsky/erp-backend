package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	userservice "github.com/Koshsky/erp-backend/internal/user/service"

	"github.com/Koshsky/erp-backend/internal/auth/dto"
	repo "github.com/Koshsky/erp-backend/internal/auth/repository"
	"github.com/Koshsky/erp-backend/internal/auth/repository/sqlc"
	"github.com/Koshsky/erp-backend/internal/security/hasher"
	"github.com/Koshsky/erp-backend/internal/security/jwt"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"
	userDTO "github.com/Koshsky/erp-backend/internal/user/dto"
)

// activeSessionSweepWindow — how long expired sessions are kept before background cleanup.
const activeSessionSweepWindow = 30 * 24 * time.Hour

// sessionSweepInterval — how often the background expired-session sweep runs.
const sessionSweepInterval = 1 * time.Hour

type AuthService struct {
	logger   *slog.Logger
	users    UserService
	jwt      *jwt.Service
	sessions *repo.AuthRepository
	tracer   *tracingpkg.Tracer

	// cleanupOnce guarantees the expired-session sweep loop starts at most once
	// (on the first auth call, see startSweep).
	cleanupOnce sync.Once
}

// NewAuthService builds the auth service.
func NewAuthService(
	logger *slog.Logger,
	users *userservice.UserService,
	jwtService *jwt.Service,
	sessions *repo.AuthRepository,
	tracer *tracingpkg.Tracer,
) *AuthService {
	return &AuthService{
		logger:   logger,
		users:    users,
		jwt:      jwtService,
		sessions: sessions,
		tracer:   tracer,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*dto.SessionResult, error) {
	ctx, end := s.tracer.Start(ctx, "auth.Login")
	defer end(nil)
	s.startSweep()

	// Logins are stored lowercase; normalize the input so case-insensitive
	// login matches the stored value without leaking case sensitivity.
	username = strings.ToLower(strings.TrimSpace(username))

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

	access, err := s.jwt.GenerateAccessToken(user.ID, user.Username)
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
	s.startSweep()

	session, err := s.findSession(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if session.RevokedAt.Valid {
		// Reusing a revoked token indicates theft: revoke all of the user's
		// active sessions. The failure is logged loudly — a failed revocation
		// must never be silently swallowed, and no new pair is issued either
		// (a still-valid stolen token would otherwise keep rotating).
		if rerr := s.sessions.RevokeAllUserSessions(ctx, session.UserID); rerr != nil {
			s.logger.ErrorContext(ctx,
				"auth: не удалось отозвать сессии при повторном использовании токена",
				"user_id", session.UserID,
				"error", rerr,
			)
		}
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
	access, err := s.jwt.GenerateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token")
	}

	return &dto.SessionResult{
		Auth:         s.newAuthResponse(user, access),
		RefreshToken: refresh,
	}, nil
}

// Logout revokes the session by refresh token (idempotently).
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	ctx, end := s.tracer.Start(ctx, "auth.Logout")
	defer end(nil)
	s.startSweep()

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

func (s *AuthService) findSession(ctx context.Context, refreshToken string) (sqlc.FindSessionByHashRow, error) {
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

// startSweep launches the background expired-session sweep exactly once (the
// first auth call triggers it): dead sessions older than the sweep window are
// deleted hourly instead of only when a refresh happens.
func (s *AuthService) startSweep() {
	s.cleanupOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(sessionSweepInterval)
			defer ticker.Stop()
			for range ticker.C {
				s.sweepExpired(context.Background())
			}
		}()
	})
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
			Preset:     user.Preset,
		},
	}
}
