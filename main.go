package main

import (
	"context"
	"github.com/ZenSam7/Education/db/sqlc"
)

func main() {
	query := db.New(?)

	ctx := context.Background()

	newUserInfo := db.CreateUserParams{
		Name:        "ZenSam7",
		Description: "pass",
		Email:       "abc@abc.abc",
	}

	if err := query.CreateUser(ctx, newUserInfo); err != nil {
		return
	}
}
