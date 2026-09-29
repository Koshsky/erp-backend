-- =============================================
-- TASK COMMENTS — OBSOLETE AS OF V5 REWRITE
-- =============================================
-- task_comments archiving is covered by the generic archive trigger
-- (fn_archive_row in V5): deleting a task cascades to its comments (FK
-- task_comments.task_id → tasks ON DELETE CASCADE), and each comment row is
-- archived via trg_archive_task_comments. The former cascade-soft-delete and
-- hard-DELETE-block triggers no longer exist.
-- This migration is kept as a version marker so existing deployments that
-- already applied the old triggers drop them here.
DROP TRIGGER IF EXISTS trigger_cascade_soft_delete_task_comments ON tasks;
DROP TRIGGER IF EXISTS block_hard_delete_on_task_comments ON task_comments;