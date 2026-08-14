-- V10: запрет кольцевых зависимостей в иерархии руководителей (users.manager_id).
-- Триггер обходит цепочку manager_id вверх от нового руководителя; если она
-- достигает самого пользователя — кольцевая зависимость, ошибка. Запрещает и
-- прямое самоподчинение (manager_id = id). Срабатывает на любом пути записи
-- (API, SQL, миграции/сиды).
CREATE OR REPLACE FUNCTION prevent_manager_cycle() RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
    cur   BIGINT;
    depth INT := 0;
BEGIN
    IF NEW.manager_id IS NULL THEN
        RETURN NEW;
    END IF;

    IF NEW.manager_id = NEW.id THEN
        RAISE EXCEPTION 'manager cannot be the user themself';
    END IF;

    -- Обход вверх по цепочке руководителей (без фильтра по deleted_at, чтобы
    -- цикл не «проскочил» через мягко удалённые записи).
    cur := NEW.manager_id;
    WHILE cur IS NOT NULL LOOP
        depth := depth + 1;
        IF depth > 1000 THEN
            RAISE EXCEPTION 'manager hierarchy is too deep';
        END IF;
        IF cur = NEW.id THEN
            RAISE EXCEPTION 'circular manager dependency is not allowed';
        END IF;
        SELECT manager_id INTO cur FROM users WHERE id = cur;
    END LOOP;

    RETURN NEW;
END;
$$;

CREATE TRIGGER trg_users_prevent_manager_cycle
BEFORE INSERT OR UPDATE OF manager_id ON users
FOR EACH ROW
WHEN (NEW.manager_id IS NOT NULL)
EXECUTE FUNCTION prevent_manager_cycle();
