

-- name: CanUserCreateProject :one
SELECT EXISTS (
    SELECT 1 FROM users 
    WHERE id = @user_id::bigint 
      AND role IN ('admin', 'dp')
      AND deleted_at IS NULL
) AS can_create;

-- name: CreateProject :one
INSERT INTO projects (code, start_date, end_date, priority, owner_id)
VALUES (
  @code::text,
  @start_date::date,
  @end_date::date,
  @priority::bigint,
  @owner_id
)
RETURNING *;

-- TODO: write CanUserViewProject

-- name: ListProjects :many
SELECT *
FROM projects
WHERE deleted_at IS NULL
ORDER BY id ASC;

-- name: FindProject :one
SELECT *
FROM projects
WHERE deleted_at IS NULL
  AND id = @project_id::bigint;

-- name: CanUserUpdateProject :one
SELECT EXISTS (
    SELECT 1 FROM projects p
    WHERE p.id = @project_id::bigint
      AND p.deleted_at IS NULL
      AND EXISTS (
          SELECT 1 FROM users u
          WHERE u.id = @user_id::bigint
            AND u.role IN ('admin', 'dp')
            AND u.deleted_at IS NULL
      )
) AS can_manage;

-- name: UpdateProject :one
UPDATE projects
SET
	code = @code,
	priority = @priority::bigint,
	start_date = @start_date,
	end_date = @end_date,
  owner_id = @owner_id,
	updated_at = NOW()
WHERE deleted_at IS NULL
	AND id = @project_id::bigint
RETURNING *;

-- name: CanUserDeleteProject :one
SELECT EXISTS (
    SELECT 1 FROM projects p
    WHERE p.id = @project_id::bigint
      AND p.deleted_at IS NULL
      AND EXISTS (
          SELECT 1 FROM users u
          WHERE u.id = @user_id::bigint
            AND u.role IN ('admin', 'dp')
            AND u.deleted_at IS NULL
      )
) AS can_manage;


-- name: DeleteProject :exec
UPDATE projects
SET deleted_at = NOW(), updated_at = NOW()
WHERE deleted_at IS NULL
	AND id = @project_id::bigint;