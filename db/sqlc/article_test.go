package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// createRandomArticle Создаём случайную статью (заодно тестируем её) и возвращаем её
func createRandomArticle() (Article, *Queries, func()) {
	queries, closeConn := GetQueries()

	arg := CreateArticleParams{
		Title:   GetRandomString(),
		Text:    GetRandomString(),
		Authors: []int32{GetRandomInt()},
	}
	newArticle, _ := queries.CreateArticle(context.Background(), arg)

	return newArticle, queries, closeConn
}

func TestCreateArticle(t *testing.T) {
	queries, closeConn := GetQueries()
	defer closeConn()

	arg := CreateArticleParams{
		Title:   GetRandomString(),
		Text:    GetRandomString(),
		Authors: []int32{GetRandomInt()},
	}

	newArticle, err := queries.CreateArticle(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newArticle)

	require.NotZero(t, newArticle.IDArticle)
	require.WithinDuration(t, newArticle.CreatedAt.Time, time.Now(), time.Second)
	require.Equal(t, newArticle.EditedAt.Time, time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC))
	require.Equal(t, newArticle.Title, arg.Title)
	require.Equal(t, newArticle.Text, arg.Text)
	require.Nil(t, newArticle.Comments)
	require.Equal(t, newArticle.Authors, arg.Authors)
	require.Zero(t, newArticle.Evaluation)
}

func TestGetArticle(t *testing.T) {
	article, queries, closeConn := createRandomArticle()
	defer closeConn() // Не забываем закрыть соединение

	findedArticle, err := queries.GetArticle(context.Background(), article.IDArticle)
	require.NoError(t, err)
	require.NotEmpty(t, findedArticle)
	require.Equal(t, findedArticle.IDArticle, article.IDArticle)
	require.Equal(t, findedArticle.CreatedAt, article.CreatedAt)
	require.Equal(t, findedArticle.EditedAt, article.EditedAt)
	require.Equal(t, findedArticle.Title, article.Title)
	require.Equal(t, findedArticle.Text, article.Text)
	require.Equal(t, findedArticle.Comments, article.Comments)
	require.Equal(t, findedArticle.Authors, article.Authors)
	require.Equal(t, findedArticle.Evaluation, article.Evaluation)
}

func TestEditArticleParam(t *testing.T) {
	article, queries, closeConn := createRandomArticle()
	defer closeConn() // Не забываем закрыть соединение

	// Измяняем Заголовок
	arg := EditArticleParamParams{
		IDArticle: article.IDArticle,
		Title:     GetRandomString(),
	}

	editedArticle, err := queries.EditArticleParam(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, editedArticle)

	require.Equal(t, editedArticle.IDArticle, article.IDArticle)
	require.Equal(t, editedArticle.CreatedAt, article.CreatedAt)
	require.Equal(t, editedArticle.EditedAt, article.EditedAt)
	require.NotEqual(t, editedArticle.Title, article.Title)
	require.Equal(t, editedArticle.Text, article.Text)
	require.Equal(t, editedArticle.Comments, article.Comments)
	require.Equal(t, editedArticle.Authors, article.Authors)
	require.Equal(t, editedArticle.Evaluation, article.Evaluation)
	article = editedArticle

	// Измяняем Оценку
	arg = EditArticleParamParams{
		IDArticle:  article.IDArticle,
		Evaluation: GetRandomInt(),
	}

	editedArticle, err = queries.EditArticleParam(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, editedArticle)

	require.Equal(t, editedArticle.IDArticle, article.IDArticle)
	require.Equal(t, editedArticle.CreatedAt, article.CreatedAt)
	require.Equal(t, editedArticle.EditedAt, article.EditedAt)
	require.Equal(t, editedArticle.Title, article.Title)
	require.Equal(t, editedArticle.Text, article.Text)
	require.Equal(t, editedArticle.Comments, article.Comments)
	require.Equal(t, editedArticle.Authors, article.Authors)
	require.NotEqual(t, editedArticle.Evaluation, article.Evaluation)
}

func TestGetArticlesWithAttribute(t *testing.T) {
	_, queries, closeConn := createRandomArticle()
	defer closeConn() // Не забываем закрыть соединение

	// Создаём 10 статей с оценкой 0
	var createdArticles [10]Article
	for i := 0; i < 10; i++ {
		art, _, cC := createRandomArticle()
		cC() // Закрываем лишние соединения

		createdArticles[i] = art
	}

	arg := GetArticlesWithAttributeParams{
		Evaluation: 0,
		Offset:     0,
		Limit:      10,
	}

	articles, err := queries.GetArticlesWithAttribute(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, articles)

	for _, article := range articles {
		require.Zero(t, article.Evaluation)
		require.Equal(t, article.EditedAt.Time, time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC))
		require.NotEmpty(t, article.Authors)
		require.Empty(t, article.Comments)
		require.NotEmpty(t, article.Text)
		require.NotEmpty(t, article.Title)
	}
}

// TODO: Доделать все тесты для article
