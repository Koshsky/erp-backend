-- name: ListActivePresets :many
SELECT id, name, description
FROM rbac_presets
ORDER BY name;

-- name: ListActivePresetRules :many
SELECT id, preset, resource, action, scope, updated_by, updated_at
FROM rbac_preset_rules;

-- name: UpsertPresetRule :one
INSERT INTO rbac_preset_rules (preset, resource, action, scope, updated_by)
VALUES (@preset::text, @resource::text, @action::text, @scope::text, @updated_by)
ON CONFLICT (preset, resource, action) DO UPDATE
SET scope = EXCLUDED.scope,
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING id, preset, resource, action, scope, updated_by, updated_at;

-- name: DeletePresetRule :exec
DELETE FROM rbac_preset_rules
WHERE id = @id::bigint;

-- name: DeleteAllPresetRules :exec
DELETE FROM rbac_preset_rules;

-- name: ListActiveRoutePolicies :many
SELECT name, kind, params, active, updated_by, updated_at
FROM rbac_route_policies
WHERE active = TRUE;

-- name: UpsertRoutePolicy :one
INSERT INTO rbac_route_policies (name, kind, params, active, updated_by)
VALUES (@name::text, @kind::text, @params, @active, @updated_by)
ON CONFLICT (name) DO UPDATE
SET kind = EXCLUDED.kind,
    params = EXCLUDED.params,
    active = EXCLUDED.active,
    updated_by = EXCLUDED.updated_by,
    updated_at = NOW()
RETURNING name, kind, params, active, updated_by, updated_at;

-- name: DeleteRoutePolicy :exec
DELETE FROM rbac_route_policies
WHERE name = @name::text;

-- name: DeleteAllRoutePolicies :exec
DELETE FROM rbac_route_policies;

-- name: UpsertPreset :one
INSERT INTO rbac_presets (name, description)
VALUES (@name::text, @description::text)
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description,
    updated_at = NOW()
RETURNING id, name, description;

-- name: UpdatePresetDescription :one
UPDATE rbac_presets
SET description = @description::text, updated_at = NOW()
WHERE name = @name::text
RETURNING id, name, description;

-- name: DeletePreset :exec
DELETE FROM rbac_presets
WHERE name = @name::text;

-- name: DeletePresetRulesByPreset :exec
DELETE FROM rbac_preset_rules
WHERE preset = @preset::text;

-- name: ClearPresetOnUsers :exec
-- Removes the preset ref from users after its deletion (base rights vanish;
-- individual overrides survive).
UPDATE users
SET preset = NULL, updated_at = NOW()
WHERE preset = @preset::text;

-- =============================================
-- Per-user permission overrides
-- =============================================

-- name: ListUserPermissions :many
SELECT id, user_id, resource, action, scope, granted, updated_by, updated_at
FROM user_permissions
WHERE user_id = @user_id::bigint;

-- name: DeleteAllUserPermissions :exec
DELETE FROM user_permissions
WHERE user_id = @user_id::bigint;

-- name: InsertUserPermission :one
INSERT INTO user_permissions (user_id, resource, action, scope, granted, updated_by)
VALUES (@user_id::bigint, @resource::text, @action::text, @scope::text, @granted, @updated_by)
RETURNING id, user_id, resource, action, scope, granted, updated_by, updated_at;

-- name: FindUserPreset :one
SELECT preset
FROM users
WHERE id = @user_id::bigint;

-- name: ListUserPrincipals :many
-- Active users with their preset — the base of the in-memory principal snapshot.
SELECT id AS user_id, preset
FROM users;

-- name: ListAllUserPermissions :many
-- All active per-user overrides — grouped by user_id in the snapshot loader.
SELECT user_id, resource, action, scope, granted
FROM user_permissions;