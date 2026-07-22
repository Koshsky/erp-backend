-- name: GetProcessWithProject :one
SELECT
    p.id AS process_id,
    p.title AS process_title,
    p.start_date AS process_start_date,
    p.end_date AS process_end_date,
    pr.id AS project_id,
    pr.code AS project_code,
    pr.start_date AS project_start_date,
    pr.end_date AS project_end_date
FROM processes p
JOIN projects pr ON pr.id = p.project_id
WHERE p.id = $1
    AND p.deleted_at IS NULL
    AND pr.deleted_at IS NULL;


-- name: GetProcessTasksWithResources :many
SELECT
    t.id AS task_id,
    t.title AS task_title,
    t.start_date AS task_start_date,
    t.end_date AS task_end_date,
    a.id AS assignment_id,
    a.quantity AS assignment_quantity,
    r.id AS resource_id,
    r.title AS resource_title,
    r.code AS resource_code,
    r.quantity AS resource_quantity
FROM tasks t
LEFT JOIN assignments a ON a.task_id = t.id AND a.deleted_at IS NULL
LEFT JOIN resources r ON r.id = a.resource_id AND r.deleted_at IS NULL
WHERE t.process_id = $1
    AND t.deleted_at IS NULL
ORDER BY t.start_date ASC, t.end_date ASC, t.id ASC, a.id ASC;

