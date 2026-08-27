-- Конфигурируемое автосоздание процессов/задач при вставке проекта (таблица
-- project_auto_create — в V1). Конфиг — одна строка: enabled + JSONB-шаблон
--   [ { "title": "Инсталляция", "owner_id": 5, "tasks": [
--       { "title": "Осмотр объекта", "resources": [ { "resource_id": 1, "quantity": 2 } ] }, ...
--     ] }, ... ]
CREATE OR REPLACE FUNCTION fn_project_auto_create() RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
	cfg      RECORD;
	p        RECORD;
	t        RECORD;
	pid      BIGINT;
	tid      BIGINT;
BEGIN
	SELECT * INTO cfg FROM project_auto_create ORDER BY id LIMIT 1;
	IF NOT FOUND OR NOT cfg.enabled THEN
		RETURN NEW;
	END IF;

	FOR p IN
		SELECT value->>'title' AS title,
		       (value->>'owner_id')::bigint AS owner_id,
		       value->'tasks' AS tasks
		FROM jsonb_array_elements(cfg.config)
	LOOP
		INSERT INTO processes (project_id, title, start_date, end_date, owner_id)
		VALUES (NEW.id, p.title, NEW.start_date, NEW.end_date, p.owner_id)
		RETURNING id INTO pid;

		FOR t IN SELECT value FROM jsonb_array_elements(p.tasks) LOOP
			INSERT INTO tasks (process_id, title, start_date, end_date)
			VALUES (pid, t.value->>'title', NEW.start_date, NEW.end_date)
			RETURNING id INTO tid;

			INSERT INTO assignments (task_id, resource_id, quantity)
			SELECT tid, (res.value->>'resource_id')::bigint, (res.value->>'quantity')::int
			FROM jsonb_array_elements(t.value->'resources') res
			WHERE (res.value->>'resource_id')::bigint IS NOT NULL;
		END LOOP;
	END LOOP;

	RETURN NEW;
END;
$$;

CREATE TRIGGER trg_project_auto_create
AFTER INSERT ON projects
FOR EACH ROW
EXECUTE FUNCTION fn_project_auto_create();
