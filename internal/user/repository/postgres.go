//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/internal/user/repository/sqlc"
	nullable "github.com/Koshsky/erp-backend/pkg/database"
)

type UserRepository struct {
	logger *slog.Logger
	pool   *pgxpool.Pool
	db     *sqlc.Queries
}

// NewUserRepository builds the UserRepository repository.
func NewUserRepository(logger *slog.Logger, pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		logger: logger,
		pool:   pool,
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

func (r *UserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	return r.db.UsernameExists(ctx, username)
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
		LastName:        user.LastName,
		FirstName:       user.FirstName,
		MiddleName:      nullable.ToString(user.MiddleName),
		Username:        user.Username,
		Role:            user.Role,
		PasswordHash:    user.PasswordHash,
		ManagerID:       nullable.ToInt8(user.ManagerID),
		Position:        user.Position,
		HireDate:        toDate(user.HireDate),
		TerminationDate: toDate(user.TerminationDate),
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
		UserID:          user.ID,
		LastName:        user.LastName,
		FirstName:       user.FirstName,
		MiddleName:      nullable.ToString(user.MiddleName),
		Username:        user.Username,
		Role:            user.Role,
		PasswordHash:    user.PasswordHash,
		ManagerID:       nullable.ToInt8(user.ManagerID),
		Position:        user.Position,
		HireDate:        toDate(user.HireDate),
		TerminationDate: toDate(user.TerminationDate),
	})
	if err != nil {
		return nil, err
	}

	mapped := mapUser(row)
	return &mapped, nil
}

func (r *UserRepository) ListUsers(
	ctx context.Context,
	userID int64,
	role string,
	roleFilter string,
	managerID int64,
	limit, offset int,
) ([]domain.User, error) {
	rows, err := r.db.ListUsers(ctx, sqlc.ListUsersParams{
		RoleFilter: roleFilter,
		IsAdmin:    role == domain.Admin,
		UserID:     userID,
		ManagerID:  managerID,
		PageLimit:  int64(limit),
		PageOffset: int64(offset),
	})
	if err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, mapUser(row))
	}
	return users, nil
}

func (r *UserRepository) CountUsers(
	ctx context.Context,
	userID int64,
	role string,
	roleFilter string,
	managerID int64,
) (int64, error) {
	return r.db.CountUsers(ctx, sqlc.CountUsersParams{
		RoleFilter: roleFilter,
		IsAdmin:    role == domain.Admin,
		UserID:     userID,
		ManagerID:  managerID,
	})
}

func (r *UserRepository) ListAllUsers(ctx context.Context) ([]domain.User, error) {
	rows, err := r.db.ListAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	users := make([]domain.User, 0, len(rows))
	for _, row := range rows {
		users = append(users, mapUser(row))
	}
	return users, nil
}

// ================= worker days (user_states) =================

func (r *UserRepository) ListStates(
	ctx context.Context,
	userID int64,
	start, end time.Time,
) ([]domain.UserState, error) {
	rows, err := r.db.ListStatesByUserRange(ctx, sqlc.ListStatesByUserRangeParams{
		UserID:    userID,
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return nil, err
	}
	states := make([]domain.UserState, 0, len(rows))
	for _, row := range rows {
		states = append(states, mapStateRow(row))
	}
	return states, nil
}

// overlapState is an overlapping interval with its state.
type overlapState struct {
	StateID   int64
	StartDate time.Time
	EndDate   time.Time
}

// SetStateRange overwrites the [start, end] range with the given state:
// subtracts [start, end] from overlapping intervals (keeping their residues),
// deletes covered ones and inserts the new one.
func (r *UserRepository) SetStateRange(
	ctx context.Context,
	userID, stateID int64,
	start, end time.Time,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := sqlc.New(tx)
	rows, err := q.ListOverlappingStates(ctx, sqlc.ListOverlappingStatesParams{
		UserID:    userID,
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return err
	}
	if err = q.DeleteOverlapping(ctx, sqlc.DeleteOverlappingParams{
		UserID:    userID,
		StartDate: start,
		EndDate:   end,
	}); err != nil {
		return err
	}
	if err = insertResidues(ctx, q, userID, toOverlapStates(rows), start, end); err != nil {
		return err
	}
	if _, err = q.InsertStateRange(ctx, sqlc.InsertStateRangeParams{
		UserID:    userID,
		StateID:   stateID,
		StartDate: start,
		EndDate:   end,
	}); err != nil {
		return err
	}

	// Сливаем смежные/пересекающиеся диапазоны того же состояния в непрерывный:
	// fn_normalize_user_states() схлопывает соседние интервалы одинакового
	// (user_id, state_id) сразу в рамках этой транзакции.
	if err = q.NormalizeUserStates(ctx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// DeleteStateRange clears the [start, end] range (keeping residues of
// overlapping intervals); with stateID != nil, only states of that type.
func (r *UserRepository) DeleteStateRange(
	ctx context.Context,
	userID int64,
	start, end time.Time,
	stateID *int64,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := sqlc.New(tx)
	var overlaps []overlapState
	if stateID == nil {
		overlaps, err = loadAndDeleteAll(ctx, q, userID, start, end)
	} else {
		overlaps, err = loadAndDeleteByState(ctx, q, userID, *stateID, start, end)
	}
	if err != nil {
		return err
	}
	if err = insertResidues(ctx, q, userID, overlaps, start, end); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// loadAndDeleteAll loads and deletes all intervals overlapping [start, end].
func loadAndDeleteAll(
	ctx context.Context,
	q *sqlc.Queries,
	userID int64,
	start, end time.Time,
) ([]overlapState, error) {
	rows, err := q.ListOverlappingStates(ctx, sqlc.ListOverlappingStatesParams{
		UserID:    userID,
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return nil, err
	}
	return toOverlapStates(rows), q.DeleteOverlapping(ctx, sqlc.DeleteOverlappingParams{
		UserID:    userID,
		StartDate: start,
		EndDate:   end,
	})
}

// loadAndDeleteByState loads and deletes intervals of the given state overlapping [start, end].
func loadAndDeleteByState(
	ctx context.Context,
	q *sqlc.Queries,
	userID, stateID int64,
	start, end time.Time,
) ([]overlapState, error) {
	rows, err := q.ListOverlappingStatesByState(ctx, sqlc.ListOverlappingStatesByStateParams{
		UserID:    userID,
		StateID:   stateID,
		StartDate: start,
		EndDate:   end,
	})
	if err != nil {
		return nil, err
	}
	return toOverlapStatesByState(rows), q.DeleteOverlappingByState(ctx, sqlc.DeleteOverlappingByStateParams{
		UserID:    userID,
		StateID:   stateID,
		StartDate: start,
		EndDate:   end,
	})
}

// insertResidues inserts the parts of overlapping intervals outside [start, end]
// (left/right residues), keeping their original state.
func insertResidues(
	ctx context.Context,
	q *sqlc.Queries,
	userID int64,
	overlaps []overlapState,
	start, end time.Time,
) error {
	for _, o := range overlaps {
		if o.StartDate.Before(start) {
			if _, err := q.InsertStateRange(ctx, sqlc.InsertStateRangeParams{
				UserID:    userID,
				StateID:   o.StateID,
				StartDate: o.StartDate,
				EndDate:   start.AddDate(0, 0, -1),
			}); err != nil {
				return err
			}
		}
		if o.EndDate.After(end) {
			if _, err := q.InsertStateRange(ctx, sqlc.InsertStateRangeParams{
				UserID:    userID,
				StateID:   o.StateID,
				StartDate: end.AddDate(0, 0, 1),
				EndDate:   o.EndDate,
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func toOverlapStates(rows []sqlc.ListOverlappingStatesRow) []overlapState {
	result := make([]overlapState, 0, len(rows))
	for _, row := range rows {
		result = append(result, overlapState{
			StateID:   row.StateID,
			StartDate: row.StartDate,
			EndDate:   row.EndDate,
		})
	}
	return result
}

func toOverlapStatesByState(rows []sqlc.ListOverlappingStatesByStateRow) []overlapState {
	result := make([]overlapState, 0, len(rows))
	for _, row := range rows {
		result = append(result, overlapState{
			StateID:   row.StateID,
			StartDate: row.StartDate,
			EndDate:   row.EndDate,
		})
	}
	return result
}

func mapUser(row sqlc.User) domain.User {
	return domain.User{
		ID:              row.ID,
		LastName:        row.LastName,
		FirstName:       row.FirstName,
		MiddleName:      nullable.StringPtr(row.MiddleName),
		Role:            row.Role,
		Username:        row.Username,
		PasswordHash:    row.PasswordHash,
		ManagerID:       nullable.Int64Ptr(row.ManagerID),
		Position:        row.Position,
		HireDate:        fromDate(row.HireDate),
		TerminationDate: fromDate(row.TerminationDate),
	}
}

func mapStateRow(row sqlc.ListStatesByUserRangeRow) domain.UserState {
	return domain.UserState{
		ID:          row.ID,
		UserID:      row.UserID,
		StateID:     row.StateID,
		StateCode:   row.StateCode,
		StateName:   row.StateName,
		IsAvailable: row.IsAvailable,
		StartDate:   row.StartDate,
		EndDate:     row.EndDate,
	}
}

// fromDate unwraps a nullable date (pgtype.Date) into [time.Time].
func fromDate(v pgtype.Date) *time.Time {
	if !v.Valid {
		return nil
	}
	t := v.Time
	return &t
}

// toDate wraps [time.Time] into a nullable date (pgtype.Date).
func toDate(v *time.Time) pgtype.Date {
	if v == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: *v, Valid: true}
}

// OwnerChain returns the owner chain (manager_id → the vp) for RBAC checks.
func (r *UserRepository) OwnerChain(ctx context.Context, id int64) (rbac.Owners, error) {
	owner, err := r.db.OwnerChain(ctx, id)
	if err != nil {
		return rbac.Owners{}, err
	}
	return rbac.Owners{Owner: owner}, nil
}
