// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	// BlockSession Блокируем сессию по id
	BlockSession(ctx context.Context, idSession pgtype.UUID) (Session, error)
	// CountRowsArticle Считаем количество строк в таблице
	CountRowsArticle(ctx context.Context) (int64, error)
	// CountRowsComment Считаем количество строк в таблице
	CountRowsComment(ctx context.Context) (int64, error)
	// CountRowsSessions Считаем количество строк в таблице
	CountRowsSessions(ctx context.Context) (int64, error)
	// CountRowsUser Считаем количество строк в таблице
	CountRowsUser(ctx context.Context) (int64, error)
	// CreateArticle Создаём статью
	CreateArticle(ctx context.Context, arg CreateArticleParams) (Article, error)
	// CreateComment Создаём комментарий к статье
	CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error)
	// CreateSession Создаём сессию
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	// CreateUser Создаём пользователя
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	// CreateVerifyRequest Создаём новый запрос на верификацию почты
	CreateVerifyRequest(ctx context.Context, arg CreateVerifyRequestParams) (VerifyEmail, error)
	// DeleteArticle Удаляем статью и комментарии к ней
	DeleteArticle(ctx context.Context, idArticle int32) (Article, error)
	// DeleteComment Удаляем комментарий к статье
	DeleteComment(ctx context.Context, idComment int32) (Article, error)
	// DeleteSession Удаляем сессию по id
	DeleteSession(ctx context.Context, idSession pgtype.UUID) (Session, error)
	// DeleteUser Удаляем пользователя
	DeleteUser(ctx context.Context, idUser int32) (User, error)
	// DeleteVerifyRequest Удаляем запрос на верификацию
	DeleteVerifyRequest(ctx context.Context, idUser int32) (VerifyEmail, error)
	// EditArticle Изменяем параметр(ы) статьи
	EditArticle(ctx context.Context, arg EditArticleParams) (Article, error)
	// EditComment Изменяем параметр(ы) пользователя
	EditComment(ctx context.Context, arg EditCommentParams) (Comment, error)
	// EditUser Изменяем параметр(ы) пользователя
	EditUser(ctx context.Context, arg EditUserParams) (User, error)
	// GetArticle Возвращаем статью по id
	GetArticle(ctx context.Context, idArticle int32) (Article, error)
	// GetArticlesWithAttribute Возвращаем много статей взятых по какому-то признаку(ам)
	GetArticlesWithAttribute(ctx context.Context, arg GetArticlesWithAttributeParams) ([]Article, error)
	// GetComment Возвращаем комментарий
	GetComment(ctx context.Context, idComment int32) (Comment, error)
	// GetCommentsOfArticle Возвращаем комментарии
	GetCommentsOfArticle(ctx context.Context, arg GetCommentsOfArticleParams) ([]Comment, error)
	// GetManySortedArticles Возвращаем много отсортированных статей
	GetManySortedArticles(ctx context.Context, arg GetManySortedArticlesParams) ([]Article, error)
	// GetManySortedArticlesWithAttribute Возвращаем много статей взятых по признаку по
	// какому-то признаку(ам) и отсортированных по другому признаку(ам)
	GetManySortedArticlesWithAttribute(ctx context.Context, arg GetManySortedArticlesWithAttributeParams) ([]Article, error)
	// GetManySortedUsers Возвращаем слайс пользователей отсортированных по какому-то параметру
	// (можно поставить: id_user, и сортировки не будет)
	GetManySortedUsers(ctx context.Context, arg GetManySortedUsersParams) ([]User, error)
	// GetSession Получаем сессиб по id
	GetSession(ctx context.Context, idSession pgtype.UUID) (Session, error)
	// GetUser Возвращаем пользователя
	GetUser(ctx context.Context, idUser int32) (User, error)
	// GetUserFromName Возвращаем пользователя по имени
	GetUserFromName(ctx context.Context, name string) (User, error)
	// GetVerifyRequest Возвращаем запрос на верификацию
	GetVerifyRequest(ctx context.Context, idUser int32) (VerifyEmail, error)
	// SetEmailIsVerified Ставим состояние почту как подтверждённую для какого-то пользователя
	SetEmailIsVerified(ctx context.Context, idUser int32) (User, error)
}

var _ Querier = (*Queries)(nil)
