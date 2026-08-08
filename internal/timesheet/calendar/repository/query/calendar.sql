-- name: ListResources :many
SELECT id, title, code, owner_id
FROM resources
WHERE deleted_at IS NULL
ORDER BY id ASC;

-- Employees that could be active within the [start_date, end_date] window
-- (by hire_date/termination_date), without per-day expansion.
-- name: ListEmployeesForCalendar :many
SELECT e.id, e.resource_id, e.hire_date, e.termination_date
FROM employees e
WHERE e.deleted_at IS NULL
    AND (e.hire_date IS NULL OR e.hire_date <= @end_date::date)
    AND (e.termination_date IS NULL OR e.termination_date >= @start_date::date)
ORDER BY e.resource_id ASC, e.id ASC;

-- Absence intervals (is_available = false) overlapping the window, without expansion.
-- name: ListUnavailableRanges :many
SELECT em.resource_id, es.start_date, es.end_date
FROM employee_states es
JOIN employees em ON em.id = es.employee_id
JOIN states s ON s.id = es.state_id
WHERE s.is_available = FALSE
    AND em.deleted_at IS NULL
    AND es.end_date >= @start_date::date
    AND es.start_date <= @end_date::date
ORDER BY em.resource_id ASC, es.start_date ASC;
