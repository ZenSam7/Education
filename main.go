package main

import (
	"github.com/ZenSam7/Education/db/sqlc"
)

func main() {
	_, c := db.GetQueries()
	defer c()
}
