package db

import (
	"context"
	"github.com/jackc/pgx/v5"
)

// MakeTx создаём и выполняем новую транзакцию
func (q *Queries) MakeTx(ctx context.Context, fn func(*Queries) error) error {
	conn, ok := q.db.(*pgx.Conn)
	if !ok {
		return pgx.ErrTxClosed
	}

	// Начало транзакции
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	// Создание нового объекта Queries с транзакцией
	qtx := q.WithTx(tx)

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
