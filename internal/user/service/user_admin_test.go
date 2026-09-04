//nolint:testpackage // service-level tests construct UserService directly with a stub repository (unexported fields)
package service

import (
	"context"
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/Koshsky/erp-backend/internal/tracing"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/internal/user/dto"
	userctx "github.com/Koshsky/erp-backend/internal/userctx"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

// stubRepo is a minimal UserRepository for the user-admin rule tests.
type stubRepo struct {
	users      map[int64]*userdomain.User
	adminCount int64
}

func newStubRepo(users ...*userdomain.User) *stubRepo {
	r := &stubRepo{users: map[int64]*userdomain.User{}}
	for _, u := range users {
		if u != nil {
			r.users[u.ID] = u
		}
	}
	for _, u := range r.users {
		if u.PresetName() == userdomain.PresetAdmin {
			r.adminCount++
		}
	}
	return r
}

func (r *stubRepo) CreateUser(_ context.Context, user userdomain.User) (*userdomain.User, error) {
	created := user
	if created.ID == 0 {
		created.ID = int64(len(r.users) + 100)
	}
	r.users[created.ID] = &created
	return &created, nil
}

func (r *stubRepo) CreateUserWithPermissions(
	_ context.Context,
	user userdomain.User,
	_ []userdomain.UserPermission,
	_ int64,
) (*userdomain.User, error) {
	return r.CreateUser(context.Background(), user)
}

func (r *stubRepo) FindUser(_ context.Context, id int64) (*userdomain.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, errors.NotFound("not found")
	}
	cp := *u
	return &cp, nil
}

func (r *stubRepo) FindUserByUsername(_ context.Context, username string) (*userdomain.User, error) {
	for _, u := range r.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.NotFound("not found")
}

func (r *stubRepo) UsernameExists(_ context.Context, username string) (bool, error) {
	for _, u := range r.users {
		if u.Username == username {
			return true, nil
		}
	}
	return false, nil
}

func (r *stubRepo) UpdateUser(_ context.Context, user userdomain.User) (*userdomain.User, error) {
	cp := user
	r.users[user.ID] = &cp
	return &cp, nil
}

func (r *stubRepo) UpdatePassword(_ context.Context, _ int64, _ string) error { return nil }
func (r *stubRepo) DeleteUser(_ context.Context, _ int64) error               { return nil }

func (r *stubRepo) ListUsers(
	_ context.Context,
	_ int64, _ string, _ string, _ int64, _ int, _ int,
) ([]userdomain.User, error) {
	return nil, nil
}

func (r *stubRepo) CountUsers(_ context.Context, _ int64, _ string, presetFilter string, _ int64) (int64, error) {
	if presetFilter == userdomain.PresetAdmin {
		return r.adminCount, nil
	}
	return 0, nil
}

func (r *stubRepo) ListAllUsers(_ context.Context) ([]userdomain.User, error) { return nil, nil }

func (r *stubRepo) ListStates(_ context.Context, _ int64, _, _ time.Time) ([]userdomain.UserState, error) {
	return nil, nil
}

func (r *stubRepo) SetStateRange(_ context.Context, _, _ int64, _, _ time.Time) error { return nil }
func (r *stubRepo) DeleteStateRange(_ context.Context, _ int64, _, _ time.Time, _ *int64) error {
	return nil
}

func newStubService(repo *stubRepo) *UserService {
	return &UserService{
		logger:     slog.New(slog.DiscardHandler),
		repository: repo,
		mapper:     &UserMapper{},
		validator:  &UserValidator{},
		tracer:     tracing.New(nil),
	}
}

// vp is a non-admin caller principal (the vp preset).
func vp(userID int64) userctx.UserContext {
	return userctx.UserContext{ID: userID, Preset: "vp", Admin: false}
}

// admin is an admin caller principal.
func admin(userID int64) userctx.UserContext {
	return userctx.UserContext{ID: userID, Preset: "admin", Admin: true}
}

// presetPtr returns a pointer to a preset code.
func presetPtr(p string) *string {
	return &p
}

// Creating a user is covered by the grantable user_admin.create right, but
// preset assignment stays admin-only: non-admin callers may only create workers
// (a missing preset defaults to worker). Prevents privilege escalation through
// the now-grantable create endpoint.
func TestCreateUserPresetRule(t *testing.T) {
	t.Parallel()
	svc := newStubService(newStubRepo())
	req := dto.CreateUserRequest{
		LastName: "И", FirstName: "И", Username: "u1", PasswordHash: "hash",
	}

	for _, preset := range []string{
		userdomain.PresetProjectDirector, userdomain.PresetProjectManager, userdomain.PresetAdmin,
	} {
		r := req
		r.Preset = presetPtr(preset)
		if _, loopErr := svc.CreateUserWithCreds(
			context.Background(), r, vp(10),
		); !errors.IsForbidden(loopErr) {
			t.Errorf("не-админ назначает пресет %q: err=%v; want forbidden", preset, loopErr)
		}
	}

	empty := req
	created, createErr := svc.CreateUserWithCreds(context.Background(), empty, vp(10))
	if createErr != nil {
		t.Fatalf("не-админ без пресета: %v; want default worker", createErr)
	}
	if presetName(created.User.Preset) != userdomain.PresetWorker {
		t.Errorf("пустой пресет: got %q; want worker", presetName(created.User.Preset))
	}

	worker := req
	worker.Preset = presetPtr(userdomain.PresetWorker)
	if _, wErr := svc.CreateUserWithCreds(
		context.Background(), worker, vp(10),
	); wErr != nil {
		t.Errorf("не-админ создаёт worker: %v; want ok", wErr)
	}

	adminReq := req
	adminReq.Preset = presetPtr(userdomain.PresetAdmin)
	if _, aErr := svc.CreateUserWithCreds(context.Background(), adminReq, admin(1)); aErr != nil {
		t.Errorf("админ создаёт админа: %v; want ok", aErr)
	}
}

// Admin may create a user with individual overrides in one call; a non-admin
// cannot, and invalid overrides are rejected.
func TestCreateUserWithPermissions(t *testing.T) {
	t.Parallel()
	svc := newStubService(newStubRepo())
	req := dto.CreateUserRequest{
		LastName: "И", FirstName: "И", Username: "u1", PasswordHash: "hash",
		Preset: presetPtr(userdomain.PresetWorker),
		Permissions: []dto.UserPermissionInput{
			{Resource: "project", Action: "view", Scope: "own", Granted: true},
			{Resource: "task", Action: "delete", Granted: false},
		},
	}

	// Non-admin with permissions -> forbidden.
	if _, err := svc.CreateUserWithCreds(context.Background(), req, vp(10)); !errors.IsForbidden(err) {
		t.Errorf("не-админ с правами при создании: err=%v; want forbidden", err)
	}

	// Admin with valid overrides -> ok.
	r := req
	created, err := svc.CreateUserWithCreds(context.Background(), r, admin(1))
	if err != nil {
		t.Fatalf("админ с правами при создании: %v; want ok", err)
	}
	if presetName(created.User.Preset) != userdomain.PresetWorker {
		t.Errorf("пресет созданного: got %q; want worker", presetName(created.User.Preset))
	}

	// Invalid resource -> bad request.
	bad := req
	bad.Permissions = []dto.UserPermissionInput{{Resource: "no_such", Action: "view", Scope: "all", Granted: true}}
	if _, berr := svc.CreateUserWithCreds(
		context.Background(),
		bad,
		admin(1),
	); berr == nil ||
		errors.StatusCode(berr) != http.StatusBadRequest {
		t.Errorf("неверный ресурс: err=%v; want bad request", berr)
	}

	// Duplicate (resource/action) -> bad request.
	dup := req
	dup.Permissions = []dto.UserPermissionInput{
		{Resource: "project", Action: "view", Scope: "own", Granted: true},
		{Resource: "project", Action: "view", Scope: "all", Granted: true},
	}
	if _, derr := svc.CreateUserWithCreds(
		context.Background(),
		dup,
		admin(1),
	); derr == nil ||
		errors.StatusCode(derr) != http.StatusBadRequest {
		t.Errorf("дубликат права: err=%v; want bad request", derr)
	}
}

// Manager assignment is part of the user_admin.update right (grantable):
// a non-admin caller may set/clear the manager; the self and cycle rules stay.
func TestUpdateManagerGrantable(t *testing.T) {
	t.Parallel()
	repo := newStubRepo(
		&userdomain.User{
			ID:        1,
			Username:  "vp1",
			LastName:  "В",
			FirstName: "П",
			Preset:    presetPtr(userdomain.PresetProcessOwner),
		},
		&userdomain.User{
			ID:        2,
			Username:  "w2",
			LastName:  "Р",
			FirstName: "а",
			Preset:    presetPtr(userdomain.PresetWorker),
		},
	)
	svc := newStubService(repo)

	mgr := int64(1)
	if _, err := svc.UpdateManager(context.Background(), 2, &mgr); err != nil {
		t.Fatalf("назначение руководителя не-админом: %v; want ok", err)
	}
	if got := repo.users[2].ManagerID; got == nil || *got != mgr {
		t.Errorf("менеджер не сохранён: got %v; want %d", got, mgr)
	}

	self := int64(2)
	if _, err := svc.UpdateManager(context.Background(), 2, &self); !errors.IsValidationError(err) {
		t.Errorf("сам себе: err=%v; want validation error", err)
	}

	if _, err := svc.UpdateManager(context.Background(), 2, nil); err != nil {
		t.Fatalf("сброс руководителя: %v; want ok", err)
	}
	if repo.users[2].ManagerID != nil {
		t.Errorf("менеджер не сброшен: %v", repo.users[2].ManagerID)
	}
}

// UpdateUser: non-admin honors manager_id in the body (user_admin.update), but
// a preset change stays admin-only (escalation protection); the last active
// admin cannot be demoted.
func TestUpdateUserManagerAndPreset(t *testing.T) {
	t.Parallel()
	repo := newStubRepo(
		&userdomain.User{
			ID:        1,
			Username:  "vp1",
			LastName:  "В",
			FirstName: "П",
			Preset:    presetPtr(userdomain.PresetProcessOwner),
		},
		&userdomain.User{
			ID:        2,
			Username:  "w2",
			LastName:  "Р",
			FirstName: "а",
			Preset:    presetPtr(userdomain.PresetWorker),
		},
	)
	svc := newStubService(repo)

	mgr := int64(1)
	if _, err := svc.UpdateUser(
		context.Background(), 2, dto.UpdateUserRequest{ManagerID: &mgr}, vp(99), 99,
	); err != nil {
		t.Fatalf("manager_id в теле от не-админа: %v; want ok", err)
	}
	if got := repo.users[2].ManagerID; got == nil || *got != mgr {
		t.Errorf("manager_id не применён: got %v; want %d", got, mgr)
	}

	preset := userdomain.PresetProjectManager
	if _, err := svc.UpdateUser(
		context.Background(), 2, dto.UpdateUserRequest{Preset: &preset}, vp(99), 99,
	); !errors.IsForbidden(err) {
		t.Errorf("смена пресета не-админом: err=%v; want forbidden", err)
	}

	if _, err := svc.UpdateUser(
		context.Background(), 2, dto.UpdateUserRequest{Preset: &preset}, admin(1), 99,
	); err != nil {
		t.Fatalf("смена пресета админом: %v; want ok", err)
	}
	if presetName(repo.users[2].Preset) != preset {
		t.Errorf("пресет не применён: got %q; want %q", presetName(repo.users[2].Preset), preset)
	}
}

// The last active admin cannot be demoted (preset change guard).
func TestUpdateUserLastAdminGuard(t *testing.T) {
	t.Parallel()
	repo := newStubRepo(
		&userdomain.User{
			ID:        1,
			Username:  "a1",
			LastName:  "Ад",
			FirstName: "м",
			Preset:    presetPtr(userdomain.PresetAdmin),
		},
		&userdomain.User{
			ID:        2,
			Username:  "w2",
			LastName:  "Р",
			FirstName: "а",
			Preset:    presetPtr(userdomain.PresetWorker),
		},
	)
	svc := newStubService(repo)

	preset := userdomain.PresetWorker
	if _, err := svc.UpdateUser(
		context.Background(), 1, dto.UpdateUserRequest{Preset: &preset}, admin(1), 1,
	); !errors.IsValidationError(err) {
		t.Errorf("снятие последнего админа самим себе: err=%v; want validation", err)
	}

	// A different admin may not demote the last admin either.
	if _, err := svc.UpdateUser(
		context.Background(), 1, dto.UpdateUserRequest{Preset: &preset}, admin(7), 7,
	); !errors.IsValidationError(err) {
		t.Errorf("снятие последнего админа другим админом: err=%v; want validation", err)
	}
}
