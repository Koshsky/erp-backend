-- =============================================
-- Presets: split the single `name` column into `tag` (system access code —
-- the unique identity stored in users.preset / rbac_preset_rules.preset) and
-- `name` (human-readable display name). Historic rows: tag = old name;
-- built-in presets get Russian display names; custom presets take their
-- description as the display name when it was filled.
-- The preset column of users / rules used to be an FK on rbac_presets(name);
-- the catalog is now referenced by tag values with integrity enforced by the
-- RBAC service, so the FK constraints are removed for good.
-- Idempotent for fresh databases (V1 already carries the new columns).
-- =============================================

ALTER TABLE rbac_presets
    ADD COLUMN IF NOT EXISTS tag TEXT;

ALTER TABLE rbac_presets_deleted
    ADD COLUMN IF NOT EXISTS tag TEXT;

-- Backfill tags from the historic names (live table + archive).
UPDATE rbac_presets SET tag = name WHERE tag IS NULL;
UPDATE rbac_presets_deleted SET tag = name WHERE tag IS NULL;

-- From now on every preset has a tag.
ALTER TABLE rbac_presets
    ALTER COLUMN tag SET NOT NULL;
ALTER TABLE rbac_presets_deleted
    ALTER COLUMN tag SET NOT NULL;

-- Drop the old foreign keys (they referenced rbac_presets(name) historically
-- and were briefly re-pointed at tag — no FK on the tag column by design).
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_preset_fkey;
ALTER TABLE rbac_preset_rules
    DROP CONSTRAINT IF EXISTS rbac_preset_rules_preset_fkey;

-- Single UNIQUE constraint on tags (fresh databases declare UNIQUE in V1;
-- upgraded databases get the constraint here).
ALTER TABLE rbac_presets
    DROP CONSTRAINT IF EXISTS rbac_presets_tag_key;
ALTER TABLE rbac_presets
    ADD CONSTRAINT rbac_presets_tag_key UNIQUE (tag);

-- Built-in presets: Russian display names.
UPDATE rbac_presets SET name = 'Администратор' WHERE tag = 'admin';
UPDATE rbac_presets SET name = 'Директор проектов' WHERE tag = 'dp';
UPDATE rbac_presets SET name = 'Руководитель проекта' WHERE tag = 'rp';
UPDATE rbac_presets SET name = 'Владелец процесса' WHERE tag = 'vp';
UPDATE rbac_presets SET name = 'Работник' WHERE tag = 'worker';

-- Custom presets: the historic name was the code itself; take the description
-- as the display name when it was filled (empty name stays the code).
UPDATE rbac_presets
SET name = NULLIF(description, '')
WHERE tag NOT IN ('admin', 'dp', 'rp', 'vp', 'worker')
  AND name = tag
  AND NULLIF(description, '') IS NOT NULL;