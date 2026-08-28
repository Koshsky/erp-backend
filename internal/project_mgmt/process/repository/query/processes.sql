-- name: CreateProcess :one
INSERT INTO processes (project_id, title, start_date, end_date, owner_id, sort_order)
VALUES (
	@project_id::bigint,
	@title::text,
	@start_date::date,
	@end_date::date,
	@owner_id,
	-- New process goes to the end of its project group.
	(SELECT COALESCE(MAX(sort_order), 0) + 1 FROM processes WHERE project_id = @project_id::bigint)
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
    (@scope_view::text = 'ancestor' AND (p.owner_id = @user_id::bigint OR pr.owner_id = @user_id::bigint)) OR
    (@scope_view::text = 'own' AND p.owner_id = @user_id::bigint)
  )
  AND (@owner_id::bigint = 0 OR p.owner_id = @owner_id::bigint OR pr.owner_id = @owner_id::bigint)
ORDER BY p.sort_order ASC, p.id ASC
LIMIT @page_limit::bigint OFFSET @page_offset::bigint;

-- name: CountProcesses :one
SELECT COUNT(*)
FROM processes p
JOIN projects pr ON pr.id = p.project_id
WHERE p.deleted_at IS NULL
  AND (
    @scope_view::text = 'all' OR
    (@scope_view::text = 'parent' AND pr.owner_id = @user_id::bigint) OR
    (@scope_view::text = 'ancestor' AND (p.owner_id = @user_id::bigint OR pr.owner_id = @user_id::bigint)) OR
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

-- name: ReorderProcesses :exec
-- Rewrites the order of the given process ids in one statement (transactional
-- guarantee comes from a single UPDATE: the caller sends the whole group).
UPDATE processes p
SET sort_order = x.ord, updated_at = NOW()
FROM unnest(@ids::bigint[]) WITH ORDINALITY AS x(id, ord)
WHERE p.id = x.id AND p.deleted_at IS NULL;

-- name: ListProcessIdsByProject :many
-- Active process ids of a project — to validate a reorder request covers the
-- whole group.
SELECT id
FROM processes
WHERE project_id = @project_id::bigint
	AND deleted_at IS NULL
ORDER BY sort_order ASC, id ASC;
