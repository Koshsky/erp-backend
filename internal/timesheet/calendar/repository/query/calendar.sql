-- name: ListResources :many
SELECT id, title, code
FROM resources
WHERE deleted_at IS NULL
ORDER BY id ASC;

-- По-дневная мощность категории: сколько сотрудников активно в каждый день диапазона.
-- name: ListCalendarCapacity :many
SELECT r.id AS resource_id,
    g.date::date,
    COUNT(e.id)::bigint AS capacity
FROM resources r
CROSS JOIN generate_series(@start_date::date, @end_date::date, INTERVAL '1 day') AS g(date)
LEFT JOIN employees e
    ON e.resource_id = r.id
    AND e.deleted_at IS NULL
    AND (e.hire_date IS NULL OR e.hire_date <= g.date)
    AND (e.termination_date IS NULL OR e.termination_date >= g.date)
WHERE r.deleted_at IS NULL
GROUP BY r.id, g.date::date
ORDER BY r.id ASC, g.date::date ASC;

-- По-дневое количество недоступных сотрудников (состояния с is_available = false):
-- разворачиваются только пересечения интервалов состояний с запрошенным диапазоном.
-- name: ListCalendarUnavailable :many
SELECT em.resource_id,
    g.date::date,
    COUNT(*)::bigint AS unavailable
FROM employee_states es
JOIN employees em ON em.id = es.employee_id
JOIN states s ON s.id = es.state_id
CROSS JOIN LATERAL generate_series(
    GREATEST(es.start_date, @start_date::date),
    LEAST(es.end_date, @end_date::date),
    INTERVAL '1 day'
) AS g(date)
WHERE s.is_available = FALSE
    AND em.deleted_at IS NULL
    AND es.end_date >= @start_date::date
    AND es.start_date <= @end_date::date
GROUP BY em.resource_id, g.date::date;
