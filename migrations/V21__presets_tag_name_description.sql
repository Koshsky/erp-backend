-- =============================================
-- Presets: split the single `name` column into `tag` (system access code —
-- the identity referenced by users.preset / rbac_preset_rules.preset) and
-- `name` (human-readable display name). Historic rows: tag = old name;
-- built-in presets get Russian display names; custom presets take their
-- description as the display name when it was filled.
-- Idempotent for fresh databases (V1 already carries the new columns).
-- =============================================

ALTER TABLE rbac_presets
    ADD COLUMN IF NOT EXISTS tag TEXT NOT NULL DEFAULT '';

ALTER TABLE rbac_presets_deleted
    ADD COLUMN IF NOT EXISTS tag TEXT NOT NULL DEFAULT '';

-- Backfill tags from the historic names (live table + archive).
UPDATE rbac_presets SET tag = name WHERE tag = '';
UPDATE rbac_presets_deleted SET tag = name WHERE tag = '';

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

-- Re-point the foreign keys at the tag (identity) column. Values stay the same
-- (tags were copied from the historic names), only the referenced column moves.
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_preset_fkey;
ALTER TABLE users
    ADD CONSTRAINT users_preset_fkey FOREIGN KEY (preset) REFERENCES rbac_presets(tag);

ALTER TABLE rbac_preset_rules
    DROP CONSTRAINT IF EXISTS rbac_preset_rules_preset_fkey;
ALTER TABLE rbac_preset_rules
    ADD CONSTRAINT rbac_preset_rules_preset_fkey FOREIGN KEY (preset) REFERENCES rbac_presets(tag);

-- Unique index on tags (partial — empty tags exist only transiently here).
CREATE UNIQUE INDEX IF NOT EXISTS idx_rbac_presets_tag
    ON rbac_presets(tag) WHERE tag <> '';