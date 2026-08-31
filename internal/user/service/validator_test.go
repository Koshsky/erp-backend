package service_test

import (
	"strings"
	"testing"

	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/internal/user/service"
)

// Roles are configured by the rbac_roles catalog: only the form is validated,
// existence is checked by the FK (V17) — so a new role passes, while an empty
// or overly long one does not.
func TestValidateUserRole(t *testing.T) {
	t.Parallel()
	v := &service.UserValidator{}
	base := userdomain.User{LastName: "И", FirstName: "И", Username: "u1"}

	if err := v.ValidateUser(&base); err == nil {
		t.Fatal("пустая роль должна давать ошибку")
	}
	long := base
	long.Role = strings.Repeat("a", 33)
	if err := v.ValidateUser(&long); err == nil {
		t.Fatal("роль длиннее 32 должна давать ошибку")
	}
	for _, role := range []string{"auditor", "vp", "worker"} {
		r := base
		r.Role = role
		if err := v.ValidateUser(&r); err != nil {
			t.Errorf("роль %q должна проходить валидацию: %v", role, err)
		}
	}
}
