-- name: CreateTask :one
INSERT INTO tasks (process_id, owner_id, title, start_date, end_date, sort_order)
VALUES (
	@process_id,
	@owner_id,
	@title,
	@start_date,
	@end_date,
	-- New task goes to the end of its process group.
	(SELECT COALESCE(MAX(sort_order), 0) + 1 FROM tasks WHERE process_id = @process_id)
)
RETURNING *;

-- name: ListTasks :many
SELECT t.*
FROM tasks t
JOIN processes p ON p.id = t.process_id
JOIN projects pr ON pr.id = p.project_id
WHERE t.deleted_at IS NULL
  AND (
    @scope_view::text = 'all' OR
    (@scope_view::text = 'parent' AND p.owner_id = @user_id::bigint) OR
    (@scope_view::text = 'ancestor' AND (t.owner_id = @user_id::bigint OR p.owner_id = @user_id::bigint OR pr.owner_id = @user_id::bigint)) OR
    (@scope_view::text = 'own' AND t.owner_id = @user_id::bigint)
  )
  AND (@owner_id::bigint = 0 OR t.owner_id = @owner_id::bigint OR p.owner_id = @owner_id::bigint OR pr.owner_id = @owner_id::bigint)
ORDER BY t.sort_order ASC, t.id ASC
LIMIT @page_limit::bigint OFFSET @page_offset::bigint;

-- name: CountTasks :one
SELECT COUNT(*)
FROM tasks t
JOIN processes p ON p.id = t.process_id
JOIN projects pr ON pr.id = p.project_id
WHERE t.deleted_at IS NULL
  AND (
    @scope_view::text = 'all' OR
    (@scope_view::text = 'parent' AND p.owner_id = @user_id::bigint) OR
    (@scope_view::text = 'ancestor' AND (t.owner_id = @user_id::bigint OR p.owner_id = @user_id::bigint OR pr.owner_id = @user_id::bigint)) OR
    (@scope_view::text = 'own' AND t.owner_id = @user_id::bigint)
  )
  AND (@owner_id::bigint = 0 OR t.owner_id = @owner_id::bigint OR p.owner_id = @owner_id::bigint OR pr.owner_id = @owner_id::bigint);

-- name: FindTask :one
SELECT *
FROM tasks
WHERE deleted_at IS NULL
	AND id = @resource_id::bigint;

-- name: UpdateTask :one
UPDATE tasks
SET
	process_id = @process_id,
	owner_id = @owner_id,
	title = @title,
	start_date = @start_date,
	end_date = @end_date,
	updated_at = NOW()
WHERE id = @task_id
	AND deleted_at IS NULL
RETURNING *;

-- name: DeleteTask :exec
UPDATE tasks
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @task_id
	AND deleted_at IS NULL;

-- name: OwnerChain :one
SELECT COALESCE(pr.owner_id, 0)::bigint AS project_owner,
       COALESCE(p.owner_id, 0)::bigint  AS process_owner,
       COALESCE(t.owner_id, 0)::bigint  AS owner_id
FROM tasks t
JOIN processes p ON p.id = t.process_id
JOIN projects pr ON pr.id = p.project_id
WHERE t.id = @id::bigint
	AND t.deleted_at IS NULL
	AND p.deleted_at IS NULL
	AND pr.deleted_at IS NULL;

-- name: ReorderTasks :exec
-- Rewrites the order of the given task ids in one statement (the caller sends
-- the whole group).
UPDATE tasks t
SET sort_order = x.ord, updated_at = NOW()
FROM unnest(@ids::bigint[]) WITH ORDINALITY AS x(id, ord)
WHERE t.id = x.id AND t.deleted_at IS NULL;

-- name: ListTaskIdsByProcess :many
-- Active task ids of a process — to validate a reorder request covers the
-- whole group.
SELECT id
FROM tasks
WHERE process_id = @process_id::bigint
	AND deleted_at IS NULL
ORDER BY sort_order ASC, id ASC;
