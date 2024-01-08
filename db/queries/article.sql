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
    (NULLIF(edited_at, COALESCE(sqlc.arg(edited_at)::timestamp , edited_at)) IS NULL) AND

    -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
    -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
    (title = (CASE WHEN (sqlc.arg(title)::text <> '' AND sqlc.arg(title) <> NULL)
        THEN sqlc.arg(title)::text
        ELSE title END)) AND
    (text = (CASE WHEN (sqlc.arg(text)::text <> '' AND sqlc.arg(title) <> NULL)
        THEN sqlc.arg(text)::text
        ELSE text END)) AND

    (NULLIF(comments, COALESCE(sqlc.arg(comments), comments)) IS NULL) AND
    (NULLIF(authors, COALESCE(sqlc.arg(authors), authors)) IS NULL) AND
    (NULLIF(evaluation, COALESCE(sqlc.arg(evaluation), evaluation)) IS NULL)
LIMIT sqlc.arg('Limit')::integer
OFFSET sqlc.arg('Offset')::integer;

-- GetManySortedArticles Возвращаем много статей отсортированных по признаку sorted_at
-- name: GetManySortedArticles :many
SELECT * FROM articles
ORDER BY sqlc.arg(sorted_at)::text
LIMIT sqlc.arg('Limit')::integer
OFFSET sqlc.arg('Offset')::integer;

-- GetManySortedArticlesWithAttribute Возвращаем много статей взятых по признаку по
-- какому-то признаку(ам) отсортированных по признаку sortedAt
-- name: GetManySortedArticlesWithAttribute :many
SELECT * FROM articles
WHERE
    edited_at = COALESCE(sqlc.arg(edited_at)::timestamp , edited_at) AND

    -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
    -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
    title = CASE WHEN sqlc.arg(title)::text <> '' THEN sqlc.arg(title)::text ELSE title END AND
    text = CASE WHEN sqlc.arg(text)::text <> '' THEN sqlc.arg(text)::text ELSE text END     AND

    comments = COALESCE(sqlc.arg(comments), comments)                AND
    authors = COALESCE(sqlc.arg(authors), authors)                   AND
    evaluation = COALESCE(sqlc.arg(evaluation), evaluation)
ORDER BY sqlc.arg(sorted_at)::text
LIMIT sqlc.arg('Limit')::integer
OFFSET sqlc.arg('Offset')::integer;

-- EditArticleParam Изменяем параметр(ы) статьи
-- name: EditArticleParam :one
UPDATE articles
SET
  edited_at = COALESCE(sqlc.arg(edited_at)::timestamp , edited_at),
  -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
  -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
  -- (Кстати, "::text" <- эти штуки нужны чтобы вместа pgtype был string/int32)
  title = CASE WHEN sqlc.arg(title)::text <> '' THEN sqlc.arg(title)::text ELSE title END,
  text = CASE WHEN sqlc.arg(text)::text <> '' THEN sqlc.arg(text)::text ELSE text END,
  comments = COALESCE(sqlc.arg(comments), comments),
  authors = COALESCE(sqlc.arg(authors), authors),
  evaluation = COALESCE(sqlc.arg(evaluation), evaluation)
WHERE id_article = sqlc.arg(id_article)
RETURNING *;
