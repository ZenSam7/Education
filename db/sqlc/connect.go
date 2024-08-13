package db

import (
	"context"
	"fmt"
	"github.com/ZenSam7/Education/tools"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
	"time"
)

// queries Можно было бы спокойно обойтись и без этого костыля, но он нужен для MakeTx
// (чтобы не создавать лишних коннектов к бд)
var queries *Queries
var replica *Queries

// MakeQueries Создаёт коннект к бд и реплике и функцию для закрытия этого соединений с ними
func MakeQueries() (*Queries, *Queries, func()) {
	config := tools.LoadConfig()
	var dbConn *pgx.Conn
	var replicaConn *pgx.Conn

	// Создаём url для соединения через pgx
	dbConnectParameters := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=education sslmode=%s",
		config.DBHost,
		config.DBUserName,
		config.DBPassword,
		config.DBSSLMode,
	)

	// Пытаемся несколько раз подключиться к бд
	for attempt := 1; attempt <= 3; attempt++ {
		// Открываем соединение при помощи pgx
		c, err := pgx.Connect(context.Background(), dbConnectParameters)
		if err != nil {
			// Когда попытки закончились падаем окончательно
			if attempt == 3 {
				log.Fatal().Msgf("Не получается подключиться к бд (%d/3)", attempt)
			} else {
				log.Err(err).Msgf("Не получается подключиться к бд (%d/3)", attempt)
			}
			time.Sleep(time.Second)
		}
		dbConn = c
	}

	// Пытаемся подключиться к реплике
	for attempt := 1; attempt <= 3; attempt++ {
		// Открываем соединение при помощи pgx
		rc, err := pgx.Connect(context.Background(), dbConnectParameters)
		if err != nil {
			log.Err(err).Msgf("Не получается подключиться к реплике (%d/3)", attempt)
			time.Sleep(time.Second)
		}
		replicaConn = rc
	}

	// Создаём переменную для отправки запросов
	queries := New(dbConn)
	queriesReplca := New(replicaConn)

	// Создаём лямбда-замыкание для закрытия соединений
	closeConnect := func() {
		func(ctx context.Context, conn *pgx.Conn) {
			if err := conn.Close(ctx); err != nil {
				log.Fatal().Err(err).Msg("Не получается закрыть соединение с бд")
			} else {
				log.Info().Msg("Закрыли соединение с бд")
			}
		}(context.Background(), dbConn)

		func(ctx context.Context, conn *pgx.Conn) {
			if err := conn.Close(ctx); err != nil {
				log.Fatal().Err(err).Msg("Не получается закрыть соединение с репликой")
			} else {
				log.Info().Msg("Закрыли соединение с репликой")
			}
		}(context.Background(), replicaConn)
	}

	return queries, queriesReplca, closeConnect
}
