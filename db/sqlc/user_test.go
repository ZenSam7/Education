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
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	str := make([]byte, rand.Intn(len(letters)))
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
func createRandomUser(t *testing.T) (User, *Queries, func()) {
	queries, closeConn := GetQueries()

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

	return newUser, queries, closeConn
}

func TestCreateUser(t *testing.T) {
	_, _, closeConn := createRandomUser(t)
	// Не забываем закрыть соединение
	defer closeConn()
}

func TestEditUserParam(t *testing.T) {
	user, queries, closeConn := createRandomUser(t)
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
	user, queries, closeConn := createRandomUser(t)
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

// TODO: Переделать тест TestGetManyUsers и запрос как у TestGetArticlesWithAttribute
func TestGetManyUsers(t *testing.T) {
	_, queries, closeConn := createRandomUser(t)
	// Не забываем закрыть соединение
	defer closeConn()

	var listUsers [6]User
	for i := 0; i < 6; i++ {
		listUsers[i], _ = queries.GetUser(context.Background(), int32(i+1))
	}

	arg := GetManySortedUsersParams{
		Attribute: "name",
		Offset:    0,
		Limit:     1,
	}

	listFindedUsers, err := queries.GetManySortedUsers(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, listFindedUsers)

	for i, usr := range listFindedUsers {
		ind := i + int(arg.Offset)
		require.Equal(t, usr.IDUser, listUsers[ind].IDUser)
		require.Equal(t, usr.Description, listUsers[ind].Description)
		require.Equal(t, usr.Karma, listUsers[ind].Karma)
		require.Equal(t, usr.Name, listUsers[ind].Name)
		require.Equal(t, usr.CreatedAt, listUsers[ind].CreatedAt)
	}
}

func TestDeleteUser(t *testing.T) {
	user, queries, closeConn := createRandomUser(t)
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
