//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/auth/repository/sqlc"
)

// Session is an active refresh session row.
type Session struct {
	ID         int64
	UserID     int64
	TokenHash  string
	CreatedAt  time.Time
	ExpiresAt  time.Time
	RevokedAt  *time.Time
	ReplacedBy int64
}

// AuthRepository persists refresh sessions (rotation/revocation, AD-06).
type AuthRepository struct {
	db *sqlc.Queries
}

// NewAuthRepository builds the auth repository.
func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: sqlc.New(pool)}
}

func (r *AuthRepository) FindSessionByHash(ctx context.Context, tokenHash string) (Session, error) {
	row, err := r.db.FindSessionByHash(ctx, tokenHash)
	if err != nil {
		return Session{}, err
	}
	return toSession(
		row.ID,
		row.UserID,
		row.TokenHash,
		row.CreatedAt,
		row.ExpiresAt,
		row.RevokedAt,
		row.ReplacedBy,
	), nil
}

func (r *AuthRepository) CreateSession(
	ctx context.Context,
	userID int64,
	tokenHash string,
	expiresAt time.Time,
) (Session, error) {
	row, err := r.db.CreateSession(ctx, sqlc.CreateSessionParams{
		UserID:     userID,
		TokenHash:  tokenHash,
		ExpiresAt:  expiresAt,
		ReplacedBy: pgtype.Int8{}, // цепочка replaced_by не заполняется (детект reuse идёт по revoked_at)
	})
	if err != nil {
		return Session{}, err
	}
	return toSession(
		row.ID,
		row.UserID,
		row.TokenHash,
		row.CreatedAt,
		row.ExpiresAt,
		row.RevokedAt,
		row.ReplacedBy,
	), nil
}

func (r *AuthRepository) RevokeSession(ctx context.Context, id int64) error {
	return r.db.RevokeSession(ctx, id)
}

func (r *AuthRepository) RevokeAllUserSessions(ctx context.Context, userID int64) error {
	return r.db.RevokeAllUserSessions(ctx, userID)
}

func (r *AuthRepository) DeleteExpiredSessions(ctx context.Context, olderThan time.Time) error {
	return r.db.DeleteExpiredSessions(ctx, olderThan)
}

func toSession(
	id int64,
	userID int64,
	tokenHash string,
	createdAt time.Time,
	expiresAt time.Time,
	revokedAt **time.Time,
	replacedBy int64,
) Session {
	var revoked *time.Time
	if revokedAt != nil && *revokedAt != nil {
		revoked = *revokedAt
	}
	return Session{
		ID:         id,
		UserID:     userID,
		TokenHash:  tokenHash,
		CreatedAt:  createdAt,
		ExpiresAt:  expiresAt,
		RevokedAt:  revoked,
		ReplacedBy: replacedBy,
	}
}
