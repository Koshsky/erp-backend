
-- Auto-create processes and template tasks when a project is created
CREATE OR REPLACE FUNCTION fn_create_project_templates()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    installation_process_id BIGINT;
BEGIN
    INSERT INTO processes (project_id, title, start_date, end_date)
    VALUES (NEW.id, 'Инсталляция', NEW.start_date, NEW.end_date)
    RETURNING id INTO installation_process_id;

    INSERT INTO processes (project_id, title, start_date, end_date)
    VALUES (NEW.id, 'Производство', NEW.start_date, NEW.end_date);

    INSERT INTO tasks (process_id, title, start_date, end_date)
    VALUES
        (installation_process_id, 'Подписание Акта ввода в эксплуатацию', NEW.start_date, NEW.end_date),
        (installation_process_id, 'Список затраченных материалов', NEW.start_date, NEW.end_date),
        (installation_process_id, 'Тестирование комплекса систем телемедицины (MVS VEGA)', NEW.start_date, NEW.end_date),
        (installation_process_id, 'Проведение инструктажа и передача инструкций мед персоналу', NEW.start_date, NEW.end_date),
        (installation_process_id, 'Пуско-наладочные работы', NEW.start_date, NEW.end_date),
        (installation_process_id, 'Инсталляция оконечного оборудования телемедицины', NEW.start_date, NEW.end_date),
        (installation_process_id, 'Предоставить карту сети', NEW.start_date, NEW.end_date),
        (installation_process_id, 'Монтаж кабеленесущих систем и кабельных трасс', NEW.start_date, NEW.end_date),
        (installation_process_id, 'Осмотр объекта', NEW.start_date, NEW.end_date);

    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_projects_create_templates
AFTER INSERT ON projects
FOR EACH ROW
EXECUTE FUNCTION fn_create_project_templates();



