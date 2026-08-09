//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/internal/user/repository/sqlc"
)

type UserRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

// NewUserRepository builds the UserRepository repository.
func NewUserRepository(logger *slog.Logger, pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *UserRepository) FindUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	row, err := r.db.FindUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	mapped := mapUser(row)
	return &mapped, nil
}

func (r *UserRepository) FindUserByID(ctx context.Context, userID int64) (*domain.User, error) {
	return r.FindUser(ctx, userID)
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID int64, hash string) error {
	return r.db.UpdateUserPassword(ctx, sqlc.UpdateUserPasswordParams{
		UserID:       userID,
		PasswordHash: hash,
	})
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int64) error {
	return r.db.DeleteUser(ctx, id)
}

func (r *UserRepository) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	row, err := r.db.CreateUser(ctx, sqlc.CreateUserParams{
		Name:         user.Name,
		Username:     user.Username,
		Role:         user.Role,
		PasswordHash: user.PasswordHash,
	})
	if err != nil {
		return nil, err
	}

	mapped := mapUser(row)
	return &mapped, nil
}

func (r *UserRepository) FindUser(ctx context.Context, id int64) (*domain.User, error) {
	row, err := r.db.FindUser(ctx, id)
	if err != nil {
		return nil, err
	}

	mapped := mapUser(row)
	return &mapped, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	row, err := r.db.UpdateUser(ctx, sqlc.UpdateUserParams{
		UserID:       user.ID,
		Name:         user.Name,
		Username:     user.Username,
		Role:         user.Role,
		PasswordHash: user.PasswordHash,
	})
	if err != nil {
		return nil, err
	}

	mapped := mapUser(row)
	return &mapped, nil
}

func (r *UserRepository) ListUsers(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, mapUser(row))
	}

	return users, nil
}

func mapUser(row sqlc.User) domain.User {
	return domain.User{
		ID:           row.ID,
		Name:         row.Name,
		Role:         row.Role,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
	}
}
