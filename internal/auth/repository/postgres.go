//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/auth/repository/sqlc"
)

// AuthRepository persists refresh sessions (rotation/revocation, AD-06).
type AuthRepository struct {
	db *sqlc.Queries
}

// NewAuthRepository builds the auth repository.
func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: sqlc.New(pool)}
}

func (r *AuthRepository) FindSessionByHash(ctx context.Context, tokenHash string) (sqlc.FindSessionByHashRow, error) {
	return r.db.FindSessionByHash(ctx, tokenHash)
}

func (r *AuthRepository) CreateSession(
	ctx context.Context,
	userID int64,
	tokenHash string,
	expiresAt time.Time,
) (sqlc.CreateSessionRow, error) {
	return r.db.CreateSession(ctx, sqlc.CreateSessionParams{
		UserID:     userID,
		TokenHash:  tokenHash,
		ExpiresAt:  expiresAt,
		ReplacedBy: pgtype.Int8{}, // the replaced_by chain is not filled (reuse detection is based on revoked_at)
	})
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
