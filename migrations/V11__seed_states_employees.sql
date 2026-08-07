-- Справочник состояний: доступность задаётся is_available.
INSERT INTO states (code, name, is_available) VALUES
('present', 'Явка', TRUE),
('business_trip', 'Командировка', TRUE),
('time_off', 'Отгул', FALSE),
('vacation', 'Отпуск', FALSE),
('sick_leave', 'Больничный', FALSE);

-- Конкретные сотрудники (уникальные ресурсы) под категории из V10.
INSERT INTO employees (resource_id, name, hire_date, termination_date)
SELECT r.id, e.name, e.hire_date, e.termination_date
FROM (
    VALUES
        ('Инженер', 'Иванов Иван Иванович', DATE '2024-01-15', NULL),
        ('Инженер', 'Петров Пётр Петрович', DATE '2024-03-01', NULL),
        ('Инженер', 'Сидоров Сидор Сидорович', DATE '2025-05-10', NULL),
        ('Инженер', 'Кузнецов Кузьма Кузьмич', DATE '2025-08-20', NULL),
        ('Инженер', 'Смирнов Смирн Смирнович', DATE '2026-01-12', NULL),
        ('Инженер', 'Волков Волк Волкович', DATE '2026-02-01', NULL),
        ('Инженер', 'Егоров Егор Егорович', DATE '2026-06-05', NULL),
        ('Монтажник', 'Фёдоров Фёдор Фёдорович', DATE '2024-02-10', NULL),
        ('Монтажник', 'Макаров Макар Макарович', DATE '2025-04-15', NULL),
        ('Монтажник', 'Андреев Андрей Андреевич', DATE '2026-03-20', NULL),
        ('Монтажник', 'Борисов Борис Борисович', DATE '2024-11-01', DATE '2026-06-30'),
        ('Производитель работ', 'Григорьев Григорий Григорьевич', DATE '2024-05-05', NULL),
        ('Производитель работ', 'Николаев Николай Николаевич', DATE '2025-09-01', NULL),
        ('Руководитель группы', 'Дмитриев Дмитрий Дмитриевич', DATE '2023-10-01', NULL),
        ('Руководитель группы', 'Степанов Степан Степанович', DATE '2024-06-15', NULL),
        ('Руководитель группы', 'Ильин Илья Ильич', DATE '2025-01-20', NULL),
        ('Руководитель группы', 'Захаров Захар Захарович', DATE '2026-04-01', NULL),
        ('Руководитель службы инсталляции', 'Семёнов Семён Семёнович', DATE '2022-09-01', NULL)
) AS e(resource_title, name, hire_date, termination_date)
JOIN resources r ON r.title = e.resource_title;

-- Назначаем руководителей (аккаунты пользователей) по категориям.
UPDATE employees e
SET manager_id = u.id
FROM (
    VALUES
        ('Инженер', 'vp1'),
        ('Монтажник', 'vp1'),
        ('Производитель работ', 'vp1'),
        ('Руководитель группы', 'vp1'),
        ('Руководитель службы инсталляции', 'vp1')
) AS m(resource_title, manager_username)
JOIN resources r ON r.title = m.resource_title
JOIN users u ON u.username = m.manager_username
WHERE e.resource_id = r.id;

-- Демо-состояния (диапазоны) на период проекта КО-01 (2026-07-15..08-30):
-- 2 инженера в отпуске, 1 на больничном, 1 монтажник в командировке (доступен).
INSERT INTO employee_states (employee_id, state_id, start_date, end_date)
SELECT e.id, s.id, x.start_date, x.end_date
FROM (
    VALUES
        ('Иванов Иван Иванович', 'vacation', DATE '2026-07-20', DATE '2026-08-02'),
        ('Петров Пётр Петрович', 'vacation', DATE '2026-08-10', DATE '2026-08-24'),
        ('Сидоров Сидор Сидорович', 'sick_leave', DATE '2026-07-27', DATE '2026-08-05'),
        ('Фёдоров Фёдор Фёдорович', 'business_trip', DATE '2026-08-01', DATE '2026-08-07')
) AS x(employee_name, state_code, start_date, end_date)
JOIN employees e ON e.name = x.employee_name
JOIN states s ON s.code = x.state_code;
