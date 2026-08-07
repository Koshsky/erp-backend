-- Объединённая миграция V10: колонки для существующих БД + демо-данные.
-- На свежих БД колонки position/owner_id уже есть из V1 — ALTER-ы отрабатывают как no-op.

-- Должность сотрудника — свободный текст, не тип ресурса.
ALTER TABLE employees
	ADD COLUMN IF NOT EXISTS position TEXT NOT NULL DEFAULT '';

-- Владелец ресурса (обязателен): добавляем, подтягиваем владельца, фиксируем NOT NULL.
ALTER TABLE resources
	ADD COLUMN IF NOT EXISTS owner_id BIGINT REFERENCES users(id);
UPDATE resources
SET owner_id = (SELECT id FROM users WHERE username = 'vp1')
WHERE owner_id IS NULL;
ALTER TABLE resources
	ALTER COLUMN owner_id SET NOT NULL;

-- Пользователи: 1 admin, 1 dp (директор портфеля), 2 rp (руководители проектов), 4 vp (владельцы процессов), 2 worker
INSERT INTO users (username, name, password_hash, role) VALUES
('rp1', 'РП-1', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'rp'),
('rp2', 'РП-2', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'rp'),
('w1', 'Работник-1', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'worker'),
('w2', 'Работник-2', '$2a$10$xrb4V/Iq3ziY8g1xU9/s/u2dE/MdKdPVD4NdiXnHxNztoEW625lIi', 'worker');

-- Демо-ресурсы принадлежат vp1 (как и все сотрудники).
INSERT INTO resources (title, code, owner_id)
SELECT t.title, t.code, u.id
FROM (
    VALUES
        ('Инженер', 'И'),
        ('Монтажник', 'М'),
        ('Производитель работ', 'ПР'),
        ('Руководитель группы', 'РГ'),
        ('Руководитель службы инсталляции', 'РСИ')
) AS t(title, code)
CROSS JOIN (SELECT id FROM users WHERE username = 'vp1') u;

-- По 1 проекту на РП; код проекта отражает цепочку владения ДП-РП-ВП.
-- V5-триггер при вставке автосоздаст по 2 процесса (Инсталляция + Производство) на проект.
INSERT INTO projects (code, start_date, end_date, priority, owner_id)
VALUES
    ('КО-01_РП1_ВП1', DATE '2026-07-15', DATE '2026-08-30', 1, (SELECT id FROM users WHERE username = 'rp1')),
    ('КО-02_РП2_ВП3', DATE '2026-08-01', DATE '2026-09-20', 2, (SELECT id FROM users WHERE username = 'rp2'));

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

-- Справочник состояний: доступность задаётся is_available.
INSERT INTO states (code, name, is_available) VALUES
('Я', 'Явка', TRUE),
('К', 'Командировка', TRUE),
('ОТГ', 'Отгул', FALSE),
('ОТП', 'Отпуск', FALSE),
('Б', 'Больничный', FALSE);

-- Конкретные сотрудники (уникальные ресурсы) под категории ресурсов.
INSERT INTO employees (resource_id, name, position, hire_date, termination_date)
SELECT r.id, e.name, e.position, e.hire_date, e.termination_date
FROM (
    VALUES
        ('Инженер', 'Иванов Иван Иванович', 'Инженер-конструктор', DATE '2024-01-15', NULL),
        ('Инженер', 'Петров Пётр Петрович', 'Инженер-технолог', DATE '2024-03-01', NULL),
        ('Инженер', 'Сидоров Сидор Сидорович', 'Инженер-проектировщик', DATE '2025-05-10', NULL),
        ('Инженер', 'Кузнецов Кузьма Кузьмич', 'Инженер-электрик', DATE '2025-08-20', NULL),
        ('Инженер', 'Смирнов Смирн Смирнович', 'Инженер-сметчик', DATE '2026-01-12', NULL),
        ('Инженер', 'Волков Волк Волкович', 'Ведущий инженер', DATE '2026-02-01', NULL),
        ('Инженер', 'Егоров Егор Егорович', 'Инженер ПНР', DATE '2026-06-05', NULL),
        ('Монтажник', 'Фёдоров Фёдор Фёдорович', 'Монтажник слаботочных систем', DATE '2024-02-10', NULL),
        ('Монтажник', 'Макаров Макар Макарович', 'Монтажник-высотник', DATE '2025-04-15', NULL),
        ('Монтажник', 'Андреев Андрей Андреевич', 'Монтажник СКС', DATE '2026-03-20', NULL),
        ('Монтажник', 'Борисов Борис Борисович', 'Монтажник-электромонтажник', DATE '2024-11-01', DATE '2026-06-30'),
        ('Производитель работ', 'Григорьев Григорий Григорьевич', 'Производитель работ', DATE '2024-05-05', NULL),
        ('Производитель работ', 'Николаев Николай Николаевич', 'Старший производитель работ', DATE '2025-09-01', NULL),
        ('Руководитель группы', 'Дмитриев Дмитрий Дмитриевич', 'Руководитель группы монтажа', DATE '2023-10-01', NULL),
        ('Руководитель группы', 'Степанов Степан Степанович', 'Руководитель группы ПНР', DATE '2024-06-15', NULL),
        ('Руководитель группы', 'Ильин Илья Ильич', 'Руководитель группы проектирования', DATE '2025-01-20', NULL),
        ('Руководитель группы', 'Захаров Захар Захарович', 'Руководитель группы инженерии', DATE '2026-04-01', NULL),
        ('Руководитель службы инсталляции', 'Семёнов Семён Семёнович', 'Руководитель службы инсталляции', DATE '2022-09-01', NULL)
) AS e(resource_title, name, position, hire_date, termination_date)
JOIN resources r ON r.title = e.resource_title;

-- Все сотрудники закреплены за vp1 (владелец процесса Инсталляция).
UPDATE employees
SET manager_id = (SELECT id FROM users WHERE username = 'vp1');

-- Демо-состояния (диапазоны) на период проекта КО-01 (2026-07-15..08-30):
-- 2 инженера в отпуске, 1 на больничном, 1 монтажник в командировке (доступен).
INSERT INTO employee_states (employee_id, state_id, start_date, end_date)
SELECT e.id, s.id, x.start_date, x.end_date
FROM (
    VALUES
        ('Иванов Иван Иванович', 'ОТП', DATE '2026-07-20', DATE '2026-08-02'),
        ('Петров Пётр Петрович', 'ОТП', DATE '2026-08-10', DATE '2026-08-24'),
        ('Сидоров Сидор Сидорович', 'Б', DATE '2026-07-27', DATE '2026-08-05'),
        ('Фёдоров Фёдор Фёдорович', 'К', DATE '2026-08-01', DATE '2026-08-07')
) AS x(employee_name, state_code, start_date, end_date)
JOIN employees e ON e.name = x.employee_name
JOIN states s ON s.code = x.state_code;
