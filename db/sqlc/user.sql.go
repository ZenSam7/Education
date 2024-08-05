// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countRowsUser = `-- name: CountRowsUser :one
SELECT COUNT(*) FROM users
`

// CountRowsUser Считаем количество строк в таблице
func (q *Queries) CountRowsUser(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countRowsUser)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (name, email, password_hash)
VALUES ($1::text, $2::text, $3::text)
RETURNING id_user, created_at, name, description, karma, email, password_hash, email_verified, role
`

type CreateUserParams struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

// CreateUser Создаём пользователя
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Name, arg.Email, arg.PasswordHash)
	var i User
	err := row.Scan(
		&i.IDUser,
		&i.CreatedAt,
		&i.Name,
		&i.Description,
		&i.Karma,
		&i.Email,
		&i.PasswordHash,
		&i.EmailVerified,
		&i.Role,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :one
WITH deleted_session AS ( -- Объединяем 2 запроса в 1
    DELETE FROM sessions
    WHERE id_user = $1::integer
), delete_verify_request AS ( -- Объединяем 3 запроса в 1
    DELETE FROM verify_emails
    WHERE id_user = $1::integer
)
DELETE FROM users
WHERE id_user = $1::integer
RETURNING id_user, created_at, name, description, karma, email, password_hash, email_verified, role
`

// DeleteUser Удаляем пользователя
func (q *Queries) DeleteUser(ctx context.Context, idUser int32) (User, error) {
	row := q.db.QueryRow(ctx, deleteUser, idUser)
	var i User
	err := row.Scan(
		&i.IDUser,
		&i.CreatedAt,
		&i.Name,
		&i.Description,
		&i.Karma,
		&i.Email,
		&i.PasswordHash,
		&i.EmailVerified,
		&i.Role,
	)
	return i, err
}

const editUser = `-- name: EditUser :one
UPDATE users
SET
  -- Крч если через go передать в качестве текстового аргумента nil то он замениться на '',
  -- а '' != NULL поэтому она вставиться как пустая строка, хотя в go мы передали nil
  -- CASE WHEN используется когда нельзя указать нулевое значение (пустую строку), COALESCE когда можно
  name = CASE WHEN $1::text <> '' THEN $1::text ELSE name END,
  description = COALESCE($2::text, description),
  karma = COALESCE($3::integer, karma)
WHERE id_user = $4::integer
RETURNING id_user, created_at, name, description, karma, email, password_hash, email_verified, role
`

type EditUserParams struct {
	Name        string      `json:"name"`
	Description pgtype.Text `json:"description"`
	Karma       pgtype.Int4 `json:"karma"`
	IDUser      int32       `json:"id_user"`
}

// EditUser Изменяем параметр(ы) пользователя
func (q *Queries) EditUser(ctx context.Context, arg EditUserParams) (User, error) {
	row := q.db.QueryRow(ctx, editUser,
		arg.Name,
		arg.Description,
		arg.Karma,
		arg.IDUser,
	)
	var i User
	err := row.Scan(
		&i.IDUser,
		&i.CreatedAt,
		&i.Name,
		&i.Description,
		&i.Karma,
		&i.Email,
		&i.PasswordHash,
		&i.EmailVerified,
		&i.Role,
	)
	return i, err
}

const getManySortedUsers = `-- name: GetManySortedUsers :many
SELECT id_user, created_at, name, description, karma, email, password_hash, email_verified, role FROM users
ORDER BY
        CASE WHEN $1::boolean THEN id_user::integer
             WHEN $2::boolean THEN karma::integer END
        , -- запятая
        CASE WHEN $3::boolean THEN name::text
             WHEN $4::boolean THEN description::text END
LIMIT $6::integer
OFFSET $5::integer
`

type GetManySortedUsersParams struct {
	IDUser      bool  `json:"id_user"`
	Karma       bool  `json:"karma"`
	Name        bool  `json:"name"`
	Description bool  `json:"description"`
	Offset      int32 `json:"Offset"`
	Limit       int32 `json:"Limit"`
}

// GetManySortedUsers Возвращаем слайс пользователей отсортированных по какому-то параметру
// (можно поставить: id_user, и сортировки не будет)
func (q *Queries) GetManySortedUsers(ctx context.Context, arg GetManySortedUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, getManySortedUsers,
		arg.IDUser,
		arg.Karma,
		arg.Name,
		arg.Description,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.IDUser,
			&i.CreatedAt,
			&i.Name,
			&i.Description,
			&i.Karma,
			&i.Email,
			&i.PasswordHash,
			&i.EmailVerified,
			&i.Role,
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

const getUser = `-- name: GetUser :one
SELECT id_user, created_at, name, description, karma, email, password_hash, email_verified, role FROM users
WHERE id_user = $1
`

// GetUser Возвращаем пользователя
func (q *Queries) GetUser(ctx context.Context, idUser int32) (User, error) {
	row := q.db.QueryRow(ctx, getUser, idUser)
	var i User
	err := row.Scan(
		&i.IDUser,
		&i.CreatedAt,
		&i.Name,
		&i.Description,
		&i.Karma,
		&i.Email,
		&i.PasswordHash,
		&i.EmailVerified,
		&i.Role,
	)
	return i, err
}

const getUserFromName = `-- name: GetUserFromName :one
SELECT id_user, created_at, name, description, karma, email, password_hash, email_verified, role FROM users
WHERE name = $1
`

// GetUserFromName Возвращаем пользователя по имени
func (q *Queries) GetUserFromName(ctx context.Context, name string) (User, error) {
	row := q.db.QueryRow(ctx, getUserFromName, name)
	var i User
	err := row.Scan(
		&i.IDUser,
		&i.CreatedAt,
		&i.Name,
		&i.Description,
		&i.Karma,
		&i.Email,
		&i.PasswordHash,
		&i.EmailVerified,
		&i.Role,
	)
	return i, err
}

const setEmailIsVerified = `-- name: SetEmailIsVerified :one
UPDATE users
SET email_verified = true
WHERE id_user = $1::integer
RETURNING id_user, created_at, name, description, karma, email, password_hash, email_verified, role
`

// SetEmailIsVerified Ставим состояние почты как подтверждённую для какого-то пользователя
func (q *Queries) SetEmailIsVerified(ctx context.Context, idUser int32) (User, error) {
	row := q.db.QueryRow(ctx, setEmailIsVerified, idUser)
	var i User
	err := row.Scan(
		&i.IDUser,
		&i.CreatedAt,
		&i.Name,
		&i.Description,
		&i.Karma,
		&i.Email,
		&i.PasswordHash,
		&i.EmailVerified,
		&i.Role,
	)
	return i, err
}
