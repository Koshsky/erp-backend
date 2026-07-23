-- name: CanUserManageTask :one
SELECT EXISTS (
    SELECT 1 FROM tasks t
    JOIN processes p ON t.process_id = p.id
    WHERE t.id = @task_id::bigint
	  AND t.deleted_at is NULL
      AND p.owner_id = @user_id::bigint
) AS can_manage;

-- name: CanUserCreateTask :one
SELECT EXISTS (
    SELECT 1 FROM processes p
    WHERE p.id = @process_id::bigint
	  AND p.deleted_at is NULL
      AND p.owner_id = @user_id::bigint
) AS can_create;

-- name: GetTask :one
SELECT *
FROM tasks
WHERE deleted_at IS NULL
	AND id = @resource_id::bigint;

-- name: CreateTask :one
INSERT INTO tasks (process_id, title, start_date, end_date)
VALUES (@process_id, @title, @start_date, @end_date)
RETURNING *;

-- name: UpdateTask :one
UPDATE tasks
SET
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

-- name: GetDetailedTask :one
SELECT 
    t.id,
    t.process_id,
    t.title,
    t.start_date,
    t.end_date,
    COALESCE(
        (SELECT json_agg(
            json_build_object(
                'id', a.id,
                'task_id', a.task_id,
                'resource_id', a.resource_id,
                'quantity', a.quantity
            )
            ORDER BY a.id
        ) FROM assignments a WHERE a.task_id = t.id),
        '[]'::json
    )::jsonb AS assignments
FROM tasks t
WHERE t.id = $1;