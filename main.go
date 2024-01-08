package main

import (
	"github.com/ZenSam7/Education/db/sqlc"
)

func main() {
	_, closeConn := db.GetQueries()
	defer closeConn()

	type a struct {
		a string
	}
}
