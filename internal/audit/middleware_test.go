//nolint:testpackage // tests the unexported bodyWriter and buildEvent directly
package audit

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/config"
	userctx "github.com/Koshsky/erp-backend/internal/userctx"
)

func TestBodyWriterCapturesStatusAndBody(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	bw := &bodyWriter{ResponseWriter: c.Writer}
	if bw.Status() != http.StatusOK {
		t.Fatalf("default status must be 200")
	}
	bw.WriteHeader(http.StatusCreated)
	if bw.Status() != http.StatusCreated {
		t.Fatalf("status must record 201")
	}
	_, _ = bw.Write([]byte(`{"ok":true}`))
	if !bytes.Equal(bw.Body(), []byte(`{"ok":true}`)) {
		t.Fatalf("body must be captured, got %s", bw.Body())
	}
	if rec.Code != http.StatusCreated {
		t.Fatalf("real writer must receive the status, got %d", rec.Code)
	}
}

func TestBuildEventExtractsActorAndID(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	cfg := config.AuditConfig{Enabled: true, URL: "http://loki:3100"}
	mw := NewMiddleware(slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)), cfg, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"code":"P-1","password":"secret"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/project", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set(userctx.KeyUser, userctx.UserContext{ID: 3, Email: "admin@x.ru", Role: "admin"})

	bw := &bodyWriter{ResponseWriter: c.Writer}
	c.Writer = bw
	bw.WriteHeader(http.StatusCreated)
	_, _ = bw.Write([]byte(`{"data":{"id":12},"error":{}}`))

	ev := mw.buildEvent(c, routeClass{entity: entProject, action: actCreate}, time.Now(), []byte(body), bw)
	if ev == nil {
		t.Fatalf("event must be built")
	}
	if ev.Entity != "project" || ev.Action != "create" {
		t.Fatalf("unexpected entity/action: %s/%s", ev.Entity, ev.Action)
	}
	if ev.ActorUserID == nil || *ev.ActorUserID != 3 {
		t.Fatalf("actor id must be extracted, got %v", ev.ActorUserID)
	}
	if ev.ActorEmail != "admin@x.ru" || ev.ActorRole != "admin" {
		t.Fatalf("actor email/role must be extracted, got %q/%q", ev.ActorEmail, ev.ActorRole)
	}
	if ev.EntityID == nil || *ev.EntityID != 12 {
		t.Fatalf("entity id must come from the response, got %v", ev.EntityID)
	}
	if ev.Status != http.StatusCreated {
		t.Fatalf("status must be 201, got %d", ev.Status)
	}
	if bytes.Contains(ev.RequestBody, []byte("secret")) {
		t.Fatalf("request body must be masked, got %s", ev.RequestBody)
	}
}

func TestBuildEventPublicAuthUsesUsername(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	cfg := config.AuditConfig{Enabled: true, URL: "http://loki:3100"}
	mw := NewMiddleware(slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)), cfg, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	body := `{"username":"ivanov","password":"secret"}`
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(body))
	c.Request.Header.Set("Content-Type", "application/json")

	bw := &bodyWriter{ResponseWriter: c.Writer}
	c.Writer = bw
	bw.WriteHeader(http.StatusOK)
	_, _ = bw.Write([]byte(`{"data":{"user":{"id":9}},"error":{}}`))

	ev := mw.buildEvent(c, routeClass{entity: entityAuth, action: actionLogin}, time.Now(), []byte(body), bw)
	if ev == nil {
		t.Fatalf("event must be built")
	}
	if ev.ActorEmail != "ivanov" {
		t.Fatalf("login actor must be the username, got %q", ev.ActorEmail)
	}
	if ev.ActorUserID == nil || *ev.ActorUserID != 9 {
		t.Fatalf("login actor id must come from the response user, got %v", ev.ActorUserID)
	}
	if ev.EntityID != nil {
		t.Fatalf("auth events must not carry an entity id, got %v", *ev.EntityID)
	}
	if bytes.Contains(ev.RequestBody, []byte("secret")) {
		t.Fatalf("login password must be masked, got %s", ev.RequestBody)
	}
}

func TestBuildEventAuthRefreshGetsActorFromResponse(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	cfg := config.AuditConfig{Enabled: true, URL: "http://loki:3100"}
	mw := NewMiddleware(slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)), cfg, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// /auth/refresh has an empty request body (token lives in the cookie).
	c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", bytes.NewBuffer(nil))

	bw := &bodyWriter{ResponseWriter: c.Writer}
	c.Writer = bw
	bw.WriteHeader(http.StatusOK)
	_, _ = bw.Write([]byte(`{"data":{"user":{"id":4,"username":"admin"}},"error":{}}`))

	ev := mw.buildEvent(c, routeClass{entity: entityAuth, action: actionRefresh}, time.Now(), nil, bw)
	if ev == nil {
		t.Fatalf("event must be built")
	}
	if ev.ActorUserID == nil || *ev.ActorUserID != 4 {
		t.Fatalf("refresh actor id must come from the response user, got %v", ev.ActorUserID)
	}
	if ev.EntityID != nil {
		t.Fatalf("auth events must not carry an entity id, got %v", *ev.EntityID)
	}
	// Empty body → no username yet (the read-side enrichment fills the login).
	if ev.ActorEmail != "" {
		t.Fatalf("refresh has no username in the request, got %q", ev.ActorEmail)
	}
}
