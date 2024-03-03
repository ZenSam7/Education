-- CreateComment Создаём комментарий к статье
-- name: CreateComment :one
WITH add_comment AS ( -- Я знаю что есть тразнакции
    UPDATE articles
    SET comments = array_append(comments, lastval())
    WHERE id_article = $1
)
INSERT INTO comments (text, from_user)
VALUES ($3, $2)
RETURNING *;

-- DeleteComment Удаляем комментарий к статье
-- name: DeleteComment :one
WITH deleted_comment_id AS ( -- Я знаю что есть тразнакции
    DELETE FROM comments
    WHERE id_comment = sqlc.arg(id_comment)::integer
)
UPDATE articles
SET comments = array_remove(comments, sqlc.arg(id_comment)::integer)
WHERE sqlc.arg(id_comment)::integer = ANY(comments)
RETURNING *;

-- GetComment Возвращаем комментарий
-- name: GetComment :one
SELECT * FROM comments
WHERE id_comment = sqlc.arg(id_comment)::integer;

-- EditCommentParam Изменяем параметр(ы) пользователя
-- name: EditCommentParam :one
UPDATE comments
SET
    -- Если изменили текст или автора польщователя то обновляем его
    edited_at = CASE WHEN sqlc.arg(text)::text <> '' THEN NOW()
                     WHEN sqlc.arg(from_user)::integer <> from_user THEN NOW()
                     ELSE edited_at END,

    -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
    -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
    text = CASE WHEN sqlc.arg(text)::text <> '' THEN sqlc.arg(text)::text ELSE text END,
    from_user = CASE WHEN sqlc.arg(from_user)::integer <> from_user
                     THEN sqlc.arg(from_user)::integer
                     ELSE from_user END,
    evaluation = CASE WHEN sqlc.arg(evaluation)::integer <> evaluation
                      THEN sqlc.arg(evaluation)::integer
                      ELSE evaluation END
WHERE id_comment = sqlc.arg(id_comment)::integer
RETURNING *;

