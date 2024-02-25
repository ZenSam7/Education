package db

// Файл с тестами для запросов

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func GetRandomString() string {
	const letters = " abcdefghijklmnopqrstuvwxyz ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Минимальная длина: 2
	str := make([]byte, rand.Intn(len(letters)-2)+2)
	for i := range str {
		str[i] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}

// GetRandomInt Число может быть отрицательным
func GetRandomInt() int32 {
	return rand.Int31() * int32(1-2*rand.Intn(2))
}

// createRandomUser Создаём случайного пользователя (заодно тестируем его),
// queries для отправки запросов, функцию для закрытия соединения и возвращаем всё это
func createRandomUser() (User, *Queries, func()) {
	queries, closeConn := GetQueries()

	arg := CreateUserParams{
		Name:        GetRandomString(),
		Description: GetRandomString(),
	}
	newUser, _ := queries.CreateUser(context.Background(), arg)

	return newUser, queries, closeConn
}

func TestCreateUser(t *testing.T) {
	queries, closeConn := GetQueries()
	defer closeConn()

	arg := CreateUserParams{
		Name:        GetRandomString(),
		Description: GetRandomString(),
	}

	newUser, err := queries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	require.Equal(t, arg.Name, newUser.Name)
	require.Equal(t, arg.Description, newUser.Description)
	require.Zero(t, newUser.Karma)

	require.NotZero(t, newUser.IDUser)
	require.NotZero(t, newUser.CreatedAt)
}

func TestEditUserParam(t *testing.T) {
	user, queries, closeConn := createRandomUser()
	// Не забываем закрыть соединение
	defer closeConn()

	// Изменяем Имя
	agr := EditUserParamParams{
		Name:   GetRandomString(),
		IDUser: user.IDUser,
	}

	editedUser, err := queries.EditUserParam(context.Background(), agr)
	require.NoError(t, err)
	require.NotEmpty(t, editedUser)

	require.NotEqual(t, editedUser.Name, user.Name)

	require.Equal(t, editedUser.IDUser, user.IDUser)
	require.Equal(t, editedUser.Karma, user.Karma)
	require.Equal(t, editedUser.Description, user.Description)
	user = editedUser // Обновляем пользователя с которым сравниваем

	// Изменяем Описание
	agr = EditUserParamParams{
		Description: GetRandomString(),
		IDUser:      user.IDUser,
	}

	editedUser, err = queries.EditUserParam(context.Background(), agr)
	require.NoError(t, err)
	require.NotEmpty(t, editedUser)

	require.NotEqual(t, editedUser.Description, user.Description)

	require.Equal(t, editedUser.IDUser, user.IDUser)
	require.Equal(t, editedUser.Karma, user.Karma)
	require.Equal(t, editedUser.Name, user.Name)
	user = editedUser // Обновляем пользователя с которым сравниваем

	// Изменяем Карму
	agr = EditUserParamParams{
		Karma:  GetRandomInt(),
		IDUser: user.IDUser,
	}

	editedUser, err = queries.EditUserParam(context.Background(), agr)
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
	require.EqualError(t, err, pgx.ErrNoRows.Error())

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

	arg := GetManySortedUsersParams{
		Offset: 0,
		Limit:  10,
		IDUser: true,
	}

	user, err := queries.GetManySortedUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	for _, usr := range user {
		require.NotEmpty(t, usr.IDUser)
		require.NotEmpty(t, usr.CreatedAt)
		require.NotEmpty(t, usr.Name)
		require.NotEmpty(t, usr.Description)
	}
}
