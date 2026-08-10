-- name: ListEmployeesByResourceID :many
SELECT *
FROM employees
WHERE resource_id = @resource_id::bigint
	AND deleted_at IS NULL
ORDER BY id ASC;

-- name: ListEmployees :many
SELECT *
FROM employees
WHERE deleted_at IS NULL
  AND (@role::text = 'admin' OR manager_id = @user_id::bigint)
  AND (@manager_id::bigint = 0 OR manager_id = @manager_id::bigint)
ORDER BY id ASC
LIMIT @page_limit::bigint OFFSET @page_offset::bigint;

-- name: CountEmployees :one
SELECT COUNT(*)
FROM employees
WHERE deleted_at IS NULL
  AND (@role::text = 'admin' OR manager_id = @user_id::bigint)
  AND (@manager_id::bigint = 0 OR manager_id = @manager_id::bigint);

-- name: FindEmployee :one
SELECT *
FROM employees
WHERE id = @employee_id::bigint
	AND deleted_at IS NULL;

-- name: CreateEmployee :one
INSERT INTO employees (resource_id, name, position, manager_id, hire_date, termination_date)
VALUES (@resource_id, @name, @position, @manager_id, @hire_date, @termination_date)
RETURNING *;

-- name: UpdateEmployee :one
UPDATE employees
SET
	resource_id = @resource_id,
	name = @name,
	position = @position,
	manager_id = @manager_id,
	hire_date = @hire_date,
	termination_date = @termination_date,
	updated_at = NOW()
WHERE id = @employee_id
	AND deleted_at IS NULL
RETURNING *;

-- name: DeleteEmployee :exec
UPDATE employees
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = @employee_id
	AND deleted_at IS NULL;

-- name: IsResourceActive :one
SELECT EXISTS(
	SELECT 1
	FROM resources
	WHERE id = @resource_id::bigint
		AND deleted_at IS NULL
);

-- name: OwnerChain :one
SELECT COALESCE(manager_id, 0)::bigint AS owner_id
FROM employees
WHERE id = @id::bigint
	AND deleted_at IS NULL;
