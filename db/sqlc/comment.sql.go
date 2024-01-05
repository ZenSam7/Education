// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: comment.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createComment = `-- name: CreateComment :exec
WITH add_comment AS (
    INSERT INTO comments (text, from_user)
    VALUES ($3, $2)
)
UPDATE articles
SET comments = array_append(comments, currval(pg_get_serial_sequence('comments','id_comment')))
WHERE id_article = $1
`

type CreateCommentParams struct {
	IDArticle int32
	FromUser  int32
	Text      string
}

// CreateComment Создаём комментарий к статье
func (q *Queries) CreateComment(ctx context.Context, arg CreateCommentParams) error {
	_, err := q.db.Exec(ctx, createComment, arg.IDArticle, arg.FromUser, arg.Text)
	return err
}

const deleteComment = `-- name: DeleteComment :exec
WITH deleted_comment_id AS (
    DELETE FROM comments
    WHERE id_comment = $2::integer
)
UPDATE articles
SET comments = array_remove(comments, $2::integer)
WHERE id_article = $1
`

type DeleteCommentParams struct {
	IDArticle int32
	IDComment int32
}

// DeleteComment Удаляем комментарий к статье
func (q *Queries) DeleteComment(ctx context.Context, arg DeleteCommentParams) error {
	_, err := q.db.Exec(ctx, deleteComment, arg.IDArticle, arg.IDComment)
	return err
}

const editCommentParam = `-- name: EditCommentParam :one
UPDATE comments
SET
  edited_at = COALESCE($1::timestamp, edited_at),
  text = COALESCE($2::text, text),
  from_user = COALESCE($3::integer, from_user),
  evaluation = COALESCE($4::integer, evaluation)
WHERE id_comment = $5::integer
RETURNING id_comment, created_at, edited_at, text, from_user, evaluation
`

type EditCommentParamParams struct {
	EditedAt   pgtype.Timestamp
	Text       string
	FromUser   int32
	Evaluation int32
	IDComment  int32
}

// EditCommentParam Изменяем параметр(ы) пользователя
func (q *Queries) EditCommentParam(ctx context.Context, arg EditCommentParamParams) (Comment, error) {
	row := q.db.QueryRow(ctx, editCommentParam,
		arg.EditedAt,
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

const getComment = `-- name: GetComment :exec
SELECT id_comment, created_at, edited_at, text, from_user, evaluation FROM comments
WHERE id_comment = $1
`

// GetComment Возвращаем комментарий
func (q *Queries) GetComment(ctx context.Context, idComment int32) error {
	_, err := q.db.Exec(ctx, getComment, idComment)
	return err
}
