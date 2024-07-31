package db

import (
	"context"
	"fmt"
	"github.com/ZenSam7/Education/tools"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

// queries Можно было бы спокойно обойтись и без этого костыля, но он нужен для MakeTx
// (чтобы не создавать лишних коннектов к бд)
var queries *Queries

// MakeQueries Создаёт коннект к бд и функцию для закрытия этого соединения с бд
func MakeQueries() (*Queries, func()) {
	config := tools.LoadConfig()

	// Создаём url для соединения через pgx
	dbConnectParameters := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=education sslmode=%s",
		config.DBHost,
		config.DBUserName,
		config.DBPassword,
		config.DBSSLMode,
	)

	// Открываем соединение при помощи pgx
	conn, err := pgx.Connect(
		context.Background(),
		dbConnectParameters,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Не получается подключиться к бд")
	}

	// Создаём переменную для отправки запросов
	queries := New(conn)

	// Создаём лямбда-замыкание для закрытия соединения
	closeConnect := func() {
		func(ctx context.Context, conn *pgx.Conn) {
			if err := conn.Close(ctx); err != nil {
				log.Fatal().Err(err).Msg("Не получается закрыть соединение")
			} else {
				log.Info().Msg("Закрыли соединение")
			}
		}(context.Background(), conn)
	}

	return queries, closeConnect
}
