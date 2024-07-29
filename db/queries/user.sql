-- CreateUser Создаём пользователя
-- name: CreateUser :one
INSERT INTO users (name, email, password_hash)
VALUES (@name::text, @email::text, @password_hash::text)
RETURNING *;

-- DeleteUser Удаляем пользователя
-- name: DeleteUser :one
DELETE FROM users
WHERE id_user = @id_user::integer
RETURNING *;

-- GetUser Возвращаем пользователя
-- name: GetUser :one
SELECT * FROM users
WHERE id_user = $1;

-- GetUserFromName Возвращаем пользователя по имени
-- name: GetUserFromName :one
SELECT * FROM users
WHERE name = $1;

-- GetManySortedUsers Возвращаем слайс пользователей отсортированных по какому-то параметру
-- (можно поставить: id_user, и сортировки не будет)
-- name: GetManySortedUsers :many
SELECT * FROM users
ORDER BY
        CASE WHEN @id_user::boolean THEN id_user::integer
             WHEN @karma::boolean THEN karma::integer END
        , -- запятая
        CASE WHEN @name::boolean THEN name::text
             WHEN @description::boolean THEN description::text END
LIMIT sqlc.arg('Limit')::integer
OFFSET sqlc.arg('Offset')::integer;

-- EditUser Изменяем параметр(ы) пользователя
-- name: EditUser :one
UPDATE users
SET
  -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
  -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
  -- CASE WHEN используется когда нельзя указать нулевое значение (пустую строку), COALESCE когда можно
  name = CASE WHEN @name::text <> '' THEN @name::text ELSE name END,
  description = COALESCE(sqlc.narg(description)::text, description),
  karma = COALESCE(sqlc.narg(karma)::integer, karma)
WHERE id_user = @id_user::integer
RETURNING *;

-- SetEmailIsVerified Ставим состояние почту как подтверждённую для какого-то пользователя
-- name: SetEmailIsVerified :one
UPDATE users
SET email_verified = true
WHERE id_user = @id_user::integer
RETURNING *;

-- CountRowsUser Считаем количество строк в таблице
-- name: CountRowsUser :one
SELECT COUNT(*) FROM users;
