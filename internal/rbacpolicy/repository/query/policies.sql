-- name: ListActivePresets :many
SELECT id, name, description
FROM rbac_presets
WHERE deleted_at IS NULL
ORDER BY name;

-- name: ListActivePresetRules :many
SELECT id, preset, resource, action, scope, updated_by, updated_at
FROM rbac_preset_rules
WHERE deleted_at IS NULL;

-- name: UpsertPresetRule :one
INSERT INTO rbac_preset_rules (preset, resource, action, scope, updated_by)
VALUES (@preset::text, @resource::text, @action::text, @scope::text, @updated_by)
ON CONFLICT (preset, resource, action) DO UPDATE
SET scope = EXCLUDED.scope,
    updated_by = EXCLUDED.updated_by,
    deleted_at = NULL, -- upsert "revives" a soft-deleted row (reset/re-grant of a permission)
    updated_at = NOW()
RETURNING id, preset, resource, action, scope, updated_by, updated_at;

-- name: SoftDeletePresetRule :exec
UPDATE rbac_preset_rules
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @id::bigint AND deleted_at IS NULL;

-- name: SoftDeleteAllPresetRules :exec
UPDATE rbac_preset_rules
SET deleted_at = NOW(), updated_at = NOW()
WHERE deleted_at IS NULL;

-- name: ListActiveRoutePolicies :many
SELECT name, kind, params, active, updated_by, updated_at
FROM rbac_route_policies
WHERE deleted_at IS NULL AND active = TRUE;

-- name: UpsertRoutePolicy :one
INSERT INTO rbac_route_policies (name, kind, params, active, updated_by)
VALUES (@name::text, @kind::text, @params, @active, @updated_by)
ON CONFLICT (name) DO UPDATE
SET kind = EXCLUDED.kind,
    params = EXCLUDED.params,
    active = EXCLUDED.active,
    updated_by = EXCLUDED.updated_by,
    deleted_at = NULL, -- upsert "revives" a soft-deleted row (reset)
    updated_at = NOW()
RETURNING name, kind, params, active, updated_by, updated_at;

-- name: SoftDeleteRoutePolicy :exec
UPDATE rbac_route_policies
SET deleted_at = NOW(), updated_at = NOW()
WHERE name = @name::text AND deleted_at IS NULL;

-- name: SoftDeleteAllRoutePolicies :exec
UPDATE rbac_route_policies
SET deleted_at = NOW(), updated_at = NOW()
WHERE deleted_at IS NULL;

-- name: UpsertPreset :one
INSERT INTO rbac_presets (name, description)
VALUES (@name::text, @description::text)
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description,
    deleted_at = NULL, -- upsert "revives" a soft-deleted preset
    updated_at = NOW()
RETURNING id, name, description;

-- name: UpdatePresetDescription :one
UPDATE rbac_presets
SET description = @description::text, updated_at = NOW()
WHERE name = @name::text AND deleted_at IS NULL
RETURNING id, name, description;

-- name: SoftDeletePreset :exec
UPDATE rbac_presets
SET deleted_at = NOW(), updated_at = NOW()
WHERE name = @name::text AND deleted_at IS NULL;

-- name: SoftDeletePresetRulesByPreset :exec
UPDATE rbac_preset_rules
SET deleted_at = NOW(), updated_at = NOW()
WHERE preset = @preset::text AND deleted_at IS NULL;

-- name: ClearPresetOnUsers :exec
-- Removes the preset ref from users after its soft delete (base rights vanish;
-- individual overrides survive).
UPDATE users
SET preset = NULL, updated_at = NOW()
WHERE preset = @preset::text AND deleted_at IS NULL;

-- =============================================
-- Per-user permission overrides
-- =============================================

-- name: ListUserPermissions :many
SELECT id, user_id, resource, action, scope, granted, updated_by, updated_at
FROM user_permissions
WHERE user_id = @user_id::bigint AND deleted_at IS NULL;

-- name: SoftDeleteAllUserPermissions :exec
UPDATE user_permissions
SET deleted_at = NOW(), updated_at = NOW()
WHERE user_id = @user_id::bigint AND deleted_at IS NULL;

-- name: InsertUserPermission :one
INSERT INTO user_permissions (user_id, resource, action, scope, granted, updated_by)
VALUES (@user_id::bigint, @resource::text, @action::text, @scope::text, @granted, @updated_by)
RETURNING id, user_id, resource, action, scope, granted, updated_by, updated_at;

-- name: FindUserPreset :one
SELECT preset
FROM users
WHERE id = @user_id::bigint AND deleted_at IS NULL;

-- name: ListUserPrincipals :many
-- Active users with their preset — the base of the in-memory principal snapshot.
SELECT id AS user_id, preset
FROM users
WHERE deleted_at IS NULL;

-- name: ListAllUserPermissions :many
-- All active per-user overrides — grouped by user_id in the snapshot loader.
SELECT user_id, resource, action, scope, granted
FROM user_permissions
WHERE deleted_at IS NULL;