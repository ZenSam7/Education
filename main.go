package main

import (
	db "github.com/ZenSam7/Education/db/sqlc"
)

func main() {
	_, closeConn := db.GetQueries()
	defer closeConn()
}
