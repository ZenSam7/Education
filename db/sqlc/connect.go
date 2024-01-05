package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
)

const connString = "host=localhost user=root password=root dbname=education sslmode=disable"

// GetQueries Возвращаем переменую через которую можно отправить запросы,
// и функцию для закрытия соединения с бд
func GetQueries() (*Queries, func()) {
	ctx := context.Background()

	// Открываем соединение при помощи pgx
	conn, err := pgx.Connect(
		ctx,
		connString,
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
