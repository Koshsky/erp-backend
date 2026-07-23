-- name: CanUserManageProcess :one
SELECT EXISTS (
    SELECT 1 FROM processes p
    JOIN projects pr ON pr.id = p.project_id
    WHERE p.id = @process_id::bigint
      AND p.deleted_at IS NULL
      AND pr.owner_id = @user_id::bigint
) AS can_manage;

-- name: CanUserCreateProcess :one
SELECT EXISTS(
	SELECT 1 FROM projects
	WHERE id = @project_id::bigint
		AND deleted_at IS NULL
		AND owner_id = @user_id::bigint
) AS can_create;

-- name: ListProcesses :many
SELECT *
FROM processes p
WHERE p.deleted_at IS NULL
  AND (
    (@role::text = 'ДП') OR
    (p.owner_id = @user_id::bigint) OR
    (p.project_id IN (
        SELECT id FROM projects pr 
        WHERE pr.owner_id = @user_id::bigint
    ))
  )
ORDER BY 
    (SELECT priority FROM projects pr WHERE pr.id = p.project_id) ASC,
    p.start_date ASC, 
    p.end_date ASC, 
    p.id ASC;

-- name: GetProcess :one
SELECT *
FROM processes
WHERE id = @id::bigint
	AND deleted_at IS NULL;

-- name: CreateProcess :one
INSERT INTO processes (project_id, title, start_date, end_date, owner_id)
VALUES (
	@project_id::bigint,
	@title::text,
	@start_date::date,
	@end_date::date,
	@owner_id::bigint
)
RETURNING *;

-- name: UpdateProcess :one
UPDATE processes
SET
	title = @title,
	start_date = @start_date,
	end_date = @end_date,
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