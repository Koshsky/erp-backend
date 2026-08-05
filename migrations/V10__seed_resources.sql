-- Пользователи: 1 admin, 1 dp (директор портфеля), 2 rp (руководители проектов), 4 vp (владельцы процессов), 2 worker
INSERT INTO users (username, name, password_hash, role) VALUES
('admin', 'Admin Name', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'admin'),
('dp1', 'ДП-1', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'dp'),
('rp1', 'РП-1', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'rp'),
('rp2', 'РП-2', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'rp'),
('vp1', 'ВП-1', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'vp'),
('vp2', 'ВП-2', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'vp'),
('vp3', 'ВП-3', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'vp'),
('vp4', 'ВП-4', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'vp'),
('w1', 'Работник-1', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'worker'),
('w2', 'Работник-2', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'worker');

-- Иерархия подчинённости: admin -> dp -> rp -> vp -> worker.
UPDATE users SET manager_id = (SELECT id FROM users WHERE username = 'admin') WHERE username = 'dp1';
UPDATE users SET manager_id = (SELECT id FROM users WHERE username = 'dp1') WHERE username IN ('rp1', 'rp2');
UPDATE users SET manager_id = (SELECT id FROM users WHERE username = 'rp1') WHERE username IN ('vp1', 'vp2');
UPDATE users SET manager_id = (SELECT id FROM users WHERE username = 'rp2') WHERE username IN ('vp3', 'vp4');
UPDATE users SET manager_id = (SELECT id FROM users WHERE username = 'vp1') WHERE username = 'w1';
UPDATE users SET manager_id = (SELECT id FROM users WHERE username = 'vp3') WHERE username = 'w2';

INSERT INTO resources (title, code, quantity) VALUES
('Инженер', 'И', 7),
('Монтажник', 'М', 4),
('Производитель работ', 'ПР', 2),
('Руководитель группы', 'РГ', 4),
('Руководитель службы инсталляции', 'РСИ', 1);

-- По 1 проекту на РП; код проекта отражает цепочку владения ДП-РП-ВП.
-- V5-триггер при вставке автосоздаст по 2 процесса (Инсталляция + Производство) на проект.
INSERT INTO projects (code, start_date, end_date, priority, owner_id)
VALUES
    ('КО-01_РП1_ВП1', DATE '2026-07-15', DATE '2026-08-30', 1, (SELECT id FROM users WHERE username = 'rp1')),
    ('КО-02_РП2_ВП3', DATE '2026-08-01', DATE '2026-09-20', 2, (SELECT id FROM users WHERE username = 'rp2'));

-- Редактируем автосозданные процессы: проставляем владельца-ВП (имена не меняем)
UPDATE processes SET
    owner_id = (SELECT id FROM users WHERE username = 'vp1')
WHERE project_id = (SELECT id FROM projects WHERE code = 'КО-01_РП1_ВП1')
  AND title = 'Инсталляция'
  AND deleted_at IS NULL;

UPDATE processes SET
    owner_id = (SELECT id FROM users WHERE username = 'vp2')
WHERE project_id = (SELECT id FROM projects WHERE code = 'КО-01_РП1_ВП1')
  AND title = 'Производство'
  AND deleted_at IS NULL;

UPDATE processes SET
    owner_id = (SELECT id FROM users WHERE username = 'vp3')
WHERE project_id = (SELECT id FROM projects WHERE code = 'КО-02_РП2_ВП3')
  AND title = 'Инсталляция'
  AND deleted_at IS NULL;

UPDATE processes SET
    owner_id = (SELECT id FROM users WHERE username = 'vp4')
WHERE project_id = (SELECT id FROM projects WHERE code = 'КО-02_РП2_ВП3')
  AND title = 'Производство'
  AND deleted_at IS NULL;

INSERT INTO assignments (task_id, resource_id, quantity)
SELECT t.id, r.id, x.qty
FROM (
    VALUES
    ('КО-01_РП1_ВП1', 'Осмотр объекта', 'Инженер', 2),
    ('КО-01_РП1_ВП1', 'Монтаж кабеленесущих систем и кабельных трасс', 'Монтажник', 4),
    ('КО-01_РП1_ВП1', 'Инсталляция оконечного оборудования телемедицины', 'Инженер', 3),
    ('КО-01_РП1_ВП1', 'Пуско-наладочные работы', 'Производитель работ', 1),
    ('КО-01_РП1_ВП1', 'Проведение инструктажа и передача инструкций мед персоналу', 'Руководитель группы', 1),
    ('КО-01_РП1_ВП1', 'Подписание Акта ввода в эксплуатацию', 'Руководитель службы инсталляции', 1),
    ('КО-02_РП2_ВП3', 'Осмотр объекта', 'Инженер', 2),
    ('КО-02_РП2_ВП3', 'Инсталляция оконечного оборудования телемедицины', 'Монтажник', 3)
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
        ('КО-01_РП1_ВП1', 'Старт проекта', 'Начало работ по объекту КО-01_РП1_ВП1', DATE '2026-07-15'),
        ('КО-01_РП1_ВП1', 'Завершение ПНР', 'Пуско-наладочные работы завершены', DATE '2026-08-12'),
        ('КО-01_РП1_ВП1', 'Ввод в эксплуатацию', 'Подписан акт ввода в эксплуатацию', DATE '2026-08-22'),
        ('КО-02_РП2_ВП3', 'Старт проекта', 'Начало работ по объекту КО-02_РП2_ВП3', DATE '2026-08-01')
) AS m(code, title, content, milestone_date)
    ON m.code = p.code
WHERE pr.title = 'Инсталляция';
