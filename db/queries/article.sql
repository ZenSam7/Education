-- CreateArticle Создаём статью
-- name: CreateArticle :one
INSERT INTO articles (title, text, authors)
VALUES ($1, $2, $3)
RETURNING *;

-- GetArticle Возвращаем статью по id
-- name: GetArticle :one
SELECT * FROM articles
WHERE id_article = $1;

-- GetArticlesWithAttribute Возвращаем много статей взятых по какому-то признаку(ам)
-- name: GetArticlesWithAttribute :many
SELECT * FROM articles
WHERE
    -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
    -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
    -- (Кстати, "::text" <- эти штуки нужны чтобы вместа pgtype был string/int32)
    title = CASE WHEN (sqlc.arg(title)::text <> '') AND (sqlc.arg(title)::text IS NOT NULL)
        THEN sqlc.arg(title)::text
        ELSE title END AND
    text = CASE WHEN (sqlc.arg(text)::text <> '') AND (sqlc.arg(text)::text IS NOT NULL)
        THEN sqlc.arg(text)::text
        ELSE text END AND
    evaluation = CASE WHEN sqlc.arg(evaluation)::integer <> evaluation
        THEN sqlc.arg(evaluation)::integer
        ELSE evaluation END
LIMIT sqlc.arg('Limit')::integer
OFFSET sqlc.arg('Offset')::integer;

-- GetManySortedArticles Возвращаем много отсортированных статей
-- name: GetManySortedArticles :many
SELECT * FROM articles
ORDER BY
        CASE WHEN sqlc.arg(id_article)::boolean THEN id_article::integer
             WHEN sqlc.arg(evaluation)::boolean THEN evaluation::integer END
        , -- запятая
        CASE WHEN sqlc.arg(comments)::boolean THEN comments::integer[]
             WHEN sqlc.arg(authors)::boolean THEN authors::integer[] END
        , -- запятая
        CASE WHEN sqlc.arg(title)::boolean THEN title::text
             WHEN sqlc.arg(text)::boolean THEN text::text END
        , -- запятая
        CASE WHEN sqlc.arg(edited_at)::boolean THEN edited_at::timestamp
             WHEN sqlc.arg(created_at)::boolean THEN created_at::timestamp END
LIMIT sqlc.arg('Limit')::integer
OFFSET sqlc.arg('Offset')::integer;

-- GetManySortedArticlesWithAttribute Возвращаем много статей взятых по признаку по
-- какому-то признаку(ам) и отсортированных по другому признаку(ам)
-- name: GetManySortedArticlesWithAttribute :many
SELECT * FROM articles
WHERE
    -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
    -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
    -- (Кстати, "::text" <- эти штуки нужны чтобы вместа pgtype был string/int32)
    title = CASE WHEN (sqlc.arg(select_title)::text <> '') AND (sqlc.arg(select_title)::text IS NOT NULL)
        THEN sqlc.arg(select_title)::text
        ELSE title END AND
    text = CASE WHEN (sqlc.arg(select_text)::text <> '') AND (sqlc.arg(select_text)::text IS NOT NULL)
        THEN sqlc.arg(select_text)::text
        ELSE text END AND
    evaluation = CASE WHEN sqlc.arg(select_evaluation)::integer <> evaluation
        THEN sqlc.arg(select_evaluation)::integer
        ELSE evaluation END
ORDER BY
        CASE WHEN sqlc.arg(sorted_id_article)::boolean THEN id_article::integer
             WHEN sqlc.arg(sorted_evaluation)::boolean THEN evaluation::integer END
        , -- запятая
        CASE WHEN sqlc.arg(sorted_comments)::boolean THEN comments::integer[]
             WHEN sqlc.arg(sorted_authors)::boolean THEN authors::integer[] END
        , -- запятая
        CASE WHEN sqlc.arg(sorted_title)::boolean THEN title::text
             WHEN sqlc.arg(sorted_text)::boolean THEN text::text END
        , -- запятая
        CASE WHEN sqlc.arg(sorted_edited_at)::boolean THEN edited_at::timestamp
             WHEN sqlc.arg(sorted_created_at)::boolean THEN created_at::timestamp END

LIMIT sqlc.arg('Limit')::integer
OFFSET sqlc.arg('Offset')::integer;

-- EditArticle Изменяем параметр(ы) статьи
-- name: EditArticle :one
UPDATE articles
SET
  -- Если изменили текст или заголовок то обновляем время изменения
  edited_at = CASE WHEN (sqlc.arg(text)::text <> '') THEN NOW()
                   WHEN (sqlc.arg(title)::text <> '') THEN NOW()
                   ELSE edited_at END,

  -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
  -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
  -- (Кстати, "::text" <- эти штуки нужны чтобы вместа pgtype был string/int32)
  title = CASE WHEN sqlc.arg(title)::text <> '' THEN sqlc.arg(title)::text ELSE title END,
  text = CASE WHEN sqlc.arg(text)::text <> '' THEN sqlc.arg(text)::text ELSE text END,
  comments = COALESCE(sqlc.arg(comments), comments),
  authors = COALESCE(sqlc.arg(authors), authors)
WHERE id_article = sqlc.arg(id_article)::integer
RETURNING *;

-- DeleteArticle Удаляем статью и комментарии к ней
-- name: DeleteArticle :one
WITH deleted_comments AS ( -- Объединяем 2 запроса в 1
    DELETE FROM comments
    WHERE id_comment = ANY ((SELECT comments FROM articles
                            WHERE id_article = sqlc.arg(id_article)::integer)::text::integer[])
), update_id AS ( -- Объединяем 3 запроса в 1
    UPDATE articles
    SET id_article = id_article - 1
    WHERE id_article > sqlc.arg(id_article)::integer
)
DELETE FROM articles
WHERE id_article = sqlc.arg(id_article)::integer
RETURNING *;
