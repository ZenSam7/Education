package db

// Файл с тестами для запросов

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func getRandomString() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	str := make([]byte, rand.Intn(len(letters)))
	for i := range str {
		str[i] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}

// getRandomInt Число может быть отрицательным
func getRandomInt() int32 {
	return rand.Int31() * int32(1-2*rand.Intn(2))
}

// createRandomUser Создаём случайного пользователя (заодно тестируем его),
// queries для отправки запросов, функцию для закрытия соединения и возвращаем всё это
func createRandomUser(t *testing.T) (User, *Queries, func()) {
	queries, closeConn := GetQueries()

	arg := CreateUserParams{
		Name:        getRandomString(),
		Description: getRandomString(),
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
		Name:        getRandomString(),
		Description: user.Description,
		Karma:       user.Karma,
		IDUser:      user.IDUser,
	}

	editedUser, err := queries.EditUserParam(context.Background(), agr)
	require.NoError(t, err)

	require.NotEqual(t, editedUser.Name, user.Name)

	require.Equal(t, editedUser.IDUser, user.IDUser)
	require.Equal(t, editedUser.Karma, user.Karma)
	require.Equal(t, editedUser.Description, user.Description)

	// Изменяем Описание
	agr = EditUserParamParams{
		Name:        user.Name,
		Description: getRandomString(),
		Karma:       user.Karma,
		IDUser:      user.IDUser,
	}

	editedUser, err = queries.EditUserParam(context.Background(), agr)
	require.NoError(t, err)

	require.NotEqual(t, editedUser.Description, user.Description)

	require.Equal(t, editedUser.IDUser, user.IDUser)
	require.Equal(t, editedUser.Karma, user.Karma)
	require.Equal(t, editedUser.Name, user.Name)

	// Изменяем Карму
	agr = EditUserParamParams{
		Name:        user.Name,
		Description: user.Description,
		Karma:       getRandomInt(),
		IDUser:      user.IDUser,
	}

	editedUser, err = queries.EditUserParam(context.Background(), agr)
	require.NoError(t, err)

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

	require.Equal(t, findedUser.IDUser, user.IDUser)
	require.Equal(t, findedUser.Description, user.Description)
	require.Equal(t, findedUser.Karma, user.Karma)
	require.Equal(t, findedUser.Name, user.Name)
	require.Equal(t, findedUser.CreatedAt, user.CreatedAt)
}

func TestGetManyUsers(t *testing.T) {
	_, queries, closeConn := createRandomUser(t)
	// Не забываем закрыть соединение
	defer closeConn()

	// Создаём 10 тестовых пользователя
	var listCreatedUsers []User
	var topUserId int64 // Id оследнего пользователя
	for i := 0; i < 10; i++ {
		usr, _, cC := createRandomUser(t)
		cC() // Закрываем лишние соединения

		topUserId = int64(usr.IDUser)
		listCreatedUsers = append(listCreatedUsers, usr)
	}

	arg := GetManyUsersParams{
		Limit:     4,
		Offset:    topUserId,
		Attribute: "id_user",
	}

	listFindedUsers, err := queries.GetManyUsers(context.Background(), arg)
	require.NoError(t, err)

	for i, usr := range listFindedUsers {
		require.Equal(t, usr.IDUser, listCreatedUsers[i].IDUser)
		require.Equal(t, usr.Description, listCreatedUsers[i].Description)
		require.Equal(t, usr.Karma, listCreatedUsers[i].Karma)
		require.Equal(t, usr.Name, listCreatedUsers[i].Name)
		require.Equal(t, usr.CreatedAt, listCreatedUsers[i].CreatedAt)
	}
}

func TestDeleteUser(t *testing.T) {
	user, queries, closeConn := createRandomUser(t)
	// Не забываем закрыть соединение
	defer closeConn()

	deletedUser, err := queries.DeleteUser(context.Background(), user.IDUser)
	require.NoError(t, err)

	require.Equal(t, deletedUser.IDUser, user.IDUser)
	require.Equal(t, deletedUser.Description, user.Description)
	require.Equal(t, deletedUser.Karma, user.Karma)
	require.Equal(t, deletedUser.Name, user.Name)
	require.Equal(t, deletedUser.CreatedAt, user.CreatedAt)

	findedUser, err := queries.GetUser(context.Background(), user.IDUser)
	require.Error(t, err)
	require.EqualError(t, err, pgx.ErrNoRows.Error())

	require.Empty(t, findedUser)

	// Если функция удаления пользователей работает, то значит
}
