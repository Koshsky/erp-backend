package rbac_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/middleware/rbac"
	"github.com/Koshsky/erp-backend/internal/policies"
	userctx "github.com/Koshsky/erp-backend/internal/userctx"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

//nolint:gochecknoinits // gin test mode
func init() {
	gin.SetMode(gin.TestMode)
}

const (
	testAdmin = "admin"
	testRP    = "rp"
	testVP    = "vp"
)

// setUser mimics AuthMiddleware: puts the user into the gin context.
func setUser(role string, id int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(userctx.KeyUser, userctx.UserContext{ID: id, Role: role})
		c.Next()
	}
}

// chainOwners is a "chain" resolver: task/process → {Project:1, Process:2}.
func chainOwners(_ context.Context, id int64) (rbac.Owners, error) {
	if id%10 == 0 {
		return rbac.Owners{}, errors.ErrNotFound
	}
	return rbac.Owners{ProjectOwner: 1, ProcessOwner: 2}, nil
}

// rowOwner is a "row owner" resolver: resource/worker → owner=id.
func rowOwner(_ context.Context, id int64) (rbac.Owners, error) {
	if id%10 == 0 {
		return rbac.Owners{}, errors.ErrNotFound
	}
	return rbac.Owners{Owner: id}, nil
}

func testData() rbac.Data {
	return rbac.Data{
		ProjectOwners:    chainOwners,
		ProcessOwners:    chainOwners,
		TaskOwners:       chainOwners,
		MilestoneOwners:  chainOwners,
		AssignmentOwners: chainOwners,
		ResourceOwners:   rowOwner,
		WorkerOwners:     rowOwner,
	}
}

func TestCheckEntity(t *testing.T) {
	t.Parallel()
	mw := rbac.ProvideMiddleware(nil, testData(), []rbac.Policy{
		{Name: "task.view", Check: policies.EntityCheck(rbac.ResourceTask, policies.ActionView)},
		{Name: "task.update", Check: policies.EntityCheck(rbac.ResourceTask, policies.ActionUpdate)},
	})

	cases := []struct {
		name       string
		role       string
		path       string
		policy     string
		withUser   bool
		wantStatus int
	}{
		{"no user → 401", "", "/1", "task.view", false, http.StatusUnauthorized},
		{"bad id → 400", testAdmin, "/abc", "task.view", true, http.StatusBadRequest},
		{"not found → 404", testAdmin, "/10", "task.view", true, http.StatusNotFound},
		{"rp views foreign project → 404", testRP, "/1", "task.view", true, http.StatusNotFound},
		{"rp update → 403 (view-only)", testRP, "/1", "task.update", true, http.StatusForbidden},
		{"vp views own process → 200", testVP, "/1", "task.view", true, http.StatusOK},
		{"vp updates own → 200", testVP, "/1", "task.update", true, http.StatusOK},
		{"admin view → 200", testAdmin, "/1", "task.view", true, http.StatusOK},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			router := gin.New()
			if tc.withUser {
				router.Use(setUser(tc.role, 2))
			}
			router.GET("/task/:id", mw.Check(tc.policy), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"ok": true})
			})

			req := httptest.NewRequest(http.MethodGet, "/task"+tc.path, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != tc.wantStatus {
				t.Errorf("GET %s: got %d, want %d (body: %s)", tc.path, rec.Code, tc.wantStatus, rec.Body.String())
			}
		})
	}
}

func TestCheckAssignmentCreate(t *testing.T) {
	t.Parallel()
	mw := rbac.ProvideMiddleware(nil, testData(), []rbac.Policy{
		{Name: "assignment.create", Check: createAssignmentForTest},
	})

	router := gin.New()
	router.Use(setUser(testVP, 2))
	router.POST("/assignment", mw.Check("assignment.create"), func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"ok": true})
	})

	cases := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"task and resource share the owner", `{"task_id":1,"resource_id":2,"quantity":1}`, http.StatusCreated},
		{"foreign resource → 403", `{"task_id":1,"resource_id":5,"quantity":1}`, http.StatusForbidden},
		{"missing task → 404", `{"task_id":10,"resource_id":2,"quantity":1}`, http.StatusNotFound},
		{"malformed body → 400", `not-json`, http.StatusBadRequest},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodPost, "/assignment", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)
			if rec.Code != tc.wantStatus {
				t.Errorf("%s: got %d, want %d (body: %s)", tc.name, rec.Code, tc.wantStatus, rec.Body.String())
			}
		})
	}
}

// createAssignmentForTest mirrors the internal/policies rule: the matrix check
// by the task plus the shared-owner business rule.
func createAssignmentForTest(rc *rbac.CheckCtx) error {
	taskID, err := rc.BodyID("task_id")
	if err != nil {
		return err
	}
	resourceID, err := rc.BodyID("resource_id")
	if err != nil {
		return err
	}

	taskOwners, err := rc.Owners(rbac.ResourceTask, taskID)
	if err != nil {
		return err
	}
	if !policies.Authorize(rc.User.Role, rbac.ResourceAssignment, policies.ActionCreate, taskOwners, rc.User.ID) {
		return errors.ErrForbidden
	}

	resourceOwners, err := rc.Owners(rbac.ResourceResource, resourceID)
	if err != nil {
		return err
	}
	if !taskOwners.SharesOwner(resourceOwners) {
		return errors.ErrForbidden
	}
	return nil
}
