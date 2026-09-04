//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0 generate -f ./sqlc.yaml
package repository

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Koshsky/erp-backend/internal/idempotency/repository/sqlc"
)

// StoredResult is a saved response of a completed idempotency key.
type StoredResult struct {
	Status int
	Body   json.RawMessage
}

// claimLease — how long an in-flight key is considered "alive"; longer means a crash.
const claimLease = 2 * time.Minute

// claimRetryAttempts — how many times we repeat observing/re-claiming the key.
const claimRetryAttempts = 2

// IdempotencyRepository stores/retrieves idempotency keys and their responses.
type IdempotencyRepository struct {
	db *sqlc.Queries
}

// NewIdempotencyRepository builds the IdempotencyRepository repository.
func NewIdempotencyRepository(pool *pgxpool.Pool) *IdempotencyRepository {
	return &IdempotencyRepository{db: sqlc.New(pool)}
}

// Claim atomically claims the key for this (user, method, path).
//
// Returns:
//   - claimed=true when this call acquired the key (the caller must execute the
//     operation and then call Complete/Release).
//   - claimed=false with a non-nil result when the key was already completed by a
//     previous call (the caller must replay result without re-executing).
//   - claimed=false with a nil result when the key exists but the request is
//     still in flight (the caller should return a conflict, not re-execute).
//
// Stale keys are re-claimed: a completed key past its TTL or an in-flight key
// older than the lease (a crashed request) is deleted and claimed again, so a
// retry is not blocked by a dead claim.
func (r *IdempotencyRepository) Claim(
	ctx context.Context,
	key string,
	userID int64,
	method, path string,
	expiresAt time.Time,
) (*StoredResult, bool, error) {
	params := sqlc.CreateIdempotencyKeyParams{
		Key:       key,
		UserID:    userID,
		Method:    method,
		Path:      path,
		ExpiresAt: expiresAt,
	}
	_, err := r.db.CreateIdempotencyKey(ctx, params)
	if err == nil {
		// We claimed the key.
		return nil, true, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, false, err
	}
	// The key already exists — check whether the first call completed; a stale
	// (expired or stuck) key is re-claimed. A stale key lost in a re-claim race
	// is observed again on the next iteration.
	for attempt := range claimRetryAttempts {
		result, claimed, oerr := r.claimExisting(ctx, params)
		if oerr != nil {
			return nil, false, oerr
		}
		if claimed {
			return nil, true, nil
		}
		if result != nil {
			return result, false, nil
		}
		// In flight or lost a race: return 409 on the last iteration.
		if attempt == claimRetryAttempts-1 {
			return nil, false, nil
		}
	}
	return nil, false, nil
}

// claimExisting observes the current state of an existing key. A stale key
// (expired or in-flight longer than the lease — a crash) is deleted and re-claimed.
// Returns: a ready replay result (result, claimed=false),
// a claim indicator (claimed=true), or "in flight" (nil, false, nil).
func (r *IdempotencyRepository) claimExisting(
	ctx context.Context,
	params sqlc.CreateIdempotencyKeyParams,
) (*StoredResult, bool, error) {
	existing, err := r.db.GetIdempotencyKey(ctx, sqlc.GetIdempotencyKeyParams{
		Key:    params.Key,
		UserID: params.UserID,
		Method: params.Method,
		Path:   params.Path,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// The key disappeared between checks — treat it as a conflict without replay.
			return nil, false, nil
		}
		return nil, false, err
	}
	if reclaimable(existing, time.Now(), claimLease) {
		if derr := r.db.DeleteIdempotencyKey(ctx, sqlc.DeleteIdempotencyKeyParams{
			Key:    params.Key,
			UserID: params.UserID,
			Method: params.Method,
			Path:   params.Path,
		}); derr != nil {
			return nil, false, derr
		}
		// Re-claim (one of the retries wins the race).
		_, ierr := r.db.CreateIdempotencyKey(ctx, params)
		if ierr == nil {
			return nil, true, nil
		}
		if !errors.Is(ierr, pgx.ErrNoRows) {
			return nil, false, ierr
		}
		return nil, false, nil
	}
	if existing.ResponseStatus <= 0 {
		// The first request is still running.
		return nil, false, nil
	}
	return &StoredResult{Status: int(existing.ResponseStatus), Body: existing.ResponseBody}, false, nil
}

// reclaimable reports whether an existing key row should be deleted and
// re-claimed: completed but expired (cleanup has not run yet), or in-flight
// for longer than the lease (the claiming request crashed).
func reclaimable(existing sqlc.IdempotencyKey, now time.Time, lease time.Duration) bool {
	if existing.ExpiresAt.Before(now) {
		return true
	}
	return existing.ResponseStatus <= 0 && existing.CreatedAt.Before(now.Add(-lease))
}

// Complete saves the response of the completed call under the key.
func (r *IdempotencyRepository) Complete(
	ctx context.Context,
	key string,
	userID int64,
	method, path string,
	status int,
	body json.RawMessage,
) error {
	// HTTP statuses are always small and fit into int32; a safeguard against
	// range overflow (gosec G115) — we never store an incorrect status in the DB.
	if status <= 0 || status > math.MaxInt32 {
		status = http.StatusInternalServerError
	}
	return r.db.CompleteIdempotencyKey(ctx, sqlc.CompleteIdempotencyKeyParams{
		Key:            key,
		UserID:         userID,
		Method:         method,
		Path:           path,
		ResponseStatus: int32(status),
		ResponseBody:   body,
	})
}

// Release deletes the key (e.g. on 5xx so a retry can run the operation again).
func (r *IdempotencyRepository) Release(
	ctx context.Context,
	key string,
	userID int64,
	method, path string,
) error {
	return r.db.DeleteIdempotencyKey(ctx, sqlc.DeleteIdempotencyKeyParams{
		Key:    key,
		UserID: userID,
		Method: method,
		Path:   path,
	})
}

// DeleteExpired deletes keys whose TTL has expired.
func (r *IdempotencyRepository) DeleteExpired(ctx context.Context) error {
	return r.db.DeleteExpiredIdempotencyKeys(ctx)
}
