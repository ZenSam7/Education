package db

import (
	"context"
	"fmt"
	"github.com/ZenSam7/Education/tools"
	"github.com/jackc/pgx/v5"
	"log"
)

// GetQueries Возвращаем переменую через которую можно отправить запросы,
// и функцию для закрытия соединения с бд
func GetQueries() (*Queries, func()) {
	ctx := context.Background()
	config := tools.LoadConfig("../..")

	// Создаём url для соединения через pgx
	dbConnectParameters := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DBHost,
		config.DBUserName,
		config.DBPassword,
		config.DBName,
		config.DBSSLMode,
	)

	// Открываем соединение при помощи pgx
	conn, err := pgx.Connect(
		ctx,
		dbConnectParameters,
	)
	if err != nil {
		log.Fatal("Не получается подключиться к бд:", err)
	}

	// Создаём переменную для отправки запросов
	queries := New(conn)

	// Создаём лямбда-замыкание для закрытия соединения
	closeConnect := func() {
		func(ctx context.Context, conn *pgx.Conn) {
			if err := conn.Close(ctx); err != nil {
				log.Fatal("Не получается закрыть соединение:", err)
			}
		}(ctx, conn)
	}

	return queries, closeConnect
}
