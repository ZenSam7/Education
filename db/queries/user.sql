-- CreateUser Создаём пользователя
-- name: CreateUser :one
INSERT INTO users (name, description)
VALUES (sqlc.arg(name)::text, sqlc.arg(description)::text)
RETURNING *;

-- DeleteUser Удаляем пользователя и сдвигаем id
-- name: DeleteUser :one
DELETE FROM users
WHERE id_user = $1
RETURNING *;

-- GetUser Возвращаем пользователя
-- name: GetUser :one
SELECT * FROM users
WHERE id_user = $1;

-- GetManySortedUsers Возвращаем слайс пользователей отсортированных по параметру attribute
-- (можно поставить: id_user, и сортировки не будет)
-- name: GetManySortedUsers :many
SELECT * FROM users
ORDER BY sqlc.arg(attribute)::text
LIMIT sqlc.arg('Limit')::integer
OFFSET sqlc.arg('Offset')::integer;

-- EditUserParam Изменяем параметр(ы) пользователя
-- name: EditUserParam :one
UPDATE users
SET
  -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
  -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
  name = CASE WHEN sqlc.arg(name)::text <> '' THEN sqlc.arg(name)::text ELSE name END,
  description = CASE WHEN sqlc.arg(description)::text <> '' THEN sqlc.arg(description)::text ELSE description END,
  karma = COALESCE(sqlc.arg(karma)::integer, karma)
WHERE id_user = sqlc.arg(id_user)::integer
RETURNING *;
