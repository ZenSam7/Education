package db

import (
	"context"
)

// Транзакции на все случаи жизни

type TxCreateUserParams struct {
	Name         string `json:"name" binding:"required,alphanum"`
	Email        string `json:"email" binding:"required,email"`
	PasswordHash string `json:"password" binding:"required"`
	// AfterCreate Вызываем эту функцию после выполнения основного запроса
	// (в пределах одной транзакции) (отсюда вызываем создание таски)
	AfterCreate func(user User) error
}

type TxCreateUserResponse struct {
	User User `json:"user"`
}

// CreateUserTx Выполняем создание пользователя в одной транзакции с записью задачи в редиску
func (q *Queries) CreateUserTx(ctx context.Context, arg TxCreateUserParams) (TxCreateUserResponse, error) {
	var result TxCreateUserResponse

	err := q.MakeTx(ctx, func(qtx *Queries) error {
		var err error

		// Создание пользователя
		result.User, err = qtx.CreateUser(ctx, CreateUserParams{
			Name:         arg.Name,
			PasswordHash: arg.PasswordHash,
			Email:        arg.Email,
		})
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}
