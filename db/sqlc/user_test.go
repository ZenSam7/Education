package db

// Файл с тестами для запросов

import (
	"context"
	"database/sql"
	"github.com/ZenSam7/Education/tools"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"testing"
)

// createRandomUser Создаём случайного пользователя и возвращаем функцию для закрытия соединения
func createRandomUser() (User, *Queries, func()) {
	queries, _, closeConn := MakeQueries()

	arg := CreateUserParams{
		Name:         tools.GetRandomString(),
		Email:        tools.GetRandomEmail(),
		PasswordHash: tools.GetRandomHash(),
	}
	newUser, _ := queries.CreateUser(context.Background(), arg)

	return newUser, queries, closeConn
}

func TestCreateUser(t *testing.T) {
	queries, _, closeConn := MakeQueries()
	defer closeConn()

	arg := CreateUserParams{
		Name:         tools.GetRandomString(),
		Email:        tools.GetRandomEmail(),
		PasswordHash: tools.GetRandomHash(),
	}

	newUser, err := queries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, arg.Name, newUser.Name)
	require.Equal(t, arg.Email, newUser.Email)
	require.Equal(t, arg.PasswordHash, newUser.PasswordHash)
	require.Zero(t, newUser.Karma)

	require.NotZero(t, newUser.IDUser)
	require.NotZero(t, newUser.CreatedAt)
}

func TestEditUser(t *testing.T) {
	user, queries, closeConn := createRandomUser()
	// Не забываем закрыть соединение
	defer closeConn()

	// Изменяем Имя
	agr := EditUserParams{
		Name:   tools.GetRandomString(),
		IDUser: user.IDUser,
	}

	editedUser, err := queries.EditUser(context.Background(), agr)
	require.NoError(t, err)
	require.NotEmpty(t, editedUser)

	require.NotEqual(t, editedUser.Name, user.Name)

	require.Equal(t, editedUser.IDUser, user.IDUser)
	require.Equal(t, editedUser.Karma, user.Karma)
	require.Equal(t, editedUser.Description, user.Description)
	user = editedUser // Обновляем пользователя с которым сравниваем

	// Изменяем Описание
	agr = EditUserParams{
		Description: pgtype.Text{String: tools.GetRandomString(), Valid: true},
		IDUser:      user.IDUser,
	}

	editedUser, err = queries.EditUser(context.Background(), agr)
	require.NoError(t, err)
	require.NotEmpty(t, editedUser)

	require.NotEqual(t, editedUser.Description, user.Description)

	require.Equal(t, editedUser.IDUser, user.IDUser)
	require.Equal(t, editedUser.Karma, user.Karma)
	require.Equal(t, editedUser.Name, user.Name)
	user = editedUser // Обновляем пользователя с которым сравниваем

	// Изменяем Карму
	agr = EditUserParams{
		Karma:  pgtype.Int4{Int32: tools.GetRandomInt(), Valid: true},
		IDUser: user.IDUser,
	}

	editedUser, err = queries.EditUser(context.Background(), agr)
	require.NoError(t, err)
	require.NotEmpty(t, editedUser)

	require.NotEqual(t, editedUser.Karma, user.Karma)

	require.Equal(t, editedUser.IDUser, user.IDUser)
	require.Equal(t, editedUser.Description, user.Description)
	require.Equal(t, editedUser.Name, user.Name)
}

func TestGetUser(t *testing.T) {
	user, queries, closeConn := createRandomUser()
	// Не забываем закрыть соединение
	defer closeConn()

	findedUser, err := queries.GetUser(context.Background(), user.IDUser)
	require.NoError(t, err)
	require.NotEmpty(t, findedUser)

	require.Equal(t, findedUser.IDUser, user.IDUser)
	require.Equal(t, findedUser.Description, user.Description)
	require.Equal(t, findedUser.Karma, user.Karma)
	require.Equal(t, findedUser.Name, user.Name)
	require.Equal(t, findedUser.CreatedAt, user.CreatedAt)
}

func TestDeleteUser(t *testing.T) {
	user, queries, closeConn := createRandomUser()
	// Не забываем закрыть соединение
	defer closeConn()

	deletedUser, err := queries.DeleteUser(context.Background(), user.IDUser)
	require.NoError(t, err)
	require.NotEmpty(t, deletedUser)

	require.Equal(t, deletedUser.IDUser, user.IDUser)
	require.Equal(t, deletedUser.Description, user.Description)
	require.Equal(t, deletedUser.Karma, user.Karma)
	require.Equal(t, deletedUser.Name, user.Name)
	require.Equal(t, deletedUser.CreatedAt, user.CreatedAt)

	findedUser, err := queries.GetUser(context.Background(), user.IDUser)
	require.Error(t, err)
	require.EqualError(t, sql.ErrNoRows, "sql: "+err.Error())

	require.Empty(t, findedUser)
}

func TestGetManySortedUsers(t *testing.T) {
	_, queries, closeConn := createRandomUser()
	// Не забываем закрыть соединение
	defer closeConn()

	// Создаём 10 пользователей
	var createdUsers [10]User
	for i := 0; i < 10; i++ {
		usr, _, cC := createRandomUser()
		cC() // Закрываем лишние соединения

		createdUsers[i] = usr
	}

	countRows, err := queries.CountRowsUser(context.Background())
	require.NoError(t, err)

	arg := GetManySortedUsersParams{
		Offset: int32(countRows) - 10,
		Limit:  10,
		IDUser: true,
	}

	log.Info().Int64("", countRows).Msg("")
	users, err := queries.GetManySortedUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, users)

	for _, usr := range users {
		require.NotEmpty(t, usr.IDUser)
		require.NotEmpty(t, usr.CreatedAt)
		require.NotEmpty(t, usr.Name)
		require.Empty(t, usr.Description)
	}
}
