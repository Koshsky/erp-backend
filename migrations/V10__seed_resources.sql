INSERT INTO users (username, name, password_hash, role) VALUES
('admin', 'Admin Name', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'ДП'),  -- директор проектов, админ
('user1', 'User One', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'РП'),
('user2', 'User Two', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'ВП');

INSERT INTO resources (title, code, quantity) VALUES
('Инженер', 'И', 7),
('Монтажник', 'М', 4),
('Производитель работ', 'ПР', 2),
('Руководитель группы', 'РГ', 4),
('Руководитель службы инсталляции', 'РСИ', 1);

INSERT INTO projects (code, start_date, end_date, priority, owner_id)
VALUES
    ('KO-1001', DATE '2026-07-15', DATE '2026-08-30', 1, 1),
    ('KO-1002', DATE '2026-08-01', DATE '2026-09-20', 2, 1),
    ('KO-1003', DATE '2026-09-05', DATE '2026-10-25', 3, 1);

INSERT INTO assignments (task_id, resource_id, quantity)
SELECT t.id, r.id, x.qty
FROM (
    VALUES
    ('KO-1001', 'Осмотр объекта', 'Инженер', 2),
    ('KO-1001', 'Монтаж кабеленесущих систем и кабельных трасс', 'Монтажник', 4),
    ('KO-1001', 'Инсталляция оконечного оборудования телемедицины', 'Инженер', 3),
    ('KO-1001', 'Пуско-наладочные работы', 'Производитель работ', 1),
    ('KO-1001', 'Проведение инструктажа и передача инструкций мед персоналу', 'Руководитель группы', 1),
    ('KO-1001', 'Подписание Акта ввода в эксплуатацию', 'Руководитель службы инсталляции', 1),
    ('KO-1002', 'Осмотр объекта', 'Инженер', 2),
    ('KO-1002', 'Инсталляция оконечного оборудования телемедицины', 'Монтажник', 3)
) AS x(project_code, task_title, resource_title, qty)
JOIN tasks t ON t.title = x.task_title
JOIN projects p ON p.code = x.project_code AND t.process_id IN (
    SELECT id FROM processes WHERE project_id = p.id
)
JOIN resources r ON r.title = x.resource_title;


INSERT INTO milestones (process_id, title, content, date)
SELECT pr.id, m.title, m.content, m.milestone_date
FROM processes pr
JOIN projects p ON p.id = pr.project_id
JOIN (
    VALUES
        ('KO-1001', 'Старт проекта', 'Начало работ по объекту KO-1001', DATE '2026-07-15'),
        ('KO-1001', 'Завершение ПНР', 'Пуско-наладочные работы завершены', DATE '2026-08-12'),
        ('KO-1001', 'Ввод в эксплуатацию', 'Подписан акт ввода в эксплуатацию', DATE '2026-08-22'),
        ('KO-1002', 'Старт проекта', 'Начало работ по объекту KO-1002', DATE '2026-08-01')
) AS m(code, title, content, milestone_date)
    ON m.code = p.code
WHERE pr.title = 'Инсталляция';