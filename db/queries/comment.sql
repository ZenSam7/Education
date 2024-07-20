-- CreateComment Создаём комментарий к статье
-- name: CreateComment :one
WITH add_comment AS ( -- Я знаю что есть тразнакции
    UPDATE articles
    SET comments = array_append(comments, lastval())
    WHERE id_article = $1
)
INSERT INTO comments (text, author)
VALUES ($3, $2)
RETURNING *;

-- DeleteComment Удаляем комментарий к статье
-- name: DeleteComment :one
WITH deleted_comment_id AS ( -- Я знаю что есть тразнакции
    DELETE FROM comments
    WHERE id_comment = @id_comment::integer
)
UPDATE articles
SET comments = array_remove(comments, @id_comment::integer)
WHERE @id_comment::integer = ANY(comments)
RETURNING *;

-- GetComment Возвращаем комментарий
-- name: GetComment :one
SELECT * FROM comments
WHERE id_comment = @id_comment::integer;

-- EditComment Изменяем параметр(ы) пользователя
-- name: EditComment :one
UPDATE comments
SET
    -- Если изменили текст или автора польщователя то обновляем его
    edited_at = CASE WHEN @text::text <> '' THEN NOW()
                     WHEN @author::integer <> author THEN NOW()
                     ELSE edited_at END,

    -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
    -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
    text = CASE WHEN @text::text <> '' THEN @text::text ELSE text END,
    author = CASE WHEN @author::integer <> author
                     THEN @author::integer
                     ELSE author END
WHERE id_comment = @id_comment::integer
RETURNING *;

-- GetCommentsOfArticle Возвращаем комментарии
-- name: GetCommentsOfArticle :many
WITH the_article AS (
    SELECT unnest(comments) AS id_comment FROM articles
    WHERE id_article = @id_article::integer
)
SELECT * FROM comments
WHERE id_comment IN (SELECT id_comment FROM the_article)
OFFSET sqlc.arg('Offset')::integer
LIMIT sqlc.arg('Limit')::integer;
