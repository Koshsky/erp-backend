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

func NewUserRepository(logger *slog.Logger, pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		logger: logger,
		db:     sqlc.New(pool),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	row, err := r.db.CreateUser(ctx, sqlc.CreateUserParams{
		Name:         user.Name,
		Username:     user.Username,
		Role:         string(user.Role),
		PasswordHash: user.PasswordHash,
	})
	if err != nil {
		return nil, err
	}

	mapped := mapUser(row)
	return &mapped, nil
}

func (r *UserRepository) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	row, err := r.db.GetUser(ctx, id)
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
		Role:         string(user.Role),
		PasswordHash: user.PasswordHash,
	})
	if err != nil {
		return nil, err
	}

	mapped := mapUser(row)
	return &mapped, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int64) error {
	return r.db.DeleteUser(ctx, id)
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
		Role:         domain.UserRole(row.Role),
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
	}
}
