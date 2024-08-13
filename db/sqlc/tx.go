package db

import (
	"context"
	"github.com/jackc/pgx/v5"
)

// MakeTx создаём и выполняем новую транзакцию
func MakeTx(ctx context.Context, fn func(Querier) error, mayUseReplica ...bool) error {
	if queries == nil || replica == nil {
		// Не закрываем соединение (да, это плохо)
		queries, replica, _ = MakeQueries()
	}

	conn, ok := queries.db.(*pgx.Conn)
	if !ok {
		// Пробуем использовать реплику
		if mayUseReplica[0] {
			conn, ok = replica.db.(*pgx.Conn)
			if !ok {
				return pgx.ErrTxClosed
			}
		}

		return pgx.ErrTxClosed
	}

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}

	// Выполнение переданной функции fn
	err = fn(queries.WithTx(tx))

	if err != nil {
		// Откат транзакции в случае ошибки
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			return rollbackErr
		}

		// Если можно попытаться использовать реплику, то используем реплику
		if mayUseReplica[0] {
			err = fn(replica.WithTx(tx))
			// Откат транзакции
			if err != nil {
				if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
					return rollbackErr
				}
				return err
			}
			// Завершение транзакции реплики
			return tx.Commit(ctx)
		}

		return err
	}

	// Завершение транзакции
	return tx.Commit(ctx)
}
