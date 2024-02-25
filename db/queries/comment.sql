-- CreateComment Создаём комментарий к статье
-- name: CreateComment :exec
WITH add_comment AS ( -- Я знаю что есть тразнакции
    UPDATE articles
    SET comments = array_append(comments, lastval())
    WHERE id_article = $1
)
INSERT INTO comments (text, from_user)
VALUES ($3, $2);

-- DeleteComment Удаляем комментарий к статье
-- name: DeleteComment :exec
WITH deleted_comment_id AS ( -- Я знаю что есть тразнакции
    DELETE FROM comments
    WHERE id_comment = sqlc.arg(id_comment)::integer
)
UPDATE articles
SET comments = array_remove(comments, sqlc.arg(id_comment)::integer)
WHERE id_article = $1;

-- GetComment Возвращаем комментарий
-- name: GetComment :exec
SELECT * FROM comments
WHERE id_comment = $1;

-- EditCommentParam Изменяем параметр(ы) пользователя
-- name: EditCommentParam :one
UPDATE comments
SET
  edited_at = COALESCE(sqlc.arg(edited_at)::timestamp, edited_at),

  -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
  -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
  text = CASE WHEN sqlc.arg(text)::text <> '' THEN sqlc.arg(text)::text ELSE text END,
  from_user = COALESCE(sqlc.arg(from_user)::integer, from_user),
  evaluation = COALESCE(sqlc.arg(evaluation)::integer, evaluation)
WHERE id_comment = sqlc.arg(id_comment)::integer
RETURNING *;

