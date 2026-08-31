-- Триггеры для task_comments (таблица — в V1, индексы — в V2).
-- Каскадное мягкое удаление: задача → её комментарии.
CREATE OR REPLACE FUNCTION cascade_soft_delete_task_comments()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE task_comments
    SET deleted_at = NOW()
    WHERE task_id = OLD.id
      AND deleted_at IS NULL;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_cascade_soft_delete_task_comments
AFTER UPDATE OF deleted_at ON tasks
FOR EACH ROW
WHEN (OLD.deleted_at IS NULL AND NEW.deleted_at IS NOT NULL)
EXECUTE FUNCTION cascade_soft_delete_task_comments();

-- Запрет жёсткого удаления (функция block_hard_delete объявлена в V5).
CREATE TRIGGER block_hard_delete_on_task_comments
BEFORE DELETE ON task_comments
FOR EACH ROW
EXECUTE FUNCTION block_hard_delete();