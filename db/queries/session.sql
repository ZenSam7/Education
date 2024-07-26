-- CreateSession Создаём сессию
-- name: CreateSession :one
INSERT INTO sessions (expired_at, refresh_token, id_user, id_session, client_ip)
VALUES (@expired_at::timestamptz,
        @refresh_token::text,
        @id_user::integer,
        @id_session::uuid,
        @client_ip::text)
RETURNING *;

-- DeleteSession Удаляем сессию по id
-- name: DeleteSession :one
DELETE FROM sessions
WHERE id_session = @id_session::uuid
RETURNING *;

-- BlockSession Блокируем сессию по id
-- name: BlockSession :one
UPDATE sessions
SET blocked = true
WHERE id_session = @id_session::uuid
RETURNING *;

-- GetSession Получаем сессиб по id
-- name: GetSession :one
SELECT * FROM sessions
WHERE id_session = @id_session::uuid;

-- CountRowsSessions Считаем количество строк в таблице
-- name: CountRowsSessions :one
SELECT COUNT(*) FROM sessions;