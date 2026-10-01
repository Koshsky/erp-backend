//nolint:testpackage // tests exercise unexported guards (isBuiltinPreset, findRulePreset)
package service

import (
	"testing"

	"github.com/Koshsky/erp-backend/internal/authz/repository/sqlc"
	userdomain "github.com/Koshsky/erp-backend/internal/user/domain"
)

// Only the admin entry is a code invariant (deleting it would clear
// users.preset for every administrator and lock them out of /rbac/*); every
// other preset — including the seeded dp/rp/vp/worker — stays editable and
// deletable through the API (renamed by an admin or deleted with its rules).
func TestIsBuiltinPreset(t *testing.T) {
	t.Parallel()
	if !isBuiltinPreset(userdomain.PresetAdmin) {
		t.Error("пресет «admin» должен считаться встроенным")
	}
	for _, name := range []string{"dp", "rp", "vp", "worker", "auditor", "", "Admin", "admin1", "Админ"} {
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

// Preset names accept letters of any script (latin/cyrillic), digits, «-» and
// «_» — the admin UI must be able to create a preset named «Аудит».
func TestValidPresetNameCharset(t *testing.T) {
	t.Parallel()
	allowed := []string{"auditor", "Аудит", "пресет_1", "СБЫТ-2026", "mixed123_"}
	for _, name := range allowed {
		if !validPresetName(name) {
			t.Errorf("validPresetName(%q) = false, want true", name)
		}
	}
	forbidden := []string{"", "a b", "a.b", "пресет,1", "a/b", "«квота»"}
	for _, name := range forbidden {
		if validPresetName(name) {
			t.Errorf("validPresetName(%q) = true, want false", name)
		}
	}
}
