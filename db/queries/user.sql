-- CreateUser Создаём пользователя
-- name: CreateUser :exec
INSERT INTO users (name, description, email)
VALUES ($1, $2, $3);

-- DeleteUser Удаляем пользователя
-- name: DeleteUser :exec
DELETE FROM users
WHERE id_user = $1;

-- GetUser Возвращаем пользователя
-- name: GetUser :one
SELECT * FROM users
WHERE id_user = $1;

-- GetManyUsers Возвращаем слайс пользователей отсортированных по параметру Column1
-- name: GetManyUsers :many
SELECT * FROM users
ORDER BY $1::text
LIMIT $2
OFFSET $3;


-- EditUserName Изменяем имя пользователя
-- name: EditUserName :exec
UPDATE users
SET name = $2
WHERE id_user = $1;

-- EditUserDescription Изменяем описание пользователя
-- name: EditUserDescription :exec
UPDATE users
SET description = $2
WHERE id_user = $1;

-- EditUserEmail Изменяем почту пользователя
-- name: EditUserEmail :exec
UPDATE users
SET email = $2
WHERE id_user = $1;

-- EditUserKarma Изменяем карму пользователя
-- name: EditUserKarma :exec
UPDATE users
SET karma = $2
WHERE id_user = $1;
