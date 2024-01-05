-- CreateArticle Создаём статью
-- name: CreateArticle :exec
INSERT INTO articles (title, text, authors)
VALUES ($1, $2, $3)
RETURNING *;

-- GetArticle Возвращаем статью по id
-- name: GetArticle :one
SELECT * FROM articles
WHERE id_article = $1;

-- GetArticlesWithAttribute Возвращаем много статей взятых по признаку attribute
-- name: GetArticlesWithAttribute :many
SELECT * FROM articles
WHERE sqlc.arg(attribute)::text = sqlc.arg(attribute_value)::text
LIMIT $1
OFFSET $2;

-- GetManySortedArticles Возвращаем много статей отсортированных по признаку attribute
-- name: GetManySortedArticles :many
SELECT * FROM articles
ORDER BY sqlc.arg(sorted_at)::text
LIMIT $1
OFFSET $2;

-- GetManySortedArticlesWithAttribute Возвращаем много статей взятых по признаку attridute отсортированных по признаку sortedAt
-- name: GetManySortedArticlesWithAttribute :many
SELECT * FROM articles
WHERE sqlc.arg(attribute)::text = sqlc.arg(attribute_value)::text
ORDER BY sqlc.arg(sorted_at)::text
LIMIT $1
OFFSET $2;

-- EditArticleParam Изменяем параметр(ы) статьи
-- name: EditArticleParam :one
UPDATE articles
SET
  edited_at = COALESCE(sqlc.arg(edited_at)::timestamp , edited_at),
  title = COALESCE(sqlc.arg(title)::text, title),
  text = COALESCE(sqlc.arg(text), text),
  comments = COALESCE(sqlc.arg(comments), comments),
  authors = COALESCE(sqlc.arg(authors), authors),
  evaluation = COALESCE(sqlc.arg(evaluation), evaluation)
WHERE id_article = sqlc.arg(id_article)
RETURNING *;
