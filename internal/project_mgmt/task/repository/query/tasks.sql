-- name: CreateTask :one
INSERT INTO tasks (process_id, owner_id, title, start_date, end_date)
VALUES (@process_id, @owner_id, @title, @start_date, @end_date)
RETURNING *;

-- name: ListTasks :many
SELECT t.*
FROM tasks t
JOIN processes p ON p.id = t.process_id
JOIN projects pr ON pr.id = p.project_id
WHERE t.deleted_at IS NULL
  AND (@see_all::boolean OR t.owner_id = @user_id::bigint OR p.owner_id = @user_id::bigint OR pr.owner_id = @user_id::bigint)
  AND (@owner_id::bigint = 0 OR t.owner_id = @owner_id::bigint OR p.owner_id = @owner_id::bigint OR pr.owner_id = @owner_id::bigint)
ORDER BY t.id ASC
LIMIT @page_limit::bigint OFFSET @page_offset::bigint;

-- name: CountTasks :one
SELECT COUNT(*)
FROM tasks t
JOIN processes p ON p.id = t.process_id
JOIN projects pr ON pr.id = p.project_id
WHERE t.deleted_at IS NULL
  AND (@see_all::boolean OR t.owner_id = @user_id::bigint OR p.owner_id = @user_id::bigint OR pr.owner_id = @user_id::bigint)
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
       COALESCE(p.owner_id, 0)::bigint  AS process_owner
FROM tasks t
JOIN processes p ON p.id = t.process_id
JOIN projects pr ON pr.id = p.project_id
WHERE t.id = @id::bigint
	AND t.deleted_at IS NULL
	AND p.deleted_at IS NULL
	AND pr.deleted_at IS NULL;
