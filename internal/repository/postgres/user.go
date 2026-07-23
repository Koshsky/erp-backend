package postgres

import (
	"context"
	"log/slog"

	"github.com/Koshsky/erp/api/internal/domain"
	"github.com/Koshsky/erp/api/internal/repository/postgres/sqlc"
)

type UserRepository struct {
	logger *slog.Logger
	db     *sqlc.Queries
}

func NewUserRepository(logger *slog.Logger, db *sqlc.Queries) *UserRepository {
	return &UserRepository{logger: logger, db: db}
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
