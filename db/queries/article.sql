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
    title = CASE WHEN (@title::text <> '') AND (@title::text IS NOT NULL)
        THEN @title::text
        ELSE title END AND
    text = CASE WHEN (@text::text <> '') AND (@text::text IS NOT NULL)
        THEN @text::text
        ELSE text END AND
    evaluation = CASE WHEN @evaluation::integer <> evaluation
        THEN @evaluation::integer
        ELSE evaluation END  AND
    authors = CASE WHEN @authors::integer[] <> authors
        THEN @select_authors::integer[]
        ELSE authors END
LIMIT sqlc.arg('Limit')::integer
OFFSET sqlc.arg('Offset')::integer;

-- GetManySortedArticles Возвращаем много отсортированных статей
-- name: GetManySortedArticles :many
SELECT * FROM articles
ORDER BY
        CASE WHEN @id_article::boolean THEN id_article::integer
             WHEN @evaluation::boolean THEN evaluation::integer END
        , -- запятая
        CASE WHEN @comments::boolean THEN comments::integer[]
             WHEN @authors::boolean THEN authors::integer[] END
        , -- запятая
        CASE WHEN @title::boolean THEN title::text
             WHEN @text::boolean THEN text::text END
        , -- запятая
        CASE WHEN @edited_at::boolean THEN edited_at::timestamp
             WHEN @created_at::boolean THEN created_at::timestamp END
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
    title = CASE WHEN (@select_title::text <> '') AND (@select_title::text IS NOT NULL)
        THEN @select_title::text
        ELSE title END AND
    text = CASE WHEN (@select_text::text <> '') AND (@select_text::text IS NOT NULL)
        THEN @select_text::text
        ELSE text END AND
    evaluation = CASE WHEN @select_evaluation::integer <> evaluation
        THEN @select_evaluation::integer
        ELSE evaluation END AND
    authors = CASE WHEN @select_authors::integer[] <> authors
        THEN @select_authors::integer[]
        ELSE authors END

ORDER BY
        -- Объединяем в группы столбцы одного типа (в один большой CASE WHEN нельзя разные типы)
        CASE WHEN @sorted_id_article::boolean THEN id_article::integer
             WHEN @sorted_evaluation::boolean THEN evaluation::integer END
        , -- запятая
        CASE WHEN @sorted_comments::boolean THEN comments::integer[]
             WHEN @sorted_authors::boolean THEN authors::integer[] END
        , -- запятая
        CASE WHEN @sorted_title::boolean THEN title::text
             WHEN @sorted_text::boolean THEN text::text END
        , -- запятая
        CASE WHEN @sorted_edited_at::boolean THEN edited_at::timestamp
             WHEN @sorted_created_at::boolean THEN created_at::timestamp END

LIMIT sqlc.arg('Limit')::integer
OFFSET sqlc.arg('Offset')::integer;

-- EditArticle Изменяем параметр(ы) статьи
-- name: EditArticle :one
UPDATE articles
SET
  -- Если изменили текст или заголовок то обновляем время изменения
  edited_at = CASE WHEN (@text::text <> '') THEN NOW()
                   WHEN (@title::text <> '') THEN NOW()
                   ELSE edited_at END,

  -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
  -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
  -- (Кстати, "::text" <- эти штуки нужны чтобы вместа pgtype был string/int32)
  -- CASE WHEN используется когда нельзя указать нулевое значение (пустую строку), COALESCE когда можно
  title = CASE WHEN @title::text <> '' THEN @title::text ELSE title END,
  text = CASE WHEN @text::text <> '' THEN @text::text ELSE text END,
  evaluation = CASE WHEN @evaluation::integer <> evaluation THEN @evaluation::integer ELSE evaluation END,
  comments = COALESCE(@comments::integer[], comments),
  authors = CASE WHEN cardinality(@authors::integer[]) <> 0 THEN @authors ELSE authors END
WHERE id_article = @id_article::integer
RETURNING *;

-- DeleteArticle Удаляем статью и комментарии к ней
-- name: DeleteArticle :one
WITH deleted_comments AS ( -- Объединяем 2 запроса в 1
    -- TODO: переделать в транзакцию
    DELETE FROM comments
    WHERE id_comment = ANY ((SELECT comments FROM articles
                            WHERE id_article = @id_article::integer)::text::integer[])
), update_id AS ( -- Объединяем 3 запроса в 1
    UPDATE articles
    SET id_article = id_article - 1
    WHERE id_article > @id_article::integer
)
DELETE FROM articles
WHERE id_article = @id_article::integer
RETURNING *;

-- CountRowsComment Считаем количество строк в таблице
-- name: CountRowsComment :one
SELECT COUNT(*) FROM comments;
