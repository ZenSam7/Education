// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: comment.sql

package db

import (
	"context"
)

const createComment = `-- name: CreateComment :one
WITH add_comment AS ( -- Я знаю что есть тразнакции
    UPDATE articles
    SET comments = array_append(comments, lastval())
    WHERE id_article = $1
)
INSERT INTO comments (text, from_user)
VALUES ($3, $2)
RETURNING id_comment, created_at, edited_at, text, from_user, evaluation
`

type CreateCommentParams struct {
	IDArticle int32  `json:"id_article"`
	FromUser  int32  `json:"from_user"`
	Text      string `json:"text"`
}

// CreateComment Создаём комментарий к статье
func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error) {
	row := q.db.QueryRow(ctx, createComment, arg.IDArticle, arg.FromUser, arg.Text)
	var i Comment
	err := row.Scan(
		&i.IDComment,
		&i.CreatedAt,
		&i.EditedAt,
		&i.Text,
		&i.FromUser,
		&i.Evaluation,
	)
	return i, err
}

const deleteComment = `-- name: DeleteComment :one
WITH deleted_comment_id AS ( -- Я знаю что есть тразнакции
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

const editCommentParam = `-- name: EditCommentParam :one
UPDATE comments
SET
    -- Если изменили текст или автора польщователя то обновляем его
    edited_at = CASE WHEN $1::text <> '' THEN NOW()
                     WHEN $2::integer <> from_user THEN NOW()
                     ELSE edited_at END,

    -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
    -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
    text = CASE WHEN $1::text <> '' THEN $1::text ELSE text END,
    from_user = CASE WHEN $2::integer <> from_user
                     THEN $2::integer
                     ELSE from_user END,
    evaluation = CASE WHEN $3::integer <> evaluation
                      THEN $3::integer
                      ELSE evaluation END
WHERE id_comment = $4::integer
RETURNING id_comment, created_at, edited_at, text, from_user, evaluation
`

type EditCommentParamParams struct {
	Text       string `json:"text"`
	FromUser   int32  `json:"from_user"`
	Evaluation int32  `json:"evaluation"`
	IDComment  int32  `json:"id_comment"`
}

// EditCommentParam Изменяем параметр(ы) пользователя
func (q *Queries) EditCommentParam(ctx context.Context, arg EditCommentParamParams) (Comment, error) {
	row := q.db.QueryRow(ctx, editCommentParam,
		arg.Text,
		arg.FromUser,
		arg.Evaluation,
		arg.IDComment,
	)
	var i Comment
	err := row.Scan(
		&i.IDComment,
		&i.CreatedAt,
		&i.EditedAt,
		&i.Text,
		&i.FromUser,
		&i.Evaluation,
	)
	return i, err
}

const getComment = `-- name: GetComment :one
SELECT id_comment, created_at, edited_at, text, from_user, evaluation FROM comments
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
		&i.FromUser,
		&i.Evaluation,
	)
	return i, err
}
