package service_test

import (
	"strings"
	"testing"

	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
	"github.com/Koshsky/erp-backend/internal/user/service"
)

// Presets are configured by the rbac_presets catalog: only the form is
// validated, existence is checked by the FK — so a new preset passes, a
// missing one (nil) is allowed, and an overly long one does not.
func TestValidateUserPreset(t *testing.T) {
	t.Parallel()
	v := &service.UserValidator{}
	base := userdomain.User{LastName: "И", FirstName: "И", Username: "u1"}

	// No preset is valid: a user may have only individual permissions.
	if err := v.ValidateUser(&base); err != nil {
		t.Fatalf("пустой пресет должен проходить валидацию: %v", err)
	}
	long := base
	longPreset := strings.Repeat("a", 33)
	long.Preset = &longPreset
	if err := v.ValidateUser(&long); err == nil {
		t.Fatal("пресет длиннее 32 должен давать ошибку")
	}
	for _, preset := range []string{"auditor", "vp", "worker"} {
		r := base
		r.Preset = &preset
		if err := v.ValidateUser(&r); err != nil {
			t.Errorf("пресет %q должен проходить валидацию: %v", preset, err)
		}
	}
}
