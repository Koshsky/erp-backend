-- =============================================
-- TASK SUBTASKS (OPERATIONS) AND EXECUTION STATUS
-- =============================================
-- The columns are also declared in V1 (the sqlc schema source), so a fresh
-- database already has them and the ALTERs below are no-ops there; existing
-- deployments get the columns from this migration.
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS parent_id BIGINT REFERENCES tasks(id) ON DELETE CASCADE;
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS status TEXT NOT NULL DEFAULT 'not_started';

-- Status catalog CHECK. A fresh database already has the inline constraint
-- from V1 (auto-named tasks_status_check), so guard the ALTER with a DO block
-- to stay idempotent for both fresh and existing deployments.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'tasks_status_check'
    ) THEN
        ALTER TABLE tasks ADD CONSTRAINT tasks_status_check
        CHECK (status IN ('not_started', 'in_progress', 'done'));
    END IF;
END $$;

-- =============================================
-- ORDER UNIQUENESS PER PARENT GROUP
-- =============================================
-- Replaces the V2 idx_tasks_process_sort_order (process_id, sort_order):
-- top-level tasks (parent_id IS NULL → COALESCE 0) must stay unique within
-- their process, subtasks within their parent task. Since subtask process_id
-- always equals the parent's (enforced in the service), (process_id, parent
-- group, sort_order) covers both levels in a single unique index. Deleted
-- rows are physically removed (archive triggers, V5), so no partial predicate
-- is needed.
DROP INDEX IF EXISTS idx_tasks_process_sort_order;
CREATE UNIQUE INDEX idx_tasks_parent_sort_order
ON tasks(process_id, COALESCE(parent_id, 0), sort_order);

-- FK lookup for subtasks.
CREATE INDEX IF NOT EXISTS idx_tasks_parent_id
ON tasks(parent_id) WHERE parent_id IS NOT NULL;