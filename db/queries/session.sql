-- CreateSession Создаём сессию
-- name: CreateSession :one
INSERT INTO sessions (expired_at, refresh_token, id_user, id_session, client_ip)
VALUES (sqlc.arg(expired_at)::timestamptz,
        sqlc.arg(refresh_token)::text,
        sqlc.arg(id_user)::integer,
        sqlc.arg(id_session)::uuid,
        sqlc.arg(client_ip)::text)
RETURNING *;

-- DeleteSession Удаляем сессию по id
-- name: DeleteSession :one
DELETE FROM sessions
WHERE id_session = sqlc.arg(id_session)::uuid
RETURNING *;

-- BlockSession Блокируем сессию по id
-- name: BlockSession :one
UPDATE sessions
SET blocked = true
WHERE id_session = sqlc.arg(id_session)::uuid
RETURNING *;

-- GetSession Получаем сессиб по id
-- name: GetSession :one
SELECT * FROM sessions
WHERE id_session = sqlc.arg(id_session)::uuid;
