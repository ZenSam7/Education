-- CreateUser Создаём пользователя
-- name: CreateUser :one
INSERT INTO users (name, email, password_hash)
VALUES (sqlc.arg(name)::text, sqlc.arg(email)::text, sqlc.arg(password_hash)::text)
RETURNING *;

-- DeleteUser Удаляем пользователя и сдвигаем id
-- name: DeleteUser :one
WITH update_id AS ( -- Объединяем 2 запроса в 1
    UPDATE users
    SET id_user = id_user - 1
    WHERE id_user > sqlc.arg(id_user)::integer
)
DELETE FROM users
WHERE id_user = sqlc.arg(id_user)::integer
RETURNING *;

-- GetUser Возвращаем пользователя
-- name: GetUser :one
SELECT * FROM users
WHERE id_user = $1;

-- GetManySortedUsers Возвращаем слайс пользователей отсортированных по какому-то параметру
-- (можно поставить: id_user, и сортировки не будет)
-- name: GetManySortedUsers :many
SELECT * FROM users
ORDER BY
        CASE WHEN sqlc.arg(id_user)::boolean THEN id_user::integer
             WHEN sqlc.arg(karma)::boolean THEN karma::integer END
        , -- запятая
        CASE WHEN sqlc.arg(name)::boolean THEN name::text
             WHEN sqlc.arg(description)::boolean THEN description::text END
LIMIT sqlc.arg('Limit')::integer
OFFSET sqlc.arg('Offset')::integer;

-- EditUser Изменяем параметр(ы) пользователя
-- name: EditUser :one
UPDATE users
SET
  -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
  -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
  name = CASE WHEN sqlc.arg(name)::text <> '' THEN sqlc.arg(name)::text ELSE name END,
  description = CASE WHEN sqlc.arg(description)::text <> '' THEN sqlc.arg(description)::text ELSE description END,
  karma = COALESCE(sqlc.arg(karma)::integer, karma)
WHERE id_user = sqlc.arg(id_user)::integer
RETURNING *;
