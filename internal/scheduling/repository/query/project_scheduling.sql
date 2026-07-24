-- name: GetProjectScheduling :many
SELECT * FROM projects
WHERE deleted_at IS NULL
AND @role = 'ДП'
ORDER BY priority ASC;
