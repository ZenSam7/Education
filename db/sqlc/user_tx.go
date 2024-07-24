package db

import (
	"context"
)

// Транзакции на все случаи жизни

type TxGetUserParams struct {
	IdUser int32 `json:"id_user"`
	// AfterCreate Вызываем эту функцию после выполнения основного запроса (в пределах одной транзакции)
	AfterCreate func(user User) error
}

type TxGetUserResponse struct {
	User User `json:"user"`
}

// GetUserTx Выполняем получение пользователя в одной транзакции с записью задачи в редиску
func (q *Queries) GetUserTx(ctx context.Context, arg TxGetUserParams) (TxGetUserResponse, error) {
	var result TxGetUserResponse

	err := q.MakeTx(ctx, func(qtx *Queries) error {
		var err error

		// Получение пользователя
		result.User, err = qtx.GetUser(ctx, arg.IdUser)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}
