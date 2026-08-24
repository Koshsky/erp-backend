package idempotency_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	idem "github.com/Koshsky/erp-backend/internal/idempotency"
	"github.com/Koshsky/erp-backend/internal/idempotency/repository"
	userctx "github.com/Koshsky/erp-backend/internal/userctx"
)

// headerKey is the wire name of the Idempotency-Key header.
const headerKey = "Idempotency-Key"

// keyScope is the composite identifier of an idempotency claim.
type keyScope struct {
	key    string
	userID int64
	method string
	path   string
}

// fakeRepo is an in-memory Repo used to exercise the middleware without a DB.
type fakeRepo struct {
	mu        sync.Mutex
	complete  map[keyScope]repository.StoredResult
	inflight  map[keyScope]bool
	claims    int
	completes int
	releases  int
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		complete: map[keyScope]repository.StoredResult{},
		inflight: map[keyScope]bool{},
	}
}

func scope(key string, userID int64, method, path string) keyScope {
	return keyScope{key: key, userID: userID, method: method, path: path}
}

func (f *fakeRepo) Claim(
	_ context.Context,
	key string,
	userID int64,
	method, path string,
	_ time.Time,
) (*repository.StoredResult, bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	s := scope(key, userID, method, path)
	if res, ok := f.complete[s]; ok {
		r := res
		return &r, false, nil
	}
	if f.inflight[s] {
		return nil, false, nil
	}
	f.inflight[s] = true
	f.claims++
	return nil, true, nil
}

func (f *fakeRepo) Complete(
	_ context.Context,
	key string,
	userID int64,
	method, path string,
	status int,
	body json.RawMessage,
) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	s := scope(key, userID, method, path)
	f.complete[s] = repository.StoredResult{Status: status, Body: body}
	delete(f.inflight, s)
	f.completes++
	return nil
}

func (f *fakeRepo) Release(_ context.Context, key string, userID int64, method, path string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	s := scope(key, userID, method, path)
	delete(f.complete, s)
	delete(f.inflight, s)
	f.releases++
	return nil
}

func (f *fakeRepo) DeleteExpired(context.Context) error { return nil }

func (f *fakeRepo) completesCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.completes
}

// setUser installs the authenticated user (id) into the gin context.
func setUser(id int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(userctx.KeyUser, userctx.UserContext{ID: id})
		c.Next()
	}
}

// buildRouter wires an idempotency handler around a POST endpoint that counts
// executions and responds with a JSON body, returning the built router.
func buildRouter(mw *idem.Middleware, counter *int) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(setUser(7))
	router.Use(mw.Handler())
	router.POST("/tasks", func(c *gin.Context) {
		*counter++
		c.JSON(http.StatusCreated, gin.H{"id": *counter})
	})
	return router
}

// doPost issues a request through the router with an optional Idempotency-Key.
func doPost(router http.Handler, key string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	if key != "" {
		req.Header.Set(headerKey, key)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestIdempotencyNoKeyPassesThrough(t *testing.T) {
	t.Parallel()
	repo := newFakeRepo()
	mw := idem.New(repo, nil, nil)
	counter := 0
	router := buildRouter(mw, &counter)

	rec := doPost(router, "")
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if counter != 1 {
		t.Fatalf("handler executed %d times, want 1", counter)
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.claims != 0 {
		t.Fatalf("claims = %d, want 0 (no key => no claim)", repo.claims)
	}
}

func TestIdempotencyReplaysSameKey(t *testing.T) {
	t.Parallel()
	repo := newFakeRepo()
	mw := idem.New(repo, nil, nil)
	counter := 0
	router := buildRouter(mw, &counter)

	first := doPost(router, "k-1")
	if first.Code != http.StatusCreated {
		t.Fatalf("first status = %d, want %d", first.Code, http.StatusCreated)
	}
	if counter != 1 {
		t.Fatalf("handler executed %d times after first, want 1", counter)
	}
	if repo.completesCount() != 1 {
		t.Fatalf("completes after first = %d, want 1", repo.completesCount())
	}

	second := doPost(router, "k-1")
	if second.Code != http.StatusCreated {
		t.Fatalf("replay status = %d, want %d", second.Code, http.StatusCreated)
	}
	if counter != 1 {
		t.Fatalf("handler executed %d times after replay, want still 1 (no duplicate)", counter)
	}
	if second.Body.String() != first.Body.String() {
		t.Fatalf(
			"replay body = %q, want first body %q (same saved response)",
			second.Body.String(), first.Body.String(),
		)
	}
}

func TestIdempotencyInFlightReturnsConflict(t *testing.T) {
	t.Parallel()
	repo := newFakeRepo()
	mw := idem.New(repo, nil, nil)
	counter := 0
	router := buildRouter(mw, &counter)

	// Эмулируем «в полёте»: первый запрос с ключом ещё не завершился.
	repo.mu.Lock()
	repo.inflight[scope("k-x", 7, http.MethodPost, "/tasks")] = true
	repo.mu.Unlock()

	rec := doPost(router, "k-x")
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d (409 in-flight)", rec.Code, http.StatusConflict)
	}
	if counter != 0 {
		t.Fatalf("handler executed %d times, want 0 (must not re-execute)", counter)
	}
}

func TestIdempotencyFiveHundredReleasesKey(t *testing.T) {
	t.Parallel()
	repo := newFakeRepo()
	mw := idem.New(repo, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(setUser(7))
	router.Use(mw.Handler())
	router.POST("/boom", func(c *gin.Context) {
		c.String(http.StatusInternalServerError, "boom")
	})

	req := httptest.NewRequest(http.MethodPost, "/boom", nil)
	req.Header.Set(headerKey, "k-5")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.releases != 1 {
		t.Fatalf("releases = %d, want 1 (5xx must release the key)", repo.releases)
	}
	if _, ok := repo.complete[scope("k-5", 7, http.MethodPost, "/boom")]; ok {
		t.Fatalf("5xx response must NOT be cached for replay")
	}
}

func TestIdempotencyWithoutAuthPassesThrough(t *testing.T) {
	t.Parallel()
	repo := newFakeRepo()
	mw := idem.New(repo, nil, nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	// Без setUser — userctx.GetUserID вернёт ошибку, мидлварь пропускает.
	router.Use(mw.Handler())
	counter := 0
	router.POST("/tasks", func(c *gin.Context) {
		counter++
		c.Status(http.StatusCreated)
	})

	req := httptest.NewRequest(http.MethodPost, "/tasks", nil)
	req.Header.Set(headerKey, "k")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", rec.Code)
	}
	if counter != 1 {
		t.Fatalf("handler executed %d times, want 1 (no auth => pass through)", counter)
	}
	repo.mu.Lock()
	defer repo.mu.Unlock()
	if repo.claims != 0 {
		t.Fatalf("claims = %d, want 0", repo.claims)
	}
}
