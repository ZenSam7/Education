package db

import (
	"context"
	"github.com/ZenSam7/Education/tools"
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

// CreateUserTx Выполняем создание пользователя в одной транзакции с записью задачи на верификацию почты в редиску
func (q *Queries) CreateUserTx(ctx context.Context, arg TxCreateUserParams) (TxCreateUserResponse, error) {
	var result TxCreateUserResponse

	err := q.MakeTx(ctx, func(q *Queries) error {
		var err error

		// Создание пользователя
		result.User, err = q.CreateUser(ctx, CreateUserParams{
			Name:         arg.Name,
			PasswordHash: arg.PasswordHash,
			Email:        arg.Email,
		})
		if err != nil {
			return err
		}

		// Запрос на верификацию
		_, err = q.CreateVerifyRequest(ctx, CreateVerifyRequestParams{
			IDUser:    result.User.IDUser,
			SecretKey: tools.GetRandomString(32),
		})
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}
