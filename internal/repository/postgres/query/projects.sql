-- name: ListProjects :many
SELECT *
FROM projects
WHERE deleted_at IS NULL
  AND (
    @role::text = 'ДП' OR
    owner_id = @user_id::bigint
  )
ORDER BY priority ASC, start_date ASC, id ASC;

-- name: GetProject :one
SELECT *
FROM projects
WHERE deleted_at IS NULL
  AND id = @project_id::bigint;

-- name: CanUserManageProject :one
SELECT EXISTS (
    SELECT 1 FROM projects p
    WHERE p.id = @project_id::bigint
      AND p.deleted_at IS NULL
      AND EXISTS (
          SELECT 1 FROM users u
          WHERE u.id = @user_id::bigint
            AND u.role = 'ДП'
            AND u.deleted_at IS NULL
      )
) AS can_manage;

-- name: CanUserCreateProject :one
SELECT EXISTS (
    SELECT 1 FROM users 
    WHERE id = @user_id::bigint 
      AND role = 'ДП'
      AND deleted_at IS NULL
) AS can_create;

-- name: CreateProject :one
INSERT INTO projects (code, start_date, end_date, priority, owner_id)
VALUES (
  @code::text,
  @start_date::date,
  @end_date::date,
  @priority,
  @owner_id::bigint
)
RETURNING *;

-- name: UpdateProject :one
UPDATE projects
SET
	code = @code,
	priority = @priority,
	start_date = @start_date,
	end_date = @end_date,
  owner_id = @owner_id,
	updated_at = NOW()
WHERE deleted_at IS NULL
	AND id = @project_id::bigint
RETURNING *;

-- name: DeleteProject :exec
UPDATE projects
SET deleted_at = NOW(), updated_at = NOW()
WHERE deleted_at IS NULL
	AND id = @project_id::bigint;