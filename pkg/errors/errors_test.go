package errors_test

import (
	stdErrors "errors"
	"net/http"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	errapi "github.com/Koshsky/erp-backend/pkg/errors"
)

func TestConflictCodeAndStatus(t *testing.T) {
	t.Parallel()
	err := errapi.Conflict("project with this code already exists")

	var de *errapi.DomainError
	if !stdErrors.As(err, &de) {
		t.Fatalf("Conflict() returned %T, want *DomainError", err)
	}
	if de.StatusCode() != http.StatusConflict {
		t.Fatalf("StatusCode = %d, want %d", de.StatusCode(), http.StatusConflict)
	}
	if de.ErrorCode() != errapi.CodeConflict {
		t.Fatalf("ErrorCode = %s, want %s", de.ErrorCode(), errapi.CodeConflict)
	}
	if de.Code.String() != "CONFLICT" {
		t.Fatalf("wire code = %q, want %q", de.Code.String(), "CONFLICT")
	}
}

func TestCodeForConflictStatus(t *testing.T) {
	t.Parallel()
	if got := errapi.CodeFor(http.StatusConflict); got != errapi.CodeConflict {
		t.Fatalf("CodeFor(409) = %s, want %s", got, errapi.CodeConflict)
	}
}

func TestConflictIsSentinel(t *testing.T) {
	t.Parallel()
	if !errapi.IsConflictError(errapi.Conflict("collision")) {
		t.Fatalf("IsConflictError(Conflict()) = false, want true")
	}
}

func TestStatusCodeMapsConflict(t *testing.T) {
	t.Parallel()
	if got := errapi.StatusCode(errapi.Conflict("x")); got != http.StatusConflict {
		t.Fatalf("StatusCode(Conflict) = %d, want %d", got, http.StatusConflict)
	}
}

func TestFromPgInvalidParam(t *testing.T) {
	t.Parallel()
	pgErr := &pgconn.PgError{Code: "22023", Message: "Сроки процесса выходят за границы проекта"}
	got := errapi.FromPgInvalidParam(pgErr)
	if errapi.StatusCode(got) != http.StatusBadRequest {
		t.Fatalf("StatusCode = %d, want 400", errapi.StatusCode(got))
	}
	if got.Error() != pgErr.Message {
		t.Errorf("Message = %q, want %q", got.Error(), pgErr.Message)
	}
	// Другие коды и чужие ошибки — без изменений.
	other := &pgconn.PgError{Code: "23505", Message: "dup"}
	if !stdErrors.Is(errapi.FromPgInvalidParam(other), other) {
		t.Errorf("23505 не должен мапиться")
	}
	sentinel := stdErrors.New("boom")
	if !stdErrors.Is(errapi.FromPgInvalidParam(sentinel), sentinel) {
		t.Errorf("сторонняя ошибка не должна мапиться")
	}
}

func TestMapPgConstraint(t *testing.T) {
	t.Parallel()
	// 23505 unique (username) -> 409
	dup := errapi.MapPgConstraint(&pgconn.PgError{Code: "23505", ConstraintName: "users_username_unique_active"})
	if errapi.StatusCode(dup) != http.StatusConflict {
		t.Errorf("23505: StatusCode = %d, want 409", errapi.StatusCode(dup))
	}
	// 23503 role fk -> 400 "неизвестная роль"
	role := errapi.MapPgConstraint(&pgconn.PgError{Code: "23503", ConstraintName: "users_role_fk"})
	if errapi.StatusCode(role) != http.StatusBadRequest || !strings.Contains(role.Error(), "каталоге ролей") {
		t.Errorf("23503 role: got %v", role)
	}
	// 23514 check -> 400
	chk := errapi.MapPgConstraint(&pgconn.PgError{Code: "23514", ConstraintName: "tasks_dates_check"})
	if errapi.StatusCode(chk) != http.StatusBadRequest {
		t.Errorf("23514: StatusCode = %d, want 400", errapi.StatusCode(chk))
	}
	// чужие ошибки без изменений
	other := &pgconn.PgError{Code: "22023", Message: "x"}
	if !stdErrors.Is(errapi.MapPgConstraint(other), other) {
		t.Errorf("22023 не должен мапиться")
	}
	sentinel := stdErrors.New("boom")
	if !stdErrors.Is(errapi.MapPgConstraint(sentinel), sentinel) {
		t.Errorf("сторонняя ошибка не должна мапиться")
	}
}
