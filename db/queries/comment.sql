-- CreateComment Создаём комментарий к статье
-- name: CreateComment :one
WITH add_comment AS ( -- TODO: переделать в транзакцию
    UPDATE articles
    SET comments = array_append(comments, lastval())
    WHERE id_article = $1
)
INSERT INTO comments (text, author)
VALUES ($3, $2)
RETURNING *;

-- DeleteComment Удаляем комментарий к статье
-- name: DeleteComment :one
WITH deleted_comment_id AS ( -- TODO: переделать в транзакцию
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
    -- Если изменили текст или автора пользователя, то обновляем его
    edited_at = CASE WHEN @text::text <> '' THEN NOW()
                     WHEN @author::integer <> author THEN NOW()
                     ELSE edited_at END,

    -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
    -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
    text = CASE WHEN @text::text <> '' THEN @text::text ELSE text END,
    evaluation = CASE WHEN @evaluation::integer <> evaluation THEN @evaluation::integer ELSE evaluation END
WHERE id_comment = @id_comment::integer
RETURNING *;

-- GetCommentsOfArticle Возвращаем комментарии к статье
-- name: GetCommentsOfArticle :many
WITH the_article AS ( -- TODO: переделать в транзакцию
    SELECT unnest(comments) AS id_comment FROM articles
    WHERE id_article = @id_article::integer
)
SELECT * FROM comments
WHERE id_comment IN (SELECT id_comment FROM the_article)
OFFSET sqlc.arg('Offset')::integer
LIMIT sqlc.arg('Limit')::integer;

-- CountRowsArticle Считаем количество строк в таблице
-- name: CountRowsArticle :one
SELECT COUNT(*) FROM articles;
