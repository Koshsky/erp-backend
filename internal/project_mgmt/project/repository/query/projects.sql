-- name: CreateProject :one
-- Idempotent create by business key code: if an active code already exists
-- we insert nothing; the calling code (repository) turns the conflict into 409.
INSERT INTO projects (code, start_date, end_date, priority, owner_id)
VALUES (
  @code::text,
  @start_date::date,
  @end_date::date,
  @priority::bigint,
  @owner_id
)
ON CONFLICT (code) WHERE deleted_at IS NULL
DO NOTHING
RETURNING *;

-- name: ListProjects :many
SELECT *
FROM projects
WHERE deleted_at IS NULL
  AND (
    @scope_view::text = 'all' OR
    (@scope_view::text = 'own' AND owner_id = @user_id::bigint)
  )
  AND (@owner_id::bigint = 0 OR owner_id = @owner_id::bigint)
ORDER BY id ASC
LIMIT @page_limit::bigint OFFSET @page_offset::bigint;

-- name: CountProjects :one
SELECT COUNT(*)
FROM projects
WHERE deleted_at IS NULL
  AND (
    @scope_view::text = 'all' OR
    (@scope_view::text = 'own' AND owner_id = @user_id::bigint)
  )
  AND (@owner_id::bigint = 0 OR owner_id = @owner_id::bigint);

-- name: FindProject :one
SELECT *
FROM projects
WHERE deleted_at IS NULL
  AND id = @project_id::bigint;

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

-- name: DeleteProject :exec
UPDATE projects
SET deleted_at = NOW(), updated_at = NOW()
WHERE deleted_at IS NULL
	AND id = @project_id::bigint;

-- name: OwnerChain :one
SELECT COALESCE(owner_id, 0)::bigint AS owner_id
FROM projects
WHERE id = @id::bigint
	AND deleted_at IS NULL;
