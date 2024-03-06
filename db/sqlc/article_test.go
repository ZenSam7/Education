package db

import (
	"context"
	"github.com/ZenSam7/Education/tools"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// createRandomArticle Создаём случайную статью и возвращаем её
func createRandomArticle() (Article, *Queries, func()) {
	queries, closeConn := GetQueries()

	arg := CreateArticleParams{
		Title:   tools.GetRandomString(),
		Text:    tools.GetRandomString(),
		Authors: []int32{tools.GetRandomUint()},
	}
	newArticle, _ := queries.CreateArticle(context.Background(), arg)

	return newArticle, queries, closeConn
}

func TestCreateArticle(t *testing.T) {
	queries, closeConn := GetQueries()
	defer closeConn()

	arg := CreateArticleParams{
		Title:   tools.GetRandomString(),
		Text:    tools.GetRandomString(),
		Authors: []int32{tools.GetRandomUint()},
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

func TestEditArticle(t *testing.T) {
	article, queries, closeConn := createRandomArticle()
	defer closeConn() // Не забываем закрыть соединение

	// Измяняем Заголовок
	arg := EditArticleParams{
		IDArticle: article.IDArticle,
		Title:     tools.GetRandomString(),
	}

	editedArticle, err := queries.EditArticle(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, editedArticle)

	require.Equal(t, editedArticle.IDArticle, article.IDArticle)
	require.Equal(t, editedArticle.CreatedAt, article.CreatedAt)
	require.WithinDuration(t, editedArticle.EditedAt.Time, time.Now(), time.Second)
	require.NotEqual(t, editedArticle.Title, article.Title)
	require.Equal(t, editedArticle.Text, article.Text)
	require.Equal(t, editedArticle.Comments, article.Comments)
	require.Equal(t, editedArticle.Authors, article.Authors)
	require.Equal(t, editedArticle.Evaluation, article.Evaluation)
	article = editedArticle

	// Измяняем Оценку
	arg = EditArticleParams{
		IDArticle:  article.IDArticle,
		Evaluation: tools.GetRandomInt(),
	}

	editedArticle, err = queries.EditArticle(context.Background(), arg)
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
	require.WithinDuration(t, editedArticle.EditedAt.Time, time.Now(), time.Second)
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

	// По сути просто берём 1 статью
	arg := GetArticlesWithAttributeParams{
		Title:  createdArticles[0].Title,
		Offset: 0,
		Limit:  10000000,
	}

	article, err := queries.GetArticlesWithAttribute(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, article)

	require.Equal(t, createdArticles[0].EditedAt.Time,
		time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC))
	require.Equal(t, createdArticles[0].Evaluation, article[0].Evaluation)
	require.Equal(t, createdArticles[0].Authors, article[0].Authors)
	require.Equal(t, createdArticles[0].Comments, article[0].Comments)
	require.Equal(t, createdArticles[0].Text, article[0].Text)
	require.Equal(t, createdArticles[0].Title, article[0].Title)
}

func TestGetManySortedArticlesWithAttribute(t *testing.T) {
	_, queries, closeConn := createRandomArticle()
	defer closeConn() // Не забываем закрыть соединение

	// Создаём 10 статей с оценкой 0
	var createdArticles [10]Article
	for i := 0; i < 10; i++ {
		art, _, cC := createRandomArticle()
		cC() // Закрываем лишние соединения

		createdArticles[i] = art
	}

	arg := GetManySortedArticlesWithAttributeParams{
		SelectTitle:     createdArticles[0].Title,
		SortedIDArticle: true,
		Offset:          0,
		Limit:           10000000,
	}

	article, err := queries.GetManySortedArticlesWithAttribute(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, article)

	require.Equal(t, createdArticles[0].EditedAt.Time,
		time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC))
	require.Equal(t, createdArticles[0].Evaluation, article[0].Evaluation)
	require.Equal(t, createdArticles[0].Authors, article[0].Authors)
	require.Equal(t, createdArticles[0].Comments, article[0].Comments)
	require.Equal(t, createdArticles[0].Text, article[0].Text)
	require.Equal(t, createdArticles[0].Title, article[0].Title)

	// Сортируем статьи по ID
	arg = GetManySortedArticlesWithAttributeParams{
		SortedIDArticle: true,
		Offset:          0,
		Limit:           10,
	}

	article, err = queries.GetManySortedArticlesWithAttribute(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, article)

	for ind, art := range article[:9] {
		require.Equal(t, art.EditedAt.Time, time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC))
		require.True(t, article[ind].IDArticle < article[ind+1].IDArticle)
		require.NotEmpty(t, art.Title)
		require.NotEmpty(t, art.Text)
		require.NotEmpty(t, art.Authors)
		require.Zero(t, art.Evaluation)
		require.Empty(t, art.Comments)
	}
}
