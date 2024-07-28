-- CreateVerifyRequest Создаём новый запрос на верификацию почты
-- name: CreateVerifyRequest :one
INSERT INTO verify_emails (id_user, secret_key)
VALUES (@id_user::integer, @secret_key::text)
RETURNING *;

-- DeleteVerifyRequest Удаляем запрос на верификацию
-- name: DeleteVerifyRequest :one
DELETE FROM verify_emails
WHERE id_user = @id_user::integer
RETURNING *;

-- GetVerifyRequest Возвращаем запрос на верификацию
-- name: GetVerifyRequest :one
SELECT * FROM verify_emails
WHERE id_user = @id_user::integer;
