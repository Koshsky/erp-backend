-- V13: idempotency keys.
-- Позволяет сделать create-эндпоинты идемпотентными: клиент присылает
-- заголовок Idempotency-Key (UUID); сервер атомарно захватывает ключ и при
-- повторе с тем же ключом возвращает сохранённый ответ вместо повторного
-- выполнения операции (не создавая новую запись).
--
-- response_status/response_body хранят ответ первого завершённого вызова и
-- заполняются после выполнения операции (claim-фаза: response_status = 0,
-- response_body = NULL). expires_at — TTL очистки старых ключей.
--
-- PK (key, user_id, method, path): изоляция по клиенту/маршруту — один и тот же
-- ключ у разных пользователей или на разных маршрутах не пересекается
-- (исключает кросс-пользовательский и кросс-маршрутный replay).
CREATE TABLE idempotency_keys (
    key             TEXT NOT NULL,
    user_id         BIGINT NOT NULL,
    method          TEXT NOT NULL,
    path            TEXT NOT NULL,
    response_status INTEGER NOT NULL DEFAULT 0,
    response_body   JSONB,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at      TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (key, user_id, method, path)
);

CREATE INDEX idx_idempotency_keys_expires ON idempotency_keys (expires_at);
