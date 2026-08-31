// Package idempotency provides an Idempotency-Key mechanism that makes create
// endpoints replay-safe: a client sends an Idempotency-Key header; the server
// atomically claims the key and on a repeat with the same key returns the saved
// response instead of executing the operation again (no duplicate records).
package idempotency

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/idempotency/repository"
	"github.com/Koshsky/erp-backend/internal/response"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"
	userctx "github.com/Koshsky/erp-backend/internal/userctx"
	errapi "github.com/Koshsky/erp-backend/pkg/errors"
)

// HTTP header the client uses to mark an idempotent request.
const headerIdempotencyKey = "Idempotency-Key"

// maxKeyLen — верхний предел длины ключа (защита PK от мусорных/гигантских значений).
const maxKeyLen = 256

// keyTTL — сколько жить ключу идемпотентности до авто-очистки.
const keyTTL = 24 * time.Hour

// cleanupInterval — как часто вычищаем просроченные ключи.
const cleanupInterval = 1 * time.Hour

// Repo is the storage the middleware uses to claim, complete and release
// idempotency keys. Kept as an interface so it can be faked in unit tests.
type Repo interface {
	Claim(
		ctx context.Context,
		key string,
		userID int64,
		method, path string,
		expiresAt time.Time,
	) (*repository.StoredResult, bool, error)
	Complete(
		ctx context.Context,
		key string,
		userID int64,
		method, path string,
		status int,
		body json.RawMessage,
	) error
	Release(ctx context.Context, key string, userID int64, method, path string) error
	DeleteExpired(ctx context.Context) error
}

// Middleware implements the Idempotency-Key mechanism for a route.
type Middleware struct {
	repo   Repo
	logger *slog.Logger
	tracer *tracingpkg.Tracer

	cleanupOnce sync.Once
}

// New builds the idempotency middleware.
func New(
	repo Repo,
	logger *slog.Logger,
	tracer *tracingpkg.Tracer,
) *Middleware {
	if logger == nil {
		logger = slog.Default()
	}
	if tracer == nil {
		tracer = tracingpkg.New(nil)
	}
	return &Middleware{repo: repo, logger: logger, tracer: tracer}
}

// Handler wraps the route handler with the idempotency logic.
func (m *Middleware) Handler() gin.HandlerFunc {
	m.startCleanup()
	return func(c *gin.Context) {
		key := c.GetHeader(headerIdempotencyKey)
		if key == "" || len(key) > maxKeyLen {
			// Без ключа (или с неоправданно длинным) — обычный неидемпотентный запрос.
			c.Next()
			return
		}

		userID, err := userctx.GetUserID(c)
		if err != nil {
			// Авторизованного пользователя нет — привязать ключ не к кому;
			// выполняем без идемпотентной гарантии.
			c.Next()
			return
		}

		ctx, end := m.tracer.Start(c.Request.Context(), "middleware.idempotency")
		c.Request = c.Request.WithContext(ctx)
		defer end(nil)

		method := c.Request.Method
		path := c.Request.URL.Path

		result, claimed, err := m.repo.Claim(ctx, key, userID, method, path, time.Now().Add(keyTTL))
		if err != nil {
			m.logger.Error("idempotency claim failed", "error", err, "key", key)
			response.InternalError(c, m.logger, "idempotency storage failure", err)
			c.Abort()
			return
		}

		if !claimed {
			if result == nil {
				// Ключ существует, но первый запрос ещё выполняется — не дублируем.
				response.Error(c, m.logger, errapi.Conflict("request already in progress"))
				c.Abort()
				return
			}
			// Повтор с тем же ключом: возвращаем сохранённый ответ.
			m.logger.Debug("idempotency replay", "key", key)
			replay(c, result.Status, result.Body)
			// Возврат без c.Next() в gin не останавливает цепочку хендлеров,
			// поэтому прерываем её явно, чтобы не выполнить операцию повторно.
			c.Abort()
			return
		}

		// Мы захватили ключ: выполняем операцию, перехватывая ответ.
		cw := &captureWriter{ResponseWriter: c.Writer}
		c.Writer = cw
		c.Next()

		m.finalize(ctx, cw, key, userID, method, path)
	}
}

// finalize сохраняет 2xx-ответ для replay либо освобождает ключ, чтобы ретрай
// (4xx/5xx) выполнил операцию заново и получил свежую оценку.
func (m *Middleware) finalize(
	ctx context.Context,
	cw *captureWriter,
	key string,
	userID int64,
	method, path string,
) {
	status := cw.status
	if status == 0 {
		status = http.StatusOK
	}
	if status >= http.StatusOK && status < http.StatusMultipleChoices {
		if cerr := m.repo.Complete(ctx, key, userID, method, path, status, json.RawMessage(cw.body)); cerr != nil {
			m.logger.ErrorContext(ctx, "idempotency complete failed", "error", cerr, "key", key)
		}
		return
	}
	if rerr := m.repo.Release(ctx, key, userID, method, path); rerr != nil {
		m.logger.ErrorContext(ctx, "idempotency release failed", "error", rerr, "key", key)
	}
}

// startCleanup запускает фоновую очистку просроченных ключей на время жизни
// процесса (аналогично фоновой чистке rate-limit бабокетов).
func (m *Middleware) startCleanup() {
	m.cleanupOnce.Do(func() {
		go func() {
			for {
				time.Sleep(cleanupInterval)
				if err := m.repo.DeleteExpired(context.Background()); err != nil {
					m.logger.Error("idempotency cleanup failed", "error", err)
				}
			}
		}()
	})
}

// replay отправляет сохранённый ответ клиенту. Все ответы приложения — JSON
// ({data,error}), поэтому Content-Type фиксируем; Idempotency-Replayed —
// диагностический маркер для клиента.
func replay(c *gin.Context, status int, body json.RawMessage) {
	c.Header("Content-Type", "application/json")
	c.Header("Idempotency-Replayed", "true")
	if len(body) == 0 {
		c.Status(status)
		return
	}
	c.Data(status, "application/json", body)
}

// captureWriter перехватывает статус и тело ответа текущего запроса.
type captureWriter struct {
	gin.ResponseWriter

	status int
	body   []byte
}

func (w *captureWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *captureWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	w.body = append(w.body, b...)
	return w.ResponseWriter.Write(b)
}

func (w *captureWriter) WriteString(s string) (int, error) {
	return w.Write([]byte(s))
}
