-- name: CanUserCreateProcess :one
SELECT EXISTS(
	SELECT 1 FROM projects
	WHERE id = @project_id::bigint
		AND deleted_at IS NULL
		AND owner_id = @user_id::bigint
) AS can_create;

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

-- TODO: write CanUserViewProcess

-- name: ListProcesss :many
SELECT *
FROM processes
WHERE deleted_at IS NULL
ORDER BY id ASC;

-- name: FindProcess :one
SELECT *
FROM processes
WHERE id = @id::bigint
	AND deleted_at IS NULL;
    

-- name: CanUserUpdateProcess :one
SELECT EXISTS (
    SELECT 1 FROM processes p
    JOIN projects pr ON pr.id = p.project_id
    WHERE p.id = @process_id::bigint
      AND p.deleted_at IS NULL
      AND pr.owner_id = @user_id::bigint
) AS can_manage;

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

-- name: CanUserDeleteProcess :one
SELECT EXISTS (
    SELECT 1 FROM processes p
    JOIN projects pr ON pr.id = p.project_id
    WHERE p.id = @process_id::bigint
      AND p.deleted_at IS NULL
      AND pr.owner_id = @user_id::bigint
) AS can_manage;

-- name: DeleteProcess :exec
UPDATE processes
SET deleted_at = NOW(), updated_at = NOW()
WHERE deleted_at IS NULL
	AND id = @process_id::bigint;