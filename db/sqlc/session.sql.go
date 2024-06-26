// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: session.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const blockSession = `-- name: BlockSession :one
UPDATE sessions
SET blocked = true
WHERE id_session = $1::uuid
RETURNING id_session, issued_at, expired_at, refresh_token, id_user, client_ip, blocked
`

// BlockSession Блокируем сессию по id
func (q *Queries) BlockSession(ctx context.Context, idSession pgtype.UUID) (Session, error) {
	row := q.db.QueryRow(ctx, blockSession, idSession)
	var i Session
	err := row.Scan(
		&i.IDSession,
		&i.IssuedAt,
		&i.ExpiredAt,
		&i.RefreshToken,
		&i.IDUser,
		&i.ClientIp,
		&i.Blocked,
	)
	return i, err
}

const createSession = `-- name: CreateSession :one
INSERT INTO sessions (expired_at, refresh_token, id_user, id_session, client_ip)
VALUES ($1::timestamptz,
        $2::text,
        $3::integer,
        $4::uuid,
        $5::text)
RETURNING id_session, issued_at, expired_at, refresh_token, id_user, client_ip, blocked
`

type CreateSessionParams struct {
	ExpiredAt    pgtype.Timestamptz `json:"expired_at"`
	RefreshToken string             `json:"refresh_token"`
	IDUser       int32              `json:"id_user"`
	IDSession    pgtype.UUID        `json:"id_session"`
	ClientIp     string             `json:"client_ip"`
}

// CreateSession Создаём сессию
func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRow(ctx, createSession,
		arg.ExpiredAt,
		arg.RefreshToken,
		arg.IDUser,
		arg.IDSession,
		arg.ClientIp,
	)
	var i Session
	err := row.Scan(
		&i.IDSession,
		&i.IssuedAt,
		&i.ExpiredAt,
		&i.RefreshToken,
		&i.IDUser,
		&i.ClientIp,
		&i.Blocked,
	)
	return i, err
}

const deleteSession = `-- name: DeleteSession :one
DELETE FROM sessions
WHERE id_session = $1::uuid
RETURNING id_session, issued_at, expired_at, refresh_token, id_user, client_ip, blocked
`

// DeleteSession Удаляем сессию по id
func (q *Queries) DeleteSession(ctx context.Context, idSession pgtype.UUID) (Session, error) {
	row := q.db.QueryRow(ctx, deleteSession, idSession)
	var i Session
	err := row.Scan(
		&i.IDSession,
		&i.IssuedAt,
		&i.ExpiredAt,
		&i.RefreshToken,
		&i.IDUser,
		&i.ClientIp,
		&i.Blocked,
	)
	return i, err
}

const getSession = `-- name: GetSession :one
SELECT id_session, issued_at, expired_at, refresh_token, id_user, client_ip, blocked FROM sessions
WHERE id_session = $1::uuid
`

// GetSession Получаем сессиб по id
func (q *Queries) GetSession(ctx context.Context, idSession pgtype.UUID) (Session, error) {
	row := q.db.QueryRow(ctx, getSession, idSession)
	var i Session
	err := row.Scan(
		&i.IDSession,
		&i.IssuedAt,
		&i.ExpiredAt,
		&i.RefreshToken,
		&i.IDUser,
		&i.ClientIp,
		&i.Blocked,
	)
	return i, err
}
