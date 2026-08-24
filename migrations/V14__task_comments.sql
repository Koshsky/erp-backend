-- Комментарии к задачам: цепочки обсуждений (ответы на ответы — через parent_id).
CREATE TABLE task_comments (
	id BIGSERIAL PRIMARY KEY,
	task_id BIGINT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
	-- Отправитель комментария (автор берётся из контекста авторизации).
	author_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
	-- Ответ на другой комментарий той же задачи; NULL — корневой комментарий.
	parent_id BIGINT REFERENCES task_comments(id) ON DELETE CASCADE,
	content TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	deleted_at TIMESTAMPTZ DEFAULT NULL
);

CREATE INDEX idx_task_comments_task_id ON task_comments(task_id);
CREATE INDEX idx_task_comments_task_created ON task_comments(task_id, created_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_task_comments_parent_id ON task_comments(parent_id) WHERE deleted_at IS NULL;

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

-- Запрет жёсткого удаления (функция block_hard_delete объявлена в V7).
CREATE TRIGGER block_hard_delete_on_task_comments
BEFORE DELETE ON task_comments
FOR EACH ROW
EXECUTE FUNCTION block_hard_delete();