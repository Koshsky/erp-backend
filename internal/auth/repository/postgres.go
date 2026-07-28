package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{pool: pool}
}

func (r *AuthRepository) FindUserByUsername(ctx context.Context, username string) (id int64, name, role, passwordHash string, err error) {
	err = r.pool.QueryRow(ctx,
		`SELECT id, name, role, password_hash FROM users WHERE username = $1 AND deleted_at IS NULL`,
		username,
	).Scan(&id, &name, &role, &passwordHash)
	return
}

func (r *AuthRepository) FindUserByID(ctx context.Context, userID int64) (id int64, name, role, passwordHash string, err error) {
	err = r.pool.QueryRow(ctx,
		`SELECT id, name, role, password_hash FROM users WHERE id = $1 AND deleted_at IS NULL`,
		userID,
	).Scan(&id, &name, &role, &passwordHash)
	return
}

func (r *AuthRepository) CreateUser(ctx context.Context, name, username, role, passwordHash string) (int64, error) {
	var id int64
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (name, username, role, password_hash) VALUES ($1, $2, $3, $4) RETURNING id`,
		name, username, role, passwordHash,
	).Scan(&id)
	return id, err
}

func (r *AuthRepository) UpdatePassword(ctx context.Context, userID int64, newHash string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2 AND deleted_at IS NULL`,
		newHash, userID,
	)
	return err
}

func (r *AuthRepository) SaveRefreshToken(ctx context.Context, userID int64, token string) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token) VALUES ($1, $2)
		 ON CONFLICT (user_id) DO UPDATE SET token = $2, created_at = NOW()`,
		userID, token,
	)
	return err
}

func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, userID int64) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM refresh_tokens WHERE user_id = $1`, userID,
	)
	return err
}
