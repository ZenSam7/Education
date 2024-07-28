// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: article.sql

package db

import (
	"context"
)

const countRowsComment = `-- name: CountRowsComment :one
SELECT COUNT(*) FROM comments
`

// CountRowsComment Считаем количество строк в таблице
func (q *Queries) CountRowsComment(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countRowsComment)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createArticle = `-- name: CreateArticle :one
INSERT INTO articles (title, text, authors)
VALUES ($1, $2, $3)
RETURNING id_article, created_at, edited_at, title, text, comments, authors, evaluation
`

type CreateArticleParams struct {
	Title   string  `json:"title"`
	Text    string  `json:"text"`
	Authors []int32 `json:"authors"`
}

// CreateArticle Создаём статью
func (q *Queries) CreateArticle(ctx context.Context, arg CreateArticleParams) (Article, error) {
	row := q.db.QueryRow(ctx, createArticle, arg.Title, arg.Text, arg.Authors)
	var i Article
	err := row.Scan(
		&i.IDArticle,
		&i.CreatedAt,
		&i.EditedAt,
		&i.Title,
		&i.Text,
		&i.Comments,
		&i.Authors,
		&i.Evaluation,
	)
	return i, err
}

const deleteArticle = `-- name: DeleteArticle :one
WITH deleted_comments AS ( -- Объединяем 2 запроса в 1 -- TODO: переделать в транзакцию
    DELETE FROM comments
    WHERE id_comment = ANY ((SELECT comments FROM articles
                            WHERE id_article = $1::integer)::text::integer[])
), update_id AS ( -- Объединяем 3 запроса в 1
    UPDATE articles
    SET id_article = id_article - 1
    WHERE id_article > $1::integer
)
DELETE FROM articles
WHERE id_article = $1::integer
RETURNING id_article, created_at, edited_at, title, text, comments, authors, evaluation
`

// DeleteArticle Удаляем статью и комментарии к ней
func (q *Queries) DeleteArticle(ctx context.Context, idArticle int32) (Article, error) {
	row := q.db.QueryRow(ctx, deleteArticle, idArticle)
	var i Article
	err := row.Scan(
		&i.IDArticle,
		&i.CreatedAt,
		&i.EditedAt,
		&i.Title,
		&i.Text,
		&i.Comments,
		&i.Authors,
		&i.Evaluation,
	)
	return i, err
}

const editArticle = `-- name: EditArticle :one
UPDATE articles
SET
  -- Если изменили текст или заголовок то обновляем время изменения
  edited_at = CASE WHEN ($1::text <> '') THEN NOW()
                   WHEN ($2::text <> '') THEN NOW()
                   ELSE edited_at END,

  -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
  -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
  -- (Кстати, "::text" <- эти штуки нужны чтобы вместа pgtype был string/int32)
  -- CASE WHEN используется когда нельзя указать нулевое значение (пустую строку), COALESCE когда можно
  title = CASE WHEN $2::text <> '' THEN $2::text ELSE title END,
  text = CASE WHEN $1::text <> '' THEN $1::text ELSE text END,
  comments = COALESCE($3::integer[], comments),
  authors = CASE WHEN cardinality($4::integer[]) <> 0 THEN $4 ELSE authors END
WHERE id_article = $5::integer
RETURNING id_article, created_at, edited_at, title, text, comments, authors, evaluation
`

type EditArticleParams struct {
	Text      string  `json:"text"`
	Title     string  `json:"title"`
	Comments  []int32 `json:"comments"`
	Authors   []int32 `json:"authors"`
	IDArticle int32   `json:"id_article"`
}

// EditArticle Изменяем параметр(ы) статьи
func (q *Queries) EditArticle(ctx context.Context, arg EditArticleParams) (Article, error) {
	row := q.db.QueryRow(ctx, editArticle,
		arg.Text,
		arg.Title,
		arg.Comments,
		arg.Authors,
		arg.IDArticle,
	)
	var i Article
	err := row.Scan(
		&i.IDArticle,
		&i.CreatedAt,
		&i.EditedAt,
		&i.Title,
		&i.Text,
		&i.Comments,
		&i.Authors,
		&i.Evaluation,
	)
	return i, err
}

const getArticle = `-- name: GetArticle :one
SELECT id_article, created_at, edited_at, title, text, comments, authors, evaluation FROM articles
WHERE id_article = $1
`

// GetArticle Возвращаем статью по id
func (q *Queries) GetArticle(ctx context.Context, idArticle int32) (Article, error) {
	row := q.db.QueryRow(ctx, getArticle, idArticle)
	var i Article
	err := row.Scan(
		&i.IDArticle,
		&i.CreatedAt,
		&i.EditedAt,
		&i.Title,
		&i.Text,
		&i.Comments,
		&i.Authors,
		&i.Evaluation,
	)
	return i, err
}

const getArticlesWithAttribute = `-- name: GetArticlesWithAttribute :many
SELECT id_article, created_at, edited_at, title, text, comments, authors, evaluation FROM articles
WHERE
    -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
    -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
    -- (Кстати, "::text" <- эти штуки нужны чтобы вместа pgtype был string/int32)
    title = CASE WHEN ($1::text <> '') AND ($1::text IS NOT NULL)
        THEN $1::text
        ELSE title END AND
    text = CASE WHEN ($2::text <> '') AND ($2::text IS NOT NULL)
        THEN $2::text
        ELSE text END AND
    evaluation = CASE WHEN $3::integer <> evaluation
        THEN $3::integer
        ELSE evaluation END
LIMIT $5::integer
OFFSET $4::integer
`

type GetArticlesWithAttributeParams struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	Evaluation int32  `json:"evaluation"`
	Offset     int32  `json:"Offset"`
	Limit      int32  `json:"Limit"`
}

// GetArticlesWithAttribute Возвращаем много статей взятых по какому-то признаку(ам)
func (q *Queries) GetArticlesWithAttribute(ctx context.Context, arg GetArticlesWithAttributeParams) ([]Article, error) {
	rows, err := q.db.Query(ctx, getArticlesWithAttribute,
		arg.Title,
		arg.Text,
		arg.Evaluation,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Article{}
	for rows.Next() {
		var i Article
		if err := rows.Scan(
			&i.IDArticle,
			&i.CreatedAt,
			&i.EditedAt,
			&i.Title,
			&i.Text,
			&i.Comments,
			&i.Authors,
			&i.Evaluation,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getManySortedArticles = `-- name: GetManySortedArticles :many
SELECT id_article, created_at, edited_at, title, text, comments, authors, evaluation FROM articles
ORDER BY
        CASE WHEN $1::boolean THEN id_article::integer
             WHEN $2::boolean THEN evaluation::integer END
        , -- запятая
        CASE WHEN $3::boolean THEN comments::integer[]
             WHEN $4::boolean THEN authors::integer[] END
        , -- запятая
        CASE WHEN $5::boolean THEN title::text
             WHEN $6::boolean THEN text::text END
        , -- запятая
        CASE WHEN $7::boolean THEN edited_at::timestamp
             WHEN $8::boolean THEN created_at::timestamp END
LIMIT $10::integer
OFFSET $9::integer
`

type GetManySortedArticlesParams struct {
	IDArticle  bool  `json:"id_article"`
	Evaluation bool  `json:"evaluation"`
	Comments   bool  `json:"comments"`
	Authors    bool  `json:"authors"`
	Title      bool  `json:"title"`
	Text       bool  `json:"text"`
	EditedAt   bool  `json:"edited_at"`
	CreatedAt  bool  `json:"created_at"`
	Offset     int32 `json:"Offset"`
	Limit      int32 `json:"Limit"`
}

// GetManySortedArticles Возвращаем много отсортированных статей
func (q *Queries) GetManySortedArticles(ctx context.Context, arg GetManySortedArticlesParams) ([]Article, error) {
	rows, err := q.db.Query(ctx, getManySortedArticles,
		arg.IDArticle,
		arg.Evaluation,
		arg.Comments,
		arg.Authors,
		arg.Title,
		arg.Text,
		arg.EditedAt,
		arg.CreatedAt,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Article{}
	for rows.Next() {
		var i Article
		if err := rows.Scan(
			&i.IDArticle,
			&i.CreatedAt,
			&i.EditedAt,
			&i.Title,
			&i.Text,
			&i.Comments,
			&i.Authors,
			&i.Evaluation,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getManySortedArticlesWithAttribute = `-- name: GetManySortedArticlesWithAttribute :many
SELECT id_article, created_at, edited_at, title, text, comments, authors, evaluation FROM articles
WHERE
    -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
    -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
    -- (Кстати, "::text" <- эти штуки нужны чтобы вместа pgtype был string/int32)
    title = CASE WHEN ($1::text <> '') AND ($1::text IS NOT NULL)
        THEN $1::text
        ELSE title END AND
    text = CASE WHEN ($2::text <> '') AND ($2::text IS NOT NULL)
        THEN $2::text
        ELSE text END AND
    evaluation = CASE WHEN $3::integer <> evaluation
        THEN $3::integer
        ELSE evaluation END
ORDER BY
        -- Объединяем в группы столбцы одного типа (в один большой CASE WHEN нельзя разные типы)
        CASE WHEN $4::boolean THEN id_article::integer
             WHEN $5::boolean THEN evaluation::integer END
        , -- запятая
        CASE WHEN $6::boolean THEN comments::integer[]
             WHEN $7::boolean THEN authors::integer[] END
        , -- запятая
        CASE WHEN $8::boolean THEN title::text
             WHEN $9::boolean THEN text::text END
        , -- запятая
        CASE WHEN $10::boolean THEN edited_at::timestamp
             WHEN $11::boolean THEN created_at::timestamp END

LIMIT $13::integer
OFFSET $12::integer
`

type GetManySortedArticlesWithAttributeParams struct {
	SelectTitle      string `json:"select_title"`
	SelectText       string `json:"select_text"`
	SelectEvaluation int32  `json:"select_evaluation"`
	SortedIDArticle  bool   `json:"sorted_id_article"`
	SortedEvaluation bool   `json:"sorted_evaluation"`
	SortedComments   bool   `json:"sorted_comments"`
	SortedAuthors    bool   `json:"sorted_authors"`
	SortedTitle      bool   `json:"sorted_title"`
	SortedText       bool   `json:"sorted_text"`
	SortedEditedAt   bool   `json:"sorted_edited_at"`
	SortedCreatedAt  bool   `json:"sorted_created_at"`
	Offset           int32  `json:"Offset"`
	Limit            int32  `json:"Limit"`
}

// GetManySortedArticlesWithAttribute Возвращаем много статей взятых по признаку по
// какому-то признаку(ам) и отсортированных по другому признаку(ам)
func (q *Queries) GetManySortedArticlesWithAttribute(ctx context.Context, arg GetManySortedArticlesWithAttributeParams) ([]Article, error) {
	rows, err := q.db.Query(ctx, getManySortedArticlesWithAttribute,
		arg.SelectTitle,
		arg.SelectText,
		arg.SelectEvaluation,
		arg.SortedIDArticle,
		arg.SortedEvaluation,
		arg.SortedComments,
		arg.SortedAuthors,
		arg.SortedTitle,
		arg.SortedText,
		arg.SortedEditedAt,
		arg.SortedCreatedAt,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Article{}
	for rows.Next() {
		var i Article
		if err := rows.Scan(
			&i.IDArticle,
			&i.CreatedAt,
			&i.EditedAt,
			&i.Title,
			&i.Text,
			&i.Comments,
			&i.Authors,
			&i.Evaluation,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
