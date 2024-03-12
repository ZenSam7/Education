// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Article struct {
	IDArticle  int32              `json:"id_article"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
	EditedAt   pgtype.Timestamptz `json:"edited_at"`
	Title      string             `json:"title"`
	Text       string             `json:"text"`
	Comments   []int32            `json:"comments"`
	Authors    []int32            `json:"authors"`
	Evaluation int32              `json:"evaluation"`
}

type Comment struct {
	IDComment  int32              `json:"id_comment"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
	EditedAt   pgtype.Timestamptz `json:"edited_at"`
	Text       string             `json:"text"`
	FromUser   int32              `json:"from_user"`
	Evaluation int32              `json:"evaluation"`
}

type User struct {
	IDUser       int32              `json:"id_user"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	Name         string             `json:"name"`
	Description  pgtype.Text        `json:"description"`
	Karma        int32              `json:"karma"`
	Email        string             `json:"email"`
	PasswordHash string             `json:"password_hash"`
}
