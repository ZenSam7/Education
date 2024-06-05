package db

import (
	"context"
	"github.com/ZenSam7/Education/tools"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateComment(t *testing.T) {
	queries, closeConn := GetQueries()
	defer closeConn()

	arg := CreateCommentParams{
		Text:      tools.GetRandomString(),
		Author:    tools.GetRandomUint(),
		IDArticle: tools.GetRandomUint(),
	}

	newComment, err := queries.CreateComment(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newComment)

	require.NotZero(t, newComment.IDComment)
	require.WithinDuration(t, newComment.CreatedAt.Time, time.Now(), time.Second)
	require.Equal(t, newComment.EditedAt.Time, time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC))
	require.Equal(t, newComment.Text, arg.Text)
	require.Equal(t, newComment.Author, arg.Author)
}

func TestGetComment(t *testing.T) {
	queries, closeConn := GetQueries()
	defer closeConn()

	// Создаём комментарий
	arg := CreateCommentParams{
		Text:      tools.GetRandomString(),
		Author:    tools.GetRandomUint(),
		IDArticle: tools.GetRandomUint(),
	}
	comm, err := queries.CreateComment(context.Background(), arg)
	require.NoError(t, err)

	newComm, err := queries.GetComment(context.Background(), comm.IDComment)
	require.NoError(t, err)
	require.Equal(t, comm, newComm)
	require.Equal(t, newComm.IDComment, comm.IDComment)
	require.Equal(t, newComm.CreatedAt, comm.CreatedAt)
	require.Equal(t, newComm.EditedAt, comm.EditedAt)
	require.Equal(t, newComm.Text, comm.Text)
	require.Equal(t, newComm.Author, comm.Author)
}

func TestDeleteComment(t *testing.T) {
	article, queries, closeConn := createRandomArticle()
	defer closeConn()

	// Создаём комментарий
	arg := CreateCommentParams{
		Text:      tools.GetRandomString(),
		Author:    tools.GetRandomUint(),
		IDArticle: article.IDArticle,
	}
	comm, err := queries.CreateComment(context.Background(), arg)
	require.NoError(t, err)

	// Проверяем что коммент создался
	article, err = queries.GetArticle(context.Background(), article.IDArticle)
	require.NoError(t, err)
	require.NotEmpty(t, article.Comments)
	require.NotZero(t, article.Comments[0])
	_, err = queries.GetComment(context.Background(), article.Comments[0])
	require.NoError(t, err)
	require.Equal(t, comm.IDComment, article.Comments[0])

	// Ничего кроме комментария статье не изменилось
	editedArticle, err := queries.DeleteComment(context.Background(), article.Comments[0])
	require.NoError(t, err)
	require.Equal(t, editedArticle.Comments, []int32{})
	require.Equal(t, editedArticle.Evaluation, article.Evaluation)
	require.Equal(t, editedArticle.Title, article.Title)
	require.Equal(t, editedArticle.Text, article.Text)
	require.Equal(t, editedArticle.Authors, article.Authors)
	require.Equal(t, editedArticle.CreatedAt, article.CreatedAt)

	comm, err = queries.GetComment(context.Background(), article.Comments[0])
	require.Error(t, err)
}

func TestEditComment(t *testing.T) {
	article, queries, closeConn := createRandomArticle()
	defer closeConn()

	// Создаём комментарий
	arg := CreateCommentParams{
		Text:      tools.GetRandomString(),
		Author:    tools.GetRandomUint(),
		IDArticle: article.IDArticle,
	}
	comm, err := queries.CreateComment(context.Background(), arg)
	require.NoError(t, err)

	// Проверяем что коммент создался
	article, err = queries.GetArticle(context.Background(), article.IDArticle)
	require.NoError(t, err)
	require.NotEmpty(t, article.Comments)
	require.NotZero(t, article.Comments[0])
	require.Equal(t, comm.IDComment, article.Comments[0])

	argsEditedComm := EditCommentParams{
		IDComment: comm.IDComment,
		Text:      tools.GetRandomString(),
	}

	editedComm, err := queries.EditComment(context.Background(), argsEditedComm)
	require.NoError(t, err)
	require.Equal(t, editedComm.Text, argsEditedComm.Text)
	require.WithinDuration(t, editedComm.EditedAt.Time, time.Now(), time.Second)

	editedComm, err = queries.GetComment(context.Background(), comm.IDComment)
	require.NoError(t, err)
	require.Equal(t, editedComm.Text, argsEditedComm.Text)
	require.NotEqual(t, editedComm.EditedAt, comm.EditedAt)
	require.NotEqual(t, editedComm.Text, comm.Text)
}
