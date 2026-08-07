-- name: ListEmployeesByResourceID :many
SELECT e.*, r.title AS resource_title
FROM employees e
JOIN resources r ON r.id = e.resource_id
WHERE e.resource_id = @resource_id::bigint
	AND e.deleted_at IS NULL
ORDER BY e.id ASC;

-- name: ListEmployees :many
SELECT e.*, r.title AS resource_title
FROM employees e
JOIN resources r ON r.id = e.resource_id
WHERE e.deleted_at IS NULL
ORDER BY e.id ASC;

-- name: ListEmployeesByManagerID :many
SELECT e.*, r.title AS resource_title
FROM employees e
JOIN resources r ON r.id = e.resource_id
WHERE e.deleted_at IS NULL
AND e.manager_id = @manager_id::bigint
ORDER BY e.id ASC;

-- name: FindEmployee :one
SELECT e.*, r.title AS resource_title
FROM employees e
JOIN resources r ON r.id = e.resource_id
WHERE e.id = @employee_id::bigint
	AND e.deleted_at IS NULL;

-- name: CreateEmployee :one
INSERT INTO employees (resource_id, name, manager_id, hire_date, termination_date)
VALUES (@resource_id, @name, @manager_id, @hire_date, @termination_date)
RETURNING *;

-- name: UpdateEmployee :one
UPDATE employees
SET
	resource_id = @resource_id,
	name = @name,
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
