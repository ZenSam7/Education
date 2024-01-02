-- CreateUser Создаём одного пользователя
-- name: CreateUser :one
INSERT INTO users (
  name, description, email, karma
) VALUES (
  $1, $2, $3, 0
)
RETURNING *;

-- GetUser Возвращаем пользователя по его id
-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 ;

-- GetManyUsers Возвращаем слайс пользователей (сортируем по дате создания)
-- name: GetManyUsers :many
SELECT * FROM users
ORDER BY created_at
LIMIT $1
OFFSET $2;

-- UpdateUserName Обновляем имя пользователя по его id
-- name: UpdateUserName :exec
UPDATE users
SET name = $2
WHERE id = $1;

-- DeleteUser Удаляем пользователя по имени
-- name: DeleteUser :exec
DELETE FROM users WHERE name = $1;
