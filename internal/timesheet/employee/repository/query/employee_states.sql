-- name: ListStatesByEmployeeRange :many
SELECT es.id, es.employee_id, es.start_date, es.end_date, es.state_id,
	s.code AS state_code, s.name AS state_name, s.is_available
FROM employee_states es
JOIN states s ON s.id = es.state_id
WHERE es.employee_id = @employee_id::bigint
	AND es.end_date >= @start_date::date
	AND es.start_date <= @end_date::date
ORDER BY es.start_date ASC;

-- name: ListOverlappingStates :many
SELECT es.id, es.employee_id, es.start_date, es.end_date, es.state_id
FROM employee_states es
WHERE es.employee_id = @employee_id::bigint
	AND es.end_date >= @start_date::date
	AND es.start_date <= @end_date::date
ORDER BY es.start_date ASC
FOR UPDATE;

-- name: ListOverlappingStatesByState :many
SELECT es.id, es.employee_id, es.start_date, es.end_date, es.state_id
FROM employee_states es
WHERE es.employee_id = @employee_id::bigint
	AND es.state_id = @state_id::bigint
	AND es.end_date >= @start_date::date
	AND es.start_date <= @end_date::date
ORDER BY es.start_date ASC
FOR UPDATE;

-- name: DeleteOverlapping :exec
DELETE FROM employee_states
WHERE employee_id = @employee_id::bigint
	AND end_date >= @start_date::date
	AND start_date <= @end_date::date;

-- name: DeleteOverlappingByState :exec
DELETE FROM employee_states
WHERE employee_id = @employee_id::bigint
	AND state_id = @state_id::bigint
	AND end_date >= @start_date::date
	AND start_date <= @end_date::date;

-- name: InsertStateRange :one
INSERT INTO employee_states (employee_id, state_id, start_date, end_date)
VALUES (@employee_id, @state_id, @start_date::date, @end_date::date)
RETURNING *;
