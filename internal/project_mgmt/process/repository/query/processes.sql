-- name: CreateProcess :one
INSERT INTO processes (project_id, title, start_date, end_date, owner_id)
VALUES (
	@project_id::bigint,
	@title::text,
	@start_date::date,
	@end_date::date,
	@owner_id
)
RETURNING *;

-- name: ListProcesss :many
SELECT p.*
FROM processes p
JOIN projects pr ON pr.id = p.project_id
WHERE p.deleted_at IS NULL
  AND (
    @scope_view::text = 'all' OR
    (@scope_view::text = 'parent' AND pr.owner_id = @user_id::bigint) OR
    (@scope_view::text = 'own' AND p.owner_id = @user_id::bigint)
  )
  AND (@owner_id::bigint = 0 OR p.owner_id = @owner_id::bigint OR pr.owner_id = @owner_id::bigint)
ORDER BY p.id ASC
LIMIT @page_limit::bigint OFFSET @page_offset::bigint;

-- name: CountProcesses :one
SELECT COUNT(*)
FROM processes p
JOIN projects pr ON pr.id = p.project_id
WHERE p.deleted_at IS NULL
  AND (
    @scope_view::text = 'all' OR
    (@scope_view::text = 'parent' AND pr.owner_id = @user_id::bigint) OR
    (@scope_view::text = 'own' AND p.owner_id = @user_id::bigint)
  )
  AND (@owner_id::bigint = 0 OR p.owner_id = @owner_id::bigint OR pr.owner_id = @owner_id::bigint);

-- name: FindProcess :one
SELECT *
FROM processes
WHERE id = @id::bigint
	AND deleted_at IS NULL;
    

-- name: UpdateProcess :one
UPDATE processes
SET
	title = @title,
	start_date = @start_date,
	end_date = @end_date,
	project_id = COALESCE(@project_id, project_id),
	owner_id = COALESCE(@owner_id, owner_id),
	updated_at = NOW()
WHERE deleted_at IS NULL
	AND id = @process_id::bigint
RETURNING *;

-- name: DeleteProcess :exec
UPDATE processes
SET deleted_at = NOW(), updated_at = NOW()
WHERE deleted_at IS NULL
	AND id = @process_id::bigint;

-- name: OwnerChain :one
SELECT COALESCE(pr.owner_id, 0)::bigint AS project_owner,
       COALESCE(p.owner_id, 0)::bigint  AS process_owner
FROM processes p
JOIN projects pr ON pr.id = p.project_id
WHERE p.id = @id::bigint
	AND p.deleted_at IS NULL
	AND pr.deleted_at IS NULL;
