package service_test

import (
	"strings"
	"testing"

	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/internal/user/service"
)

// Роли конфигурируются каталогом rbac_roles: валидируется только форма,
// существование проверяет FK (V17) — поэтому новая роль проходит, а пустая
// и слишком длинная — нет.
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
