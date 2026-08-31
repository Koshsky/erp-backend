//nolint:testpackage // service-level tests construct UserService directly with a stub repository (unexported fields)
package service

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/Koshsky/erp-backend/internal/tracing"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/internal/user/dto"
	"github.com/Koshsky/erp-backend/pkg/errors"
)

// stubRepo is a minimal UserRepository for the user-admin rule tests.
type stubRepo struct {
	users map[int64]*userdomain.User
}

func newStubRepo(users ...*userdomain.User) *stubRepo {
	r := &stubRepo{users: map[int64]*userdomain.User{}}
	for _, u := range users {
		if u != nil {
			r.users[u.ID] = u
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

func (r *stubRepo) CountUsers(_ context.Context, _ int64, _ string, _ string, _ int64) (int64, error) {
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

// Creating a user is covered by the grantable user_admin.create right, but
// role assignment stays admin-only: non-admin callers may only create workers
// (an empty role defaults to worker). Prevents privilege escalation through
// the now-grantable create endpoint.
func TestCreateUserRoleRule(t *testing.T) {
	t.Parallel()
	svc := newStubService(newStubRepo())
	req := dto.CreateUserRequest{
		LastName: "И", FirstName: "И", Username: "u1", PasswordHash: "hash",
	}

	for _, role := range []string{
		userdomain.ProjectDirector, userdomain.ProjectManager, userdomain.Admin,
	} {
		r := req
		r.Role = role
		if _, loopErr := svc.CreateUserWithCreds(
			context.Background(), r, userdomain.ProcessOwner,
		); !errors.IsForbidden(loopErr) {
			t.Errorf("не-админ создаёт роль %q: err=%v; want forbidden", role, loopErr)
		}
	}

	empty := req
	empty.Role = ""
	created, createErr := svc.CreateUserWithCreds(context.Background(), empty, userdomain.ProcessOwner)
	if createErr != nil {
		t.Fatalf("не-админ с пустой ролью: %v; want default worker", createErr)
	}
	if created.User.Role != userdomain.Worker {
		t.Errorf("пустая роль: got %q; want worker", created.User.Role)
	}

	worker := req
	worker.Role = userdomain.Worker
	if _, wErr := svc.CreateUserWithCreds(
		context.Background(), worker, userdomain.ProcessOwner,
	); wErr != nil {
		t.Errorf("не-админ создаёт worker: %v; want ok", wErr)
	}

	admin := req
	admin.Role = userdomain.Admin
	if _, aErr := svc.CreateUserWithCreds(context.Background(), admin, userdomain.Admin); aErr != nil {
		t.Errorf("админ создаёт админа: %v; want ok", aErr)
	}
}

// Manager assignment is part of the user_admin.update right (grantable):
// a non-admin caller may set/clear the manager; the self and cycle rules stay.
func TestUpdateManagerGrantable(t *testing.T) {
	t.Parallel()
	repo := newStubRepo(
		&userdomain.User{ID: 1, Username: "vp1", LastName: "В", FirstName: "П", Role: userdomain.ProcessOwner},
		&userdomain.User{ID: 2, Username: "w2", LastName: "Р", FirstName: "а", Role: userdomain.Worker},
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
// a role change stays admin-only (escalation protection).
func TestUpdateUserManagerAndRole(t *testing.T) {
	t.Parallel()
	repo := newStubRepo(
		&userdomain.User{ID: 1, Username: "vp1", LastName: "В", FirstName: "П", Role: userdomain.ProcessOwner},
		&userdomain.User{ID: 2, Username: "w2", LastName: "Р", FirstName: "а", Role: userdomain.Worker},
	)
	svc := newStubService(repo)

	mgr := int64(1)
	if _, err := svc.UpdateUser(
		context.Background(), 2, dto.UpdateUserRequest{ManagerID: &mgr}, userdomain.ProcessOwner, 99,
	); err != nil {
		t.Fatalf("manager_id в теле от не-админа: %v; want ok", err)
	}
	if got := repo.users[2].ManagerID; got == nil || *got != mgr {
		t.Errorf("manager_id не применён: got %v; want %d", got, mgr)
	}

	role := userdomain.ProjectManager
	if _, err := svc.UpdateUser(
		context.Background(), 2, dto.UpdateUserRequest{Role: &role}, userdomain.ProcessOwner, 99,
	); !errors.IsForbidden(err) {
		t.Errorf("смена роли не-админом: err=%v; want forbidden", err)
	}

	if _, err := svc.UpdateUser(
		context.Background(), 2, dto.UpdateUserRequest{Role: &role}, userdomain.Admin, 99,
	); err != nil {
		t.Fatalf("смена роли админом: %v; want ok", err)
	}
	if repo.users[2].Role != role {
		t.Errorf("роль не применена: got %q; want %q", repo.users[2].Role, role)
	}
}
