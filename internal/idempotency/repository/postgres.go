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

// claimLease — сколько in-flight ключ считается «живым»; дольше — краш.
const claimLease = 2 * time.Minute

// claimRetryAttempts — сколько раз повторяем наблюдение/перезахват ключа.
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
		// Мы захватили ключ.
		return nil, true, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, false, err
	}
	// Ключ уже существует — смотрим, завершился ли первый вызов; устаревший
	// (истёкший или зависший) ключ перезахватываем. Устаревший ключ, потерянный
	// в гонке перезахвата, на следующей итерации наблюдается заново.
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
		// «В полёте» или потерянная гонка: на последней итерации возвращаем 409.
		if attempt == claimRetryAttempts-1 {
			return nil, false, nil
		}
	}
	return nil, false, nil
}

// claimExisting наблюдает текущее состояние существующего ключа. Устаревший
// ключ (истёкший или in-flight дольше lease — краш) удаляется и перезахватывается.
// Возвращает: готовый результат реплея (result, claimed=false),
// признак захвата (claimed=true), либо «в полёте» (nil, false, nil).
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
			// Ключ исчез между checks — считаем конфликтом без воспроизведения.
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
		// Повторный захват (гонку выигрывает один из ретраев).
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
		// Первый запрос ещё выполняется.
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

// Complete сохраняет ответ завершённого вызова по ключу.
func (r *IdempotencyRepository) Complete(
	ctx context.Context,
	key string,
	userID int64,
	method, path string,
	status int,
	body json.RawMessage,
) error {
	// HTTP-статусы всегда малы и помещаются в int32; страховка от выхода за
	// диапазон (gosec G115) — не кладём некорректный статус в БД.
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

// Release удаляет ключ (например при 5xx, чтобы ретрай мог повторить операцию).
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

// DeleteExpired удаляет ключи, у которых истёк TTL.
func (r *IdempotencyRepository) DeleteExpired(ctx context.Context) error {
	return r.db.DeleteExpiredIdempotencyKeys(ctx)
}
