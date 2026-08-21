-- Обновление 12: серверные refresh-сессии (ротация/отзыв refresh-токенов, AD-06).
-- Храним только SHA-256 токена; сам opaque-токен живёт в HttpOnly-куке (AD-05).
CREATE TABLE refresh_sessions (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    -- Цепочка ротации: id сессии, которой эта была заменена (для детекта повторного использования).
    replaced_by BIGINT REFERENCES refresh_sessions(id)
);

CREATE INDEX refresh_sessions_user_idx    ON refresh_sessions (user_id);
CREATE INDEX refresh_sessions_expires_idx ON refresh_sessions (expires_at);