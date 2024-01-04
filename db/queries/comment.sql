-- CreateComment Создаём комментарий к статье
-- name: CreateComment :exec
WITH add_comment AS (
    INSERT INTO comments (text, from_user)
    VALUES ($3, $2)
)
UPDATE articles
SET comments = array_append(comments, currval(pg_get_serial_sequence('comments','id_comment')))
WHERE id_article = $1;

-- DeleteComment Удаляем комментарий к статье
-- name: DeleteComment :exec
WITH deleted_comment_id AS (
    DELETE FROM comments
    WHERE id_comment = $2
)
UPDATE articles
SET comments = array_remove(comments, $2)
WHERE id_article = $1;

-- GetComment Возвращаем комментарий
-- name: GetComment :exec
SELECT * FROM comments
WHERE id_comment = $1;


-- EditCommentText Изменяем текст комментария и обновляем время изменения комментария (Column1 = id_comment, Column1 = text)
-- name: EditCommentText :exec
WITH update_time AS (
    UPDATE comments
    SET edited_at = now()
    WHERE id_comment = $1::integer
)
UPDATE comments
SET text = $2::text
WHERE id_comment = $1::integer;
