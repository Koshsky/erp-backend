-- =============================================================
-- Manual one-shot upgrade: roles -> presets + per-user permissions
-- =============================================================
-- WARNING: the repo's migration policy edits V1/V10 in place, so a DB that has
-- already applied the old V1..V14 cannot flyway-migrate (checksum mismatch) and
-- must be recreated (make reset) OR upgraded in-place with THIS script, run once
-- against the old schema BEFORE starting the new backend.
--
-- The script renames the RBAC catalog/matrix to presets, moves users.role into
-- users.preset, and creates the per-user permission overrides table. Run it in a
-- transaction; back up first. After it completes, the new backend can start
-- (it reads presets + user_permissions only).
-- =============================================================

BEGIN;

-- 1. Rename the role catalog and the rights matrix to presets.
ALTER TABLE rbac_roles RENAME TO rbac_presets;
ALTER TABLE rbac_role_rules RENAME TO rbac_preset_rules;
ALTER TABLE rbac_preset_rules RENAME COLUMN role TO preset;

-- 2. Index/trigger names follow the table names (Postgres renames them
--    implicitly when the table is renamed; the FK constraint name does too).

-- 3. users.role -> users.preset (nullable; the role value becomes the preset).
ALTER TABLE users ADD COLUMN preset TEXT REFERENCES rbac_presets(name);
UPDATE users SET preset = role WHERE deleted_at IS NULL;

-- 4. Drop the old role column and its FK (the FK was renamed along with the
--    table, but the users.role column still points at rbac_presets(name)).
ALTER TABLE users DROP COLUMN role;

-- 5. Per-user permission overrides (explicit grants / revokes that shadow the
--    preset rule for the same resource+action).
CREATE TABLE user_permissions (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    resource    TEXT NOT NULL,
    action      TEXT NOT NULL,
    scope       TEXT NOT NULL DEFAULT 'all',
    granted     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    updated_by  BIGINT,
    CHECK (NOT granted OR scope IN ('all', 'own', 'parent', 'ancestor'))
);
CREATE UNIQUE INDEX idx_user_permissions_active
ON user_permissions(user_id, resource, action) WHERE deleted_at IS NULL;

-- 6. Soft-delete protection (same trigger function as the other tables, V5).
CREATE TRIGGER block_hard_delete_on_user_permissions
BEFORE DELETE ON user_permissions
FOR EACH ROW EXECUTE FUNCTION block_hard_delete();

COMMIT;