package main

import (
	"github.com/ZenSam7/Education/api"
	"github.com/ZenSam7/Education/db/sqlc"
	"github.com/ZenSam7/Education/tools"
	"log"
)

func main() {
	config := tools.LoadConfig(".")
	queries, closeConn := db.GetQueries()
	defer closeConn()

	server := api.NewProcess(queries)

	if err := server.Run(config.ServerAddress); err != nil {
		log.Fatal("Не получилось поднять сервер (api):", err)
	}
}
