-- =============================================
-- AUTO-CREATE TRIGGER: template task operations (subtasks)
-- =============================================
-- Replaces the V8/V11 function: a template task may now carry an optional
-- "operations" array. Each operation becomes a subtask of the created task
-- (parent_id = the task id), inheriting the project dates and the task's
-- process, with status 'not_started' (subtasks never carry their own dates
-- or resources on the frontend, so the template only stores titles).
CREATE OR REPLACE FUNCTION fn_project_auto_create() RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
	cfg     RECORD;
	p       RECORD;
	tk      RECORD;
	op      RECORD;
	pid     BIGINT;
	tid     BIGINT;
BEGIN
	SELECT * INTO cfg FROM project_auto_create ORDER BY id LIMIT 1;
	IF NOT FOUND OR NOT cfg.enabled THEN
		RETURN NEW;
	END IF;

	FOR p IN
		SELECT value->>'title' AS title,
		       (value->>'owner_id')::bigint AS owner_id,
		       value->>'color' AS color,
		       value->'tasks' AS tasks,
		       ord
		FROM jsonb_array_elements(cfg.config) WITH ORDINALITY AS pp(value, ord)
	LOOP
		INSERT INTO processes (project_id, title, start_date, end_date, owner_id, color, sort_order)
		VALUES (NEW.id, p.title, NEW.start_date, NEW.end_date, p.owner_id, NULLIF(p.color, ''), p.ord)
		RETURNING id INTO pid;

		FOR tk IN
			SELECT value, ord
			FROM jsonb_array_elements(p.tasks) WITH ORDINALITY AS tt(value, ord)
		LOOP
			INSERT INTO tasks (process_id, title, start_date, end_date, color, sort_order)
			VALUES (pid, tk.value->>'title', NEW.start_date, NEW.end_date, NULLIF(tk.value->>'color', ''), tk.ord)
			RETURNING id INTO tid;

			INSERT INTO assignments (task_id, resource_id, quantity)
			SELECT tid, (res.value->>'resource_id')::bigint, (res.value->>'quantity')::int
			FROM jsonb_array_elements(tk.value->'resources') res
			WHERE (res.value->>'resource_id')::bigint IS NOT NULL;

			-- Operations → subtasks of the just-created task. Each inherits
			-- the project dates and the parent's process; status is the
			-- default 'not_started'; no color.
			FOR op IN
				SELECT value->>'title' AS title, ord
				FROM jsonb_array_elements(tk.value->'operations') WITH ORDINALITY AS oo(value, ord)
			LOOP
				IF op.title IS NOT NULL AND op.title <> '' THEN
					INSERT INTO tasks (process_id, parent_id, title, start_date, end_date, color, status, sort_order)
					VALUES (pid, tid, op.title, NEW.start_date, NEW.end_date, NULL, 'not_started', op.ord);
				END IF;
			END LOOP;
		END LOOP;
	END LOOP;

	RETURN NEW;
END;
$$;