// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: image.sql

package db

import (
	"context"
)

const deleteImage = `-- name: DeleteImage :one
DELETE FROM images
WHERE id_image = $1::integer
RETURNING id_image, name, content, id_user
`

// DeleteImage Удаляем картинку
func (q *Queries) DeleteImage(ctx context.Context, idImage int32) (Image, error) {
	row := q.db.QueryRow(ctx, deleteImage, idImage)
	var i Image
	err := row.Scan(
		&i.IDImage,
		&i.Name,
		&i.Content,
		&i.IDUser,
	)
	return i, err
}

const editImage = `-- name: EditImage :one
UPDATE images
SET content = $1::bytea
WHERE id_image = $2::integer
RETURNING id_image, name, content, id_user
`

type EditImageParams struct {
	Content []byte `json:"content"`
	IDImage int32  `json:"id_image"`
}

// EditImage Заменяем картинку на новую
func (q *Queries) EditImage(ctx context.Context, arg EditImageParams) (Image, error) {
	row := q.db.QueryRow(ctx, editImage, arg.Content, arg.IDImage)
	var i Image
	err := row.Scan(
		&i.IDImage,
		&i.Name,
		&i.Content,
		&i.IDUser,
	)
	return i, err
}

const getImage = `-- name: GetImage :one
SELECT id_image, name, content, id_user FROM images
WHERE id_image = $1::integer
`

// GetImage Возвращаем картинку
func (q *Queries) GetImage(ctx context.Context, idImage int32) (Image, error) {
	row := q.db.QueryRow(ctx, getImage, idImage)
	var i Image
	err := row.Scan(
		&i.IDImage,
		&i.Name,
		&i.Content,
		&i.IDUser,
	)
	return i, err
}

const loadImage = `-- name: LoadImage :one
INSERT INTO images (name, content, id_user)
VALUES ($1::text, $2::bytea, $3::integer)
RETURNING id_image, name, content, id_user
`

type LoadImageParams struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
	IDUser  int32  `json:"id_user"`
}

// LoadImage Загружаем изображение в бд
func (q *Queries) LoadImage(ctx context.Context, arg LoadImageParams) (Image, error) {
	row := q.db.QueryRow(ctx, loadImage, arg.Name, arg.Content, arg.IDUser)
	var i Image
	err := row.Scan(
		&i.IDImage,
		&i.Name,
		&i.Content,
		&i.IDUser,
	)
	return i, err
}

const renameImage = `-- name: RenameImage :one
UPDATE images
SET name = $1::text
WHERE id_image = $2::integer
RETURNING id_image, name, content, id_user
`

type RenameImageParams struct {
	Name    string `json:"name"`
	IDImage int32  `json:"id_image"`
}

// RenameImage Переименовываем картинку
func (q *Queries) RenameImage(ctx context.Context, arg RenameImageParams) (Image, error) {
	row := q.db.QueryRow(ctx, renameImage, arg.Name, arg.IDImage)
	var i Image
	err := row.Scan(
		&i.IDImage,
		&i.Name,
		&i.Content,
		&i.IDUser,
	)
	return i, err
}
