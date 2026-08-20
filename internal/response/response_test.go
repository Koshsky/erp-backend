package response_test

import (
	"encoding/json"
	stderrors "errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/response"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

type body struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type envelope struct {
	Error *body `json:"error"`
}

// run serves the given handler through a real gin router and returns the body.
func run(t *testing.T, handler gin.HandlerFunc) envelope {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Any("/test", handler)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(rec, req)

	var env envelope
	if err := json.Unmarshal(rec.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal response: %v (body=%s)", err, rec.Body.String())
	}
	return env
}

func TestInternalErrorSendsGeneric(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(&strings.Builder{}, nil))
	env := run(t, func(c *gin.Context) {
		response.InternalError(c, logger, "db: failed", stderrors.New("connection refused"))
	})

	if env.Error == nil {
		t.Fatal("expected error envelope")
	}
	if env.Error.Code != "INTERNAL_ERROR" {
		t.Errorf("code = %q, want INTERNAL_ERROR", env.Error.Code)
	}
	if env.Error.Message != "internal server error" {
		t.Errorf("message = %q, want %q", env.Error.Message, "internal server error")
	}
	if strings.Contains(env.Error.Message, "db") || strings.Contains(env.Error.Message, "refused") {
		t.Errorf("internal details leaked: %q", env.Error.Message)
	}
}

func TestErrorSendsGenericOn500(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(&strings.Builder{}, nil))
	env := run(t, func(c *gin.Context) {
		response.Error(c, logger, stderrors.New("dial tcp 10.0.0.1:5432: connect: connection refused"))
	})

	if env.Error.Code != "INTERNAL_ERROR" {
		t.Errorf("code = %q, want INTERNAL_ERROR", env.Error.Code)
	}
	if env.Error.Message != "internal server error" {
		t.Errorf("message = %q, want %q", env.Error.Message, "internal server error")
	}
	if strings.Contains(env.Error.Message, "dial") || strings.Contains(env.Error.Message, "5432") {
		t.Errorf("internal details leaked: %q", env.Error.Message)
	}
}

func TestErrorKeepsDomainMessageOn4xx(t *testing.T) {
	t.Parallel()
	logger := slog.New(slog.NewTextHandler(&strings.Builder{}, nil))
	env := run(t, func(c *gin.Context) {
		response.Error(c, logger, errors.NotFound("user missing"))
	})

	if env.Error.Code != "NOT_FOUND" {
		t.Errorf("code = %q, want NOT_FOUND", env.Error.Code)
	}
	if env.Error.Message != "user missing" {
		t.Errorf("message = %q, want %q", env.Error.Message, "user missing")
	}
}
