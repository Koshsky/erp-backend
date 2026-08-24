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
func (r *IdempotencyRepository) Claim(
	ctx context.Context,
	key string,
	userID int64,
	method, path string,
	expiresAt time.Time,
) (*StoredResult, bool, error) {
	_, err := r.db.CreateIdempotencyKey(ctx, sqlc.CreateIdempotencyKeyParams{
		Key:       key,
		UserID:    userID,
		Method:    method,
		Path:      path,
		ExpiresAt: expiresAt,
	})
	if err == nil {
		// Мы захватили ключ.
		return nil, true, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, false, err
	}
	// Ключ уже существует — смотрим, завершился ли первый вызов.
	return r.loadExisting(ctx, key, userID, method, path)
}

// loadExisting читает результат по уже существующему ключу: возвращает
// сохранённый ответ, если он завершился, иначе — признак «в полёте».
func (r *IdempotencyRepository) loadExisting(
	ctx context.Context,
	key string,
	userID int64,
	method, path string,
) (*StoredResult, bool, error) {
	existing, err := r.db.GetIdempotencyKey(ctx, sqlc.GetIdempotencyKeyParams{
		Key:    key,
		UserID: userID,
		Method: method,
		Path:   path,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Ключ исчез между checks — считаем конфликтом без воспроизведения.
			return nil, false, nil
		}
		return nil, false, err
	}
	if existing.ResponseStatus <= 0 {
		// Первый запрос ещё выполняется.
		return nil, false, nil
	}
	return &StoredResult{Status: int(existing.ResponseStatus), Body: existing.ResponseBody}, false, nil
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
