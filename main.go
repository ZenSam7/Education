package main

import (
	"github.com/ZenSam7/Education/api"
	"github.com/ZenSam7/Education/db/sqlc"
	"log"
)

const (
	serverAdrress = "0.0.0.0:8080"
)

func main() {
	queries, closeConn := db.GetQueries()
	defer closeConn()

	server := api.NewProcess(queries)

	if err := server.Run(serverAdrress); err != nil {
		log.Fatal("Не получилось поднять сервер (api):", err)
	}
}
