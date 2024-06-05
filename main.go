package main

import (
	"github.com/ZenSam7/Education/api"
	db "github.com/ZenSam7/Education/db/sqlc"
	"github.com/ZenSam7/Education/tools"
	"log"
)

func main() {
	config := tools.LoadConfig(".")
	queries, closeConn := db.GetQueries()
	defer closeConn()

	server, err := api.NewProcess(config, queries)
	if err != nil {
		log.Fatal("Ошибка в создании роутера:", err.Error())
	}

	if err := server.Run(config.ServerAddress); err != nil {
		log.Fatal("Не получилось поднять сервер (api):", err)
	}
}
