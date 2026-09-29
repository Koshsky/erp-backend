-- =============================================
-- TASK DEPENDENCIES (SCHEDULING LINKS BETWEEN TASKS)
-- =============================================
-- The tables/columns are also declared in V1 (the sqlc schema source), so a
-- fresh database already has them and the CREATEs below are no-ops there;
-- existing deployments get them from this migration.
--
-- A dependency is a directed edge between two top-level tasks of the same
-- process (enforced in the service): the successor (task_id) may not
-- start/finish before the predecessor (depends_on_task_id). The exact rule
-- depends on the link type:
--   fs — finish-to-start:  successor must start after the predecessor ends,
--   ss — start-to-start:   successor must start after the predecessor starts,
--   ff — finish-to-finish: successor must end after the predecessor ends,
--   sf — start-to-finish:  successor must end after the predecessor starts.

CREATE TABLE IF NOT EXISTS task_dependencies (
	id BIGSERIAL PRIMARY KEY,
	task_id BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	depends_on_task_id BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	type TEXT NOT NULL DEFAULT 'fs',
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	CONSTRAINT task_dependencies_no_self CHECK (task_id <> depends_on_task_id),
	CONSTRAINT task_dependencies_unique UNIQUE (task_id, depends_on_task_id),
	CHECK (type IN ('fs', 'ss', 'ff', 'sf'))
);

-- Successor lookup (incoming constraints of a task) and cascade scans.
CREATE INDEX IF NOT EXISTS idx_task_dependencies_task
ON task_dependencies(task_id);
CREATE INDEX IF NOT EXISTS idx_task_dependencies_depends_on
ON task_dependencies(depends_on_task_id);

-- Archive (soft delete by move, V5): the live table holds active rows only.
CREATE TABLE IF NOT EXISTS task_dependencies_deleted (
    id                  BIGINT NOT NULL,
    task_id             BIGINT NOT NULL,
    depends_on_task_id  BIGINT NOT NULL,
    type                TEXT NOT NULL DEFAULT 'fs',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Archive trigger. A fresh database already has it from V5 (edited), so guard
-- the CREATE with a DO block to stay idempotent for both fresh and existing
-- deployments.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_trigger WHERE tgname = 'trg_archive_task_dependencies'
    ) THEN
        CREATE TRIGGER trg_archive_task_dependencies
        BEFORE DELETE ON task_dependencies
        FOR EACH ROW EXECUTE FUNCTION fn_archive_row();
    END IF;
END $$;

-- =============================================
-- RBAC ROUTE POLICIES (name → kind + params)
-- =============================================
-- Nested under /task/:id (like task.comment.*): the entity kind resolves the
-- owner by the :id path param (the affected task) via the task owner chain.
INSERT INTO rbac_route_policies (name, kind, params) VALUES
    ('task.dependency.list',   'entity', '{"resource":"task","action":"view","owner":"id"}'),
    ('task.dependency.create', 'entity', '{"resource":"task","action":"update","owner":"id"}'),
    ('task.dependency.update', 'entity', '{"resource":"task","action":"update","owner":"id"}'),
    ('task.dependency.delete', 'entity', '{"resource":"task","action":"update","owner":"id"}')
ON CONFLICT (name) DO NOTHING;