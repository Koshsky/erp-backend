package rbac

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Koshsky/erp-backend/internal/common/ctx"
	"github.com/Koshsky/erp-backend/internal/common/helpers"
	"github.com/Koshsky/erp-backend/internal/common/response"
	"github.com/Koshsky/erp-backend/internal/validator"
)

// ErrInvalidBody indicates the request body cannot be parsed (no parent id).
var ErrInvalidBody = errors.New("invalid request body")

// ErrBadRequest indicates invalid request parameters (the id from the URL).
var ErrBadRequest = errors.New("invalid request")

// ErrForbidden and ErrNotFound are denial sentinels (adapted from validator).
var (
	ErrForbidden = validator.ErrForbidden
	ErrNotFound  = validator.ErrNotFound
)

// ResolveByID returns the owner chain of an entity by its id.
type ResolveByID func(ctx context.Context, id int64) (Owners, error)

// Data holds the owner resolver functions (OwnerChain implementations in the
// entity repositories). They are injected into the engine at startup.
type Data struct {
	ProjectOwners    ResolveByID
	ProcessOwners    ResolveByID
	TaskOwners       ResolveByID
	MilestoneOwners  ResolveByID
	AssignmentOwners ResolveByID
	ResourceOwners   ResolveByID
	EmployeeOwners   ResolveByID
}

// resolve returns the resolver for a resource (nil — the resource has no owner).
func (d Data) resolve(rsrc Resource) ResolveByID {
	switch rsrc {
	case ResourceProject:
		return d.ProjectOwners
	case ResourceProcess:
		return d.ProcessOwners
	case ResourceTask:
		return d.TaskOwners
	case ResourceMilestone:
		return d.MilestoneOwners
	case ResourceAssignment:
		return d.AssignmentOwners
	case ResourceResource:
		return d.ResourceOwners
	case ResourceEmployee:
		return d.EmployeeOwners
	case ResourceState:
		// States have no owner.
		return nil
	default:
		return nil
	}
}

// Policy is an access rule for a route. All logic (role, owners, business
// rules) is hidden inside Check.
type Policy struct {
	Name  string
	Check func(rc *CheckCtx) error
}

// CheckCtx is the policy evaluation context.
type CheckCtx struct {
	C    *gin.Context
	User ctx.UserContext
	Data Data
	body []byte
}

// Body returns the request body (already restored for the handler).
func (rc *CheckCtx) Body() []byte {
	return rc.body
}

// ParamID returns the entity id from the URL (:id).
func (rc *CheckCtx) ParamID() (int64, error) {
	id, err := strconv.ParseInt(rc.C.Param("id"), 10, 64)
	if err != nil {
		return 0, ErrBadRequest
	}
	return id, nil
}

// BodyID extracts a numeric id from the body by key (0 if missing/null).
func (rc *CheckCtx) BodyID(key string) (int64, error) {
	var obj map[string]json.RawMessage
	if err := json.Unmarshal(rc.body, &obj); err != nil {
		return 0, ErrInvalidBody
	}
	raw, ok := obj[key]
	if !ok || len(raw) == 0 || string(raw) == "null" {
		return 0, nil
	}
	var v int64
	if err := json.Unmarshal(raw, &v); err != nil {
		return 0, ErrInvalidBody
	}
	return v, nil
}

// Owners returns the owner chain of an entity (via the resolver from Data).
func (rc *CheckCtx) Owners(rsrc Resource, id int64) (Owners, error) {
	resolve := rc.Data.resolve(rsrc)
	if resolve == nil {
		return Owners{}, validator.ErrNotFound
	}
	return resolve(rc.C.Request.Context(), id)
}

// Middleware is the policy engine.
type Middleware struct {
	logger *slog.Logger
	data   Data
	byName map[string]Policy
}

func New(logger *slog.Logger, data Data, policies []Policy) *Middleware {
	byName := make(map[string]Policy, len(policies))
	for _, p := range policies {
		byName[p.Name] = p
	}
	return &Middleware{logger: logger, data: data, byName: byName}
}

// Check runs the policy by name.
func (m *Middleware) Check(name string) gin.HandlerFunc {
	policy, ok := m.byName[name]
	if !ok {
		panic("rbac: unknown policy " + name)
	}
	return func(c *gin.Context) {
		user, found := getUser(c)
		if !found {
			return
		}
		body, err := readBody(c)
		if err != nil {
			response.BadRequest(c, "invalid request body")
			c.Abort()
			return
		}
		rc := &CheckCtx{C: c, User: user, Data: m.data, body: body}
		if err = policy.Check(rc); err != nil {
			m.abort(c, err)
			return
		}
		c.Next()
	}
}

// abortError writes the HTTP response according to the policy error type.
func (m *Middleware) abort(c *gin.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidBody):
		response.BadRequest(c, "invalid request body")
	case errors.Is(err, ErrBadRequest):
		response.BadRequest(c, "invalid request")
	case errors.Is(err, ErrForbidden):
		response.Forbidden(c, "forbidden")
	case validator.IsNotFoundError(err):
		response.NotFound(c, "not found")
	default:
		response.InternalError(c, m.logger, err.Error(), err)
	}
	c.Abort()
}

// readBody reads the request body and restores it for the following handlers.
func readBody(c *gin.Context) ([]byte, error) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	return body, nil
}

// getUser takes the user from the context (set by AuthMiddleware).
func getUser(c *gin.Context) (ctx.UserContext, bool) {
	user, err := helpers.GetUser(c)
	if err != nil {
		response.Unauthorized(c, "authentication required")
		c.Abort()
		return ctx.UserContext{}, false
	}
	return user, true
}
