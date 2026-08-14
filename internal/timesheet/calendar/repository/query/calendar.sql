-- name: ListResources :many
SELECT id, title, code, owner_id
FROM resources
WHERE deleted_at IS NULL
ORDER BY id ASC;

-- Members that could be active within the [start_date, end_date] window
-- (by hire_date/termination_date), without per-day expansion.
-- name: ListEmployeesForCalendar :many
SELECT u.id, rm.resource_id, u.hire_date, u.termination_date
FROM resource_members rm
JOIN users u ON u.id = rm.user_id
WHERE u.deleted_at IS NULL
    AND (u.hire_date IS NULL OR u.hire_date <= @end_date::date)
    AND (u.termination_date IS NULL OR u.termination_date >= @start_date::date)
ORDER BY rm.resource_id ASC, u.id ASC;

-- Absence intervals (is_available = false) overlapping the window, without expansion.
-- name: ListUnavailableRanges :many
SELECT rm.resource_id, es.start_date, es.end_date
FROM user_states es
JOIN resource_members rm ON rm.user_id = es.user_id
JOIN states s ON s.id = es.state_id
WHERE s.is_available = FALSE
    AND es.end_date >= @start_date::date
    AND es.start_date <= @end_date::date
ORDER BY rm.resource_id ASC, es.start_date ASC;
