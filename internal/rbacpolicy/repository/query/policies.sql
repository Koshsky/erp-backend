-- name: ListActiveRoles :many
SELECT id, name, description
FROM rbac_roles
WHERE deleted_at IS NULL
ORDER BY name;

-- name: ListActiveRules :many
SELECT id, role, resource, action, scope, updated_by, updated_at
FROM rbac_role_rules
WHERE deleted_at IS NULL;

-- name: UpsertRule :one
INSERT INTO rbac_role_rules (role, resource, action, scope, updated_by)
VALUES (@role::text, @resource::text, @action::text, @scope::text, @updated_by)
ON CONFLICT (role, resource, action) DO UPDATE
SET scope = EXCLUDED.scope,
    updated_by = EXCLUDED.updated_by,
    deleted_at = NULL, -- upsert «оживляет» soft-deleted строку (reset/повторная выдача права)
    updated_at = NOW()
RETURNING id, role, resource, action, scope, updated_by, updated_at;

-- name: SoftDeleteRule :exec
UPDATE rbac_role_rules
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @id::bigint AND deleted_at IS NULL;

-- name: SoftDeleteAllRules :exec
UPDATE rbac_role_rules
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
    deleted_at = NULL, -- upsert «оживляет» soft-deleted строку (reset)
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