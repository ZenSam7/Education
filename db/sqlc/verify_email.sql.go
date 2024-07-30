// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: verify_email.sql

package db

import (
	"context"
)

const createVerifyRequest = `-- name: CreateVerifyRequest :one
INSERT INTO verify_emails (id_user, secret_key)
VALUES ($1::integer, $2::text)
RETURNING id_verify_email, id_user, secret_key, expired_at
`

type CreateVerifyRequestParams struct {
	IDUser    int32  `json:"id_user"`
	SecretKey string `json:"secret_key"`
}

// CreateVerifyRequest Создаём новый запрос на верификацию почты
func (q *Queries) CreateVerifyRequest(ctx context.Context, arg CreateVerifyRequestParams) (VerifyEmail, error) {
	row := q.db.QueryRow(ctx, createVerifyRequest, arg.IDUser, arg.SecretKey)
	var i VerifyEmail
	err := row.Scan(
		&i.IDVerifyEmail,
		&i.IDUser,
		&i.SecretKey,
		&i.ExpiredAt,
	)
	return i, err
}

const deleteVerifyRequest = `-- name: DeleteVerifyRequest :one
DELETE FROM verify_emails
WHERE id_user = $1::integer
RETURNING id_verify_email, id_user, secret_key, expired_at
`

// DeleteVerifyRequest Удаляем запрос на верификацию
func (q *Queries) DeleteVerifyRequest(ctx context.Context, idUser int32) (VerifyEmail, error) {
	row := q.db.QueryRow(ctx, deleteVerifyRequest, idUser)
	var i VerifyEmail
	err := row.Scan(
		&i.IDVerifyEmail,
		&i.IDUser,
		&i.SecretKey,
		&i.ExpiredAt,
	)
	return i, err
}

const getVerifyRequest = `-- name: GetVerifyRequest :one
SELECT id_verify_email, id_user, secret_key, expired_at FROM verify_emails
WHERE id_user = $1::integer
`

// GetVerifyRequest Возвращаем запрос на верификацию
func (q *Queries) GetVerifyRequest(ctx context.Context, idUser int32) (VerifyEmail, error) {
	row := q.db.QueryRow(ctx, getVerifyRequest, idUser)
	var i VerifyEmail
	err := row.Scan(
		&i.IDVerifyEmail,
		&i.IDUser,
		&i.SecretKey,
		&i.ExpiredAt,
	)
	return i, err
}