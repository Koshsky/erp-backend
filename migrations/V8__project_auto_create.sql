-- Configurable auto-creation of processes/tasks when a project is inserted (the
-- project_auto_create table lives in V1). The config is a single row: enabled + a JSONB template
--   [ { "title": "Installation", "owner_id": 5, "tasks": [
--       { "title": "Site inspection", "resources": [ { "resource_id": 1, "quantity": 2 } ] }, ...
--     ] }, ... ]
-- The template array position is written into sort_order (jsonb arrays keep
-- their order; WITH ORDINALITY numbers the elements), so the admin-defined
-- template order is preserved on creation (processes within the project, tasks
-- within the process).
CREATE OR REPLACE FUNCTION fn_project_auto_create() RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
	cfg      RECORD;
	p        RECORD;
	tk       RECORD;
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
		       value->'tasks' AS tasks,
		       ord
		FROM jsonb_array_elements(cfg.config) WITH ORDINALITY AS pp(value, ord)
	LOOP
		INSERT INTO processes (project_id, title, start_date, end_date, owner_id, sort_order)
		VALUES (NEW.id, p.title, NEW.start_date, NEW.end_date, p.owner_id, p.ord)
		RETURNING id INTO pid;

		FOR tk IN
			SELECT value, ord
			FROM jsonb_array_elements(p.tasks) WITH ORDINALITY AS tt(value, ord)
		LOOP
			INSERT INTO tasks (process_id, title, start_date, end_date, sort_order)
			VALUES (pid, tk.value->>'title', NEW.start_date, NEW.end_date, tk.ord)
			RETURNING id INTO tid;

			INSERT INTO assignments (task_id, resource_id, quantity)
			SELECT tid, (res.value->>'resource_id')::bigint, (res.value->>'quantity')::int
			FROM jsonb_array_elements(tk.value->'resources') res
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