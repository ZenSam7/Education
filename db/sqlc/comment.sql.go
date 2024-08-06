// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: comment.sql

package db

import (
	"context"
)

const countRowsArticle = `-- name: CountRowsArticle :one
SELECT COUNT(*) FROM articles
`

// CountRowsArticle Считаем количество строк в таблице
func (q *Queries) CountRowsArticle(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countRowsArticle)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createComment = `-- name: CreateComment :one
WITH add_comment AS ( -- TODO: переделать в транзакцию
    UPDATE articles
    SET comments = array_append(comments, lastval())
    WHERE id_article = $1
)
INSERT INTO comments (text, author)
VALUES ($3, $2)
RETURNING id_comment, created_at, edited_at, text, author, evaluation
`

type CreateCommentParams struct {
	IDArticle int32  `json:"id_article"`
	Author    int32  `json:"author"`
	Text      string `json:"text"`
}

// CreateComment Создаём комментарий к статье
func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error) {
	row := q.db.QueryRow(ctx, createComment, arg.IDArticle, arg.Author, arg.Text)
	var i Comment
	err := row.Scan(
		&i.IDComment,
		&i.CreatedAt,
		&i.EditedAt,
		&i.Text,
		&i.Author,
		&i.Evaluation,
	)
	return i, err
}

const deleteComment = `-- name: DeleteComment :one
WITH deleted_comment_id AS ( -- TODO: переделать в транзакцию
    DELETE FROM comments
    WHERE id_comment = $1::integer
)
UPDATE articles
SET comments = array_remove(comments, $1::integer)
WHERE $1::integer = ANY(comments)
RETURNING id_article, created_at, edited_at, title, text, comments, authors, evaluation
`

// DeleteComment Удаляем комментарий к статье
func (q *Queries) DeleteComment(ctx context.Context, idComment int32) (Article, error) {
	row := q.db.QueryRow(ctx, deleteComment, idComment)
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

const editComment = `-- name: EditComment :one
UPDATE comments
SET
    -- Если изменили текст или автора пользователя, то обновляем его
    edited_at = CASE WHEN $1::text <> '' THEN NOW()
                     WHEN $2::integer <> author THEN NOW()
                     ELSE edited_at END,

    -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
    -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
    text = CASE WHEN $1::text <> '' THEN $1::text ELSE text END,
    evaluation = CASE WHEN $3::integer <> evaluation THEN $3::integer ELSE evaluation END
WHERE id_comment = $4::integer
RETURNING id_comment, created_at, edited_at, text, author, evaluation
`

type EditCommentParams struct {
	Text       string `json:"text"`
	Author     int32  `json:"author"`
	Evaluation int32  `json:"evaluation"`
	IDComment  int32  `json:"id_comment"`
}

// EditComment Изменяем параметр(ы) пользователя
func (q *Queries) EditComment(ctx context.Context, arg EditCommentParams) (Comment, error) {
	row := q.db.QueryRow(ctx, editComment,
		arg.Text,
		arg.Author,
		arg.Evaluation,
		arg.IDComment,
	)
	var i Comment
	err := row.Scan(
		&i.IDComment,
		&i.CreatedAt,
		&i.EditedAt,
		&i.Text,
		&i.Author,
		&i.Evaluation,
	)
	return i, err
}

const getComment = `-- name: GetComment :one
SELECT id_comment, created_at, edited_at, text, author, evaluation FROM comments
WHERE id_comment = $1::integer
`

// GetComment Возвращаем комментарий
func (q *Queries) GetComment(ctx context.Context, idComment int32) (Comment, error) {
	row := q.db.QueryRow(ctx, getComment, idComment)
	var i Comment
	err := row.Scan(
		&i.IDComment,
		&i.CreatedAt,
		&i.EditedAt,
		&i.Text,
		&i.Author,
		&i.Evaluation,
	)
	return i, err
}

const getCommentsOfArticle = `-- name: GetCommentsOfArticle :many
WITH the_article AS ( -- TODO: переделать в транзакцию
    SELECT unnest(comments) AS id_comment FROM articles
    WHERE id_article = $3::integer
)
SELECT id_comment, created_at, edited_at, text, author, evaluation FROM comments
WHERE id_comment IN (SELECT id_comment FROM the_article)
OFFSET $1::integer
LIMIT $2::integer
`

type GetCommentsOfArticleParams struct {
	Offset    int32 `json:"Offset"`
	Limit     int32 `json:"Limit"`
	IDArticle int32 `json:"id_article"`
}

// GetCommentsOfArticle Возвращаем комментарии к статье
func (q *Queries) GetCommentsOfArticle(ctx context.Context, arg GetCommentsOfArticleParams) ([]Comment, error) {
	rows, err := q.db.Query(ctx, getCommentsOfArticle, arg.Offset, arg.Limit, arg.IDArticle)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Comment{}
	for rows.Next() {
		var i Comment
		if err := rows.Scan(
			&i.IDComment,
			&i.CreatedAt,
			&i.EditedAt,
			&i.Text,
			&i.Author,
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
