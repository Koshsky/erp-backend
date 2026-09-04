package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Koshsky/erp-backend/internal/policies"
	repo "github.com/Koshsky/erp-backend/internal/user/repository"

	"github.com/Koshsky/erp-backend/internal/security/creds"
	"github.com/Koshsky/erp-backend/internal/security/hasher"
	tracingpkg "github.com/Koshsky/erp-backend/internal/tracing"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/internal/user/dto"
	userctx "github.com/Koshsky/erp-backend/internal/userctx"
	"github.com/Koshsky/erp-backend/pkg/date"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

// RBACReloader is an optional hook to refresh the in-memory RBAC snapshot
// (presets + per-user rights) after a user mutation; implemented by the
// rbacpolicy PolicyStore. Nil — rights changes propagate via the TTL refresh.
type RBACReloader interface {
	Reload(ctx context.Context) error
}

type UserService struct {
	logger     *slog.Logger
	repository UserRepository
	mapper     *UserMapper
	validator  *UserValidator
	tracer     *tracingpkg.Tracer
	rbac       RBACReloader
}

// maxManagerDepth — guard against an infinite loop while walking the manager hierarchy.
const maxManagerDepth = 1000

// NewUserService builds the UserService service.
func NewUserService(
	logger *slog.Logger,
	tracer *tracingpkg.Tracer,
	r *repo.UserRepository,
	rbac RBACReloader,
) *UserService {
	return &UserService{
		logger:     logger,
		repository: r,
		mapper:     &UserMapper{},
		validator:  &UserValidator{},
		tracer:     tracer,
		rbac:       rbac,
	}
}

func (s *UserService) FindUserByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	ctx, end := s.tracer.Start(ctx, "user.FindUserByID")
	defer end(nil)
	return s.FindUser(ctx, id)
}

func (s *UserService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	ctx, end := s.tracer.Start(ctx, "user.ChangePassword")
	defer end(nil)

	// Complexity policy — checked before the old password (AD-09): the format of
	// the new password reveals nothing about the current one.
	if err := creds.ValidatePassword(newPassword); err != nil {
		return err
	}

	user, err := s.FindUserByID(ctx, userID)
	if err != nil {
		return errors.NotFound("user not found")
	}

	if err = hasher.Compare(user.PasswordHash, oldPassword); err != nil {
		return errors.NewValidationError("invalid current password")
	}

	newHash, err := hasher.Hash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password")
	}

	return s.repository.UpdatePassword(ctx, userID, newHash)
}

// CreateUserWithCreds creates a user and returns the generated password (if any)
// for the admin UI. The password is shown exactly once. Preset assignment on
// creation is admin-only (user_admin.create is grantable, but preset escalation
// stays an admin privilege — non-admin callers may only create workers).
func (s *UserService) CreateUserWithCreds(
	ctx context.Context,
	req dto.CreateUserRequest,
	caller userctx.UserContext,
) (*dto.CreateUserResult, error) {
	ctx, end := s.tracer.Start(ctx, "user.CreateUser")
	defer end(nil)
	return s.createUserInternal(ctx, req, caller)
}

// createUserInternal creates a user; when credentials are missing they are
// generated (the plaintext password is returned for one-time display).
func (s *UserService) createUserInternal(
	ctx context.Context,
	req dto.CreateUserRequest,
	caller userctx.UserContext,
) (*dto.CreateUserResult, error) {
	var generated string
	// Non-admin holders of user_admin.create may only create workers (a missing
	// preset defaults to worker); the other presets require the admin bypass.
	if !caller.Admin {
		if req.Preset == nil {
			req.Preset = stringPtr(userdomain.PresetWorker)
		}
		if *req.Preset != userdomain.PresetWorker {
			return nil, errors.ErrForbidden
		}
	}
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		username, err := s.generateUsername(ctx, req.LastName, presetName(req.Preset))
		if err != nil {
			return nil, err
		}
		req.Username = username
	}
	if req.PasswordHash == "" {
		raw, err := creds.RandomPassword()
		if err != nil {
			return nil, err
		}
		hash, err := hasher.Hash(raw)
		if err != nil {
			return nil, err
		}
		req.PasswordHash = hash
		generated = raw
	}
	// Individual overrides at creation: admin-only, validated like
	// /rbac/users/{id}/permissions.
	perms, permsErr := s.validateCreatePermissions(req.Permissions, caller)
	if permsErr != nil {
		return nil, permsErr
	}

	user := s.mapper.ToDomainFromCreate(req)
	if err := s.validator.ValidateUser(&user); err != nil {
		return nil, err
	}

	created, err := s.repository.CreateUserWithPermissions(ctx, user, perms, caller.ID)
	if err != nil {
		return nil, err
	}
	s.refreshRBAC(ctx)

	return &dto.CreateUserResult{User: *s.mapper.ToDTO(created), Password: generated}, nil
}

// scopeAllCode — the "all" ownership scope code (grant default; mirrors the
// policies codec).
const scopeAllCode = "all"

// validateCreatePermissions validates the individual overrides of a created
// user: admin-only (like preset assignment) and each entry checked the same
// way as the per-user permissions endpoint (resource/action/scope codecs,
// scope applicability, no duplicates). Nil on an empty list.
func (s *UserService) validateCreatePermissions(
	req []dto.UserPermissionInput,
	caller userctx.UserContext,
) ([]userdomain.UserPermission, error) {
	if len(req) == 0 {
		return nil, nil
	}
	if !caller.Admin {
		return nil, errors.ErrForbidden
	}
	seen := map[string]bool{}
	out := make([]userdomain.UserPermission, 0, len(req))
	for _, p := range req {
		res, ok := policies.ParseResource(p.Resource)
		if !ok {
			return nil, errors.BadRequest("неизвестный ресурс " + p.Resource)
		}
		if _, okAction := policies.ParseAction(p.Action); !okAction {
			return nil, errors.BadRequest("неизвестное действие " + p.Action)
		}
		key := p.Resource + "/" + p.Action
		if seen[key] {
			return nil, errors.BadRequest("дублируется право " + key)
		}
		seen[key] = true
		scope := scopeAllCode
		if p.Granted {
			parsed, okScope := policies.ParseScope(p.Scope)
			if !okScope || parsed == policies.ScopeNone {
				return nil, errors.BadRequest("недопустимая зона " + p.Scope + " (all|own|parent|ancestor)")
			}
			if !policies.ScopeApplicable(res, parsed) {
				return nil, errors.BadRequest("зона " + p.Scope + " неприменима к ресурсу " + p.Resource)
			}
			scope = p.Scope
		}
		out = append(out, userdomain.UserPermission{
			Resource: p.Resource,
			Action:   p.Action,
			Scope:    scope,
			Granted:  p.Granted,
		})
	}
	return out, nil
}

// generateUsername builds a unique login: transliteration of the last name
// (last_name); if taken, appends a numeric suffix; if there is nothing to
// transliterate — falls back to prefix+random suffix.
func (s *UserService) generateUsername(ctx context.Context, name, preset string) (string, error) {
	prefix := "user_"
	if preset == userdomain.PresetWorker {
		prefix = "worker_"
	}

	base := creds.Transliterate(name)
	if base == "" {
		suffix, err := creds.RandomUsernameSuffix()
		if err != nil {
			return "", err
		}
		return prefix + suffix, nil
	}

	username := base
	for i := 2; ; i++ {
		exists, err := s.repository.UsernameExists(ctx, username)
		if err != nil {
			return "", err
		}
		if !exists {
			return username, nil
		}
		username = fmt.Sprintf("%s%d", base, i)
	}
}

// ResetPassword generates a new random password for the user and returns it once.
func (s *UserService) ResetPassword(ctx context.Context, id int64) (*dto.ResetPasswordResponse, error) {
	ctx, end := s.tracer.Start(ctx, "user.ResetPassword")
	defer end(nil)

	user, err := s.repository.FindUser(ctx, id)
	if err != nil || user == nil {
		return nil, errors.NotFound("user not found")
	}

	raw, err := creds.RandomPassword()
	if err != nil {
		return nil, err
	}
	hash, err := hasher.Hash(raw)
	if err != nil {
		return nil, err
	}
	if err = s.repository.UpdatePassword(ctx, id, hash); err != nil {
		return nil, err
	}

	return &dto.ResetPasswordResponse{Password: raw}, nil
}

func (s *UserService) FindUserByUsername(ctx context.Context, username string) (*dto.UserResponse, error) {
	ctx, end := s.tracer.Start(ctx, "user.FindUserByUsername")
	defer end(nil)

	user, err := s.repository.FindUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(user), nil
}

func (s *UserService) FindUser(ctx context.Context, id int64) (*dto.UserResponse, error) {
	ctx, end := s.tracer.Start(ctx, "user.FindUser")
	defer end(nil)

	user, err := s.repository.FindUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.NotFound("user not found")
	}
	return s.mapper.ToDTO(user), nil
}

// UpdateUser updates a user. Preset assignment is the only admin-only business
// rule (escalation protection via checkPresetChange); manager_id and profile
// fields are covered by the user_admin.update right.
func (s *UserService) UpdateUser(
	ctx context.Context,
	id int64,
	req dto.UpdateUserRequest,
	caller userctx.UserContext,
	callerID int64,
) (*dto.UserResponse, error) {
	ctx, end := s.tracer.Start(ctx, "user.UpdateUser")
	defer end(nil)

	user, err := s.repository.FindUser(ctx, id)
	if err != nil || user == nil {
		return nil, errors.NotFound("user not found")
	}

	if err = s.checkPresetChange(ctx, user, req.Preset, caller, callerID); err != nil {
		return nil, err
	}
	if req.ManagerID != nil {
		if err = s.validateManager(ctx, id, req.ManagerID); err != nil {
			return nil, err
		}
	}

	s.mapper.ApplyUpdateToDomain(user, req)
	user.Username = strings.TrimSpace(user.Username)
	if err = s.validator.ValidateUser(user); err != nil {
		return nil, err
	}

	updated, err := s.repository.UpdateUser(ctx, *user)
	if err != nil {
		return nil, err
	}
	s.refreshRBAC(ctx)

	return s.mapper.ToDTO(updated), nil
}

// checkPresetChange validates a preset change: only admin, not on self, and
// never removing the last active admin.
func (s *UserService) checkPresetChange(
	ctx context.Context,
	user *userdomain.User,
	newPreset *string,
	caller userctx.UserContext,
	callerID int64,
) error {
	if newPreset == nil {
		return nil
	}
	if !caller.Admin {
		return errors.ErrForbidden
	}
	if user.ID == callerID {
		return errors.NewValidationError("нельзя менять пресет прав самому себе")
	}
	if *newPreset == userdomain.PresetAdmin || presetName(user.Preset) != userdomain.PresetAdmin {
		return nil
	}
	admins, err := s.repository.CountUsers(ctx, 0, scopeAllCode, userdomain.PresetAdmin, 0)
	if err != nil {
		return err
	}
	if admins <= 1 {
		return errors.NewValidationError("нельзя снять последнего админа")
	}
	return nil
}

// refreshRBAC best-effort refresh of the in-memory RBAC snapshot after a user
// mutation (preset/account changes must take effect immediately on this
// instance; the TTL refresh heals on failure).
func (s *UserService) refreshRBAC(ctx context.Context) {
	if s.rbac == nil {
		return
	}
	if err := s.rbac.Reload(ctx); err != nil {
		s.logger.WarnContext(ctx, "user: обновление RBAC-снапшота не удалось (исправится фоновым TTL)", "error", err)
	}
}

// presetName unwraps a preset pointer ("" — none).
func presetName(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// stringPtr returns a pointer to a string.
func stringPtr(s string) *string {
	return &s
}

// UpdateManager explicitly sets (or clears) a user's manager. Covered by the
// user_admin.update right (grantable); the manager cannot be the user themself
// and cycles are rejected.
func (s *UserService) UpdateManager(
	ctx context.Context,
	id int64,
	managerID *int64,
) (*dto.UserResponse, error) {
	ctx, end := s.tracer.Start(ctx, "user.UpdateManager")
	defer end(nil)

	user, err := s.repository.FindUser(ctx, id)
	if err != nil || user == nil {
		return nil, errors.NotFound("user not found")
	}
	if err = s.validateManager(ctx, id, managerID); err != nil {
		return nil, err
	}

	user.ManagerID = managerID
	if err = s.validator.ValidateUser(user); err != nil {
		return nil, err
	}
	updated, err := s.repository.UpdateUser(ctx, *user)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTO(updated), nil
}

// validateManager checks that the manager assignment is valid: it is not the
// user themself, the manager exists and is active, and the assignment does not
// create a circular dependency (walking up the manager_id chain).
// managerID == nil — clearing, allowed.
func (s *UserService) validateManager(ctx context.Context, userID int64, managerID *int64) error {
	if managerID == nil {
		return nil
	}
	if *managerID == userID {
		return errors.NewValidationError("руководитель не может быть самим пользователем")
	}

	manager, err := s.repository.FindUser(ctx, *managerID)
	if err != nil || manager == nil {
		return errors.NotFound("руководитель не найден")
	}

	cur := manager.ManagerID
	depth := 0
	for cur != nil {
		depth++
		if depth > maxManagerDepth {
			return errors.NewValidationError("иерархия руководителей слишком глубокая")
		}
		if *cur == userID {
			return errors.NewValidationError("кольцевая зависимость в руководстве не допускается")
		}
		u, ferr := s.repository.FindUser(ctx, *cur)
		if ferr != nil || u == nil {
			return errors.NotFound("руководитель не найден")
		}
		cur = u.ManagerID
	}
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	ctx, end := s.tracer.Start(ctx, "user.DeleteUser")
	defer end(nil)
	return s.repository.DeleteUser(ctx, id)
}

// ListUsers returns a paged list of users; visibility by manager is enforced
// in the middleware (vp sees only their own subordinates).
func (s *UserService) ListUsers(
	ctx context.Context,
	userID int64,
	viewScope string,
	presetFilter string,
	managerID int64,
	limit, offset int,
) ([]dto.UserResponse, int64, error) {
	ctx, end := s.tracer.Start(ctx, "user.ListUsers")
	defer end(nil)

	users, err := s.repository.ListUsers(ctx, userID, viewScope, presetFilter, managerID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repository.CountUsers(ctx, userID, viewScope, presetFilter, managerID)
	if err != nil {
		return nil, 0, err
	}
	return s.mapper.ToDTOs(users), total, nil
}

// ListAllUsers returns every active user (unscoped; used for owner pickers).
func (s *UserService) ListAllUsers(ctx context.Context) ([]dto.UserResponse, error) {
	ctx, end := s.tracer.Start(ctx, "user.ListAllUsers")
	defer end(nil)

	users, err := s.repository.ListAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return s.mapper.ToDTOs(users), nil
}

func (s *UserService) ListStates(
	ctx context.Context,
	userID int64,
	start, end date.Date,
) ([]dto.UserStateResponse, error) {
	ctx, finish := s.tracer.Start(ctx, "user.ListStates")
	defer finish(nil)

	if err := s.validator.ValidatePositiveID(userID, "user_id"); err != nil {
		return nil, err
	}
	if err := s.validator.ValidateDayRange(start, end); err != nil {
		return nil, err
	}

	if err := s.ensureUserExists(ctx, userID); err != nil {
		return nil, err
	}

	states, err := s.repository.ListStates(ctx, userID, start.Time(), end.Time())
	if err != nil {
		return nil, err
	}
	return s.mapper.ToStateDTOs(states), nil
}

func (s *UserService) SetDays(
	ctx context.Context,
	userID int64,
	req dto.SetDaysRequest,
) error {
	ctx, end := s.tracer.Start(ctx, "user.SetDays")
	defer end(nil)

	if err := s.validator.ValidatePositiveID(userID, "user_id"); err != nil {
		return err
	}
	if err := s.validator.ValidatePositiveID(req.StateID, "state_id"); err != nil {
		return err
	}
	if err := s.validator.ValidateDayRange(req.StartDate, req.EndDate); err != nil {
		return err
	}

	if err := s.ensureUserExists(ctx, userID); err != nil {
		return err
	}

	return s.repository.SetStateRange(ctx, userID, req.StateID, req.StartDate.Time(), req.EndDate.Time())
}

func (s *UserService) DeleteDays(
	ctx context.Context,
	userID int64,
	start, end date.Date,
	stateID *int64,
) error {
	ctx, finish := s.tracer.Start(ctx, "user.DeleteDays")
	defer finish(nil)

	if err := s.validator.ValidatePositiveID(userID, "user_id"); err != nil {
		return err
	}
	if stateID != nil {
		if err := s.validator.ValidatePositiveID(*stateID, "state_id"); err != nil {
			return err
		}
	}
	if err := s.validator.ValidateDayRange(start, end); err != nil {
		return err
	}

	if err := s.ensureUserExists(ctx, userID); err != nil {
		return err
	}

	return s.repository.DeleteStateRange(ctx, userID, start.Time(), end.Time(), stateID)
}

// ensureUserExists verifies the user exists (404 otherwise).
func (s *UserService) ensureUserExists(ctx context.Context, userID int64) error {
	user, err := s.repository.FindUser(ctx, userID)
	if err != nil {
		if errors.IsNotFoundError(err) {
			return errors.ErrUserNotFound
		}
		return err
	}
	if user == nil {
		return errors.ErrUserNotFound
	}
	return nil
}
