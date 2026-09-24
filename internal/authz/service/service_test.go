//nolint:testpackage // tests exercise unexported guards (isBuiltinPreset, findRulePreset)
package service

import (
	"testing"

	"github.com/Koshsky/erp-backend/internal/authz/repository/sqlc"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
)

// builtinPresets — the seeded catalog of V10__rbac_policies.sql.
//
//nolint:gochecknoglobals // rule registry
var builtinPresets = []string{"admin", "dp", "rp", "vp", "worker"}

// The seeded presets are immutable through the RBAC admin API: deleting admin
// (or any other built-in) would clear users.preset and lock administrators out
// of /rbac/*, while an upsert would let the built-in admin entry be clobbered.
func TestIsBuiltinPreset(t *testing.T) {
	t.Parallel()
	for _, name := range builtinPresets {
		if !isBuiltinPreset(name) {
			t.Errorf("пресет %q должен считаться встроенным", name)
		}
	}
	// Runtime presets created through the API stay editable and deletable.
	for _, name := range []string{"auditor", "", "Admin", "admin1", "Админ"} {
		if isBuiltinPreset(name) {
			t.Errorf("пресет %q не является встроенным", name)
		}
	}
}

// The admin preset is protected by the same guard the admin API relies on.
func TestIsBuiltinPresetAdmin(t *testing.T) {
	t.Parallel()
	if !isBuiltinPreset(userdomain.PresetAdmin) {
		t.Fatal("встроенный админ-пресет должен быть защищён")
	}
}

// A matrix row is resolved to its preset by id: the admin preset's rows must be
// recognizable so DELETE /rbac/preset-rules/{id} can refuse them.
func TestFindRulePreset(t *testing.T) {
	t.Parallel()
	rules := []sqlc.ListActivePresetRulesRow{
		{ID: 1, Preset: userdomain.PresetAdmin, Resource: "rbac_config", Action: "view", Scope: "all"},
		{ID: 2, Preset: userdomain.PresetProcessOwner, Resource: "task", Action: "view", Scope: "parent"},
	}
	if preset, ok := findRulePreset(rules, 1); !ok || preset != userdomain.PresetAdmin {
		t.Fatalf("строка 1 должна принадлежать пресету admin, получено %q (найдено: %v)", preset, ok)
	}
	if preset, ok := findRulePreset(rules, 2); !ok || preset != userdomain.PresetProcessOwner {
		t.Fatalf("строка 2 должна принадлежать пресету vp, получено %q (найдено: %v)", preset, ok)
	}
	if preset, ok := findRulePreset(rules, 3); ok {
		t.Fatalf("несуществующая строка не должна находиться, получено %q", preset)
	}
	if _, ok := findRulePreset(nil, 1); ok {
		t.Fatal("пустой набор строк не должен содержать строк")
	}
}

// The refusal is a 400 with the preset name (the established style: the
// handler passes service errors to response.Error unchanged).
func TestBuiltinPresetErr(t *testing.T) {
	t.Parallel()
	err := builtinPresetErr(userdomain.PresetAdmin)
	if err == nil {
		t.Fatal("ошибка защиты встроенного пресета не должна быть nil")
	}
	if err.Error() == "" {
		t.Fatal("сообщение об ошибке защиты встроенного пресета не должно быть пустым")
	}
}
