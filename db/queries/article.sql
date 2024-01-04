-- CreateArticle Создаём статью
-- name: CreateArticle :exec
INSERT INTO articles (title, text, authors)
VALUES ($1, $2, $3);

-- GetArticle Возвращаем статью по id
-- name: GetArticle :one
SELECT * FROM articles
WHERE id_article = $1;

-- GetArticlesWithTitle Возвращаем много статей взятых по названию
-- name: GetArticlesWithTitle :many
SELECT * FROM articles
WHERE title = $1
LIMIT $2
OFFSET $3;

-- GetArticlesWithEvalution Возвращаем много статей взятых по оценке
-- name: GetArticlesWithEvalution :many
SELECT * FROM articles
WHERE evaluation = $1
LIMIT $2
OFFSET $3;

-- GetManySortedArticles Возвращаем много статей отсортированных по признаку Column1
-- name: GetManySortedArticles :many
SELECT * FROM articles
ORDER BY $1::text
LIMIT $2
OFFSET $3;


-- EditArticleText Изменяем текст статьи и обновляем время изменения статьи (Column1 = id_article, Column2 = text)
-- name: EditArticleText :exec
WITH update_time AS (
    UPDATE articles
    SET edited_at = now()
    WHERE id_article = $1::integer
)
UPDATE articles
SET text = $2::text
WHERE id_article = $1::integer;

-- EditArticleTitle Изменяем заголовок статьи
-- name: EditArticleTitle :exec
UPDATE articles
SET title = $2
WHERE id_article = $1;

-- AddArticleAuthors Добавляем автора статьи
-- name: AddArticleAuthors :exec
UPDATE articles
SET authors = array_append(authors, $2)
WHERE id_article = $1;

-- DeleteArticleAuthors Удаляем автора статьи
-- name: DeleteArticleAuthors :exec
UPDATE articles
SET authors = array_remove(authors, $2)
WHERE id_article = $1;
