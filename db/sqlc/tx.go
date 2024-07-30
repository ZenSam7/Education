package db

import (
	"context"
	"github.com/jackc/pgx/v5"
)

// MakeTx создаём и выполняем новую транзакцию
func MakeTx(ctx context.Context, fn func(Querier) error) error {
	conn, ok := queries.db.(*pgx.Conn)
	if !ok {
		return pgx.ErrTxClosed
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	// Создание нового объекта queries с транзакцией
	qtx := queries.WithTx(tx)

	// Выполнение переданной функции fn
	err = fn(qtx)
	if err != nil {
		// Откат транзакции в случае ошибки
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return rollbackErr
		}
		return err
	}

	// Завершение транзакции
	return tx.Commit(ctx)
}
