package errors_test

import (
	stdErrors "errors"
	"net/http"
	"testing"

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
