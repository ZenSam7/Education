package main

import (
	"context"
	"github.com/ZenSam7/Education/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

var DataBase db.DBTX

func main() {
	query := db.Queries{Db: DataBase}

	ctx := context.Background()

	newUserInfo := db.CreateUserParams{
		Name:        "ZenSam7",
		Description: pgtype.Text{String: "pass", Valid: true},
		Email:       pgtype.Text{String: "abc@abc.abc", Valid: true},
	}

	if _, err := query.CreateUser(ctx, newUserInfo); err != nil {
		return
	}

}
