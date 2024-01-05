-- CreateUser Создаём пользователя
-- name: CreateUser :one
INSERT INTO users (name, description)
VALUES (sqlc.arg(name)::text, sqlc.arg(description)::text)
RETURNING *;

-- DeleteUser Удаляем пользователя
-- name: DeleteUser :one
DELETE FROM users
WHERE id_user = $1
RETURNING *;

-- GetUser Возвращаем пользователя
-- name: GetUser :one
SELECT * FROM users
WHERE id_user = $1;

-- GetManyUsers Возвращаем слайс пользователей отсортированных по параметру attribute
-- (можно поставить: id_user, и сортировки не будет)
-- name: GetManyUsers :many
SELECT * FROM users
ORDER BY sqlc.arg(attribute)::text
LIMIT $1
OFFSET $2;


-- EditUserParam Изменяем параметр(ы) пользователя
-- name: EditUserParam :one
UPDATE users
SET
  name = COALESCE(sqlc.arg(name)::text, name),
  description = COALESCE(sqlc.arg(description)::text, description),
  karma = COALESCE(sqlc.arg(karma)::integer, karma)
WHERE id_user = sqlc.arg(id_user)::integer
RETURNING *;
