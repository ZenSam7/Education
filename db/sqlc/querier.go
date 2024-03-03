// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"context"
)

type Querier interface {
	// CreateArticle Создаём статью
	CreateArticle(ctx context.Context, arg CreateArticleParams) (Article, error)
	// CreateComment Создаём комментарий к статье
	CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error)
	// CreateUser Создаём пользователя
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	// DeleteArticle Удаляем статью и комментарии к ней
	DeleteArticle(ctx context.Context, idArticle int32) (Article, error)
	// DeleteComment Удаляем комментарий к статье
	DeleteComment(ctx context.Context, idComment int32) (Article, error)
	// DeleteUser Удаляем пользователя и сдвигаем id
	DeleteUser(ctx context.Context, idUser int32) (User, error)
	// EditArticleParam Изменяем параметр(ы) статьи
	EditArticleParam(ctx context.Context, arg EditArticleParamParams) (Article, error)
	// EditCommentParam Изменяем параметр(ы) пользователя
	EditCommentParam(ctx context.Context, arg EditCommentParamParams) (Comment, error)
	// EditUserParam Изменяем параметр(ы) пользователя
	EditUserParam(ctx context.Context, arg EditUserParamParams) (User, error)
	// GetArticle Возвращаем статью по id
	GetArticle(ctx context.Context, idArticle int32) (Article, error)
	// GetArticlesWithAttribute Возвращаем много статей взятых по какому-то признаку(ам)
	GetArticlesWithAttribute(ctx context.Context, arg GetArticlesWithAttributeParams) ([]Article, error)
	// GetComment Возвращаем комментарий
	GetComment(ctx context.Context, idComment int32) (Comment, error)
	// GetManySortedArticles Возвращаем много отсортированных статей
	GetManySortedArticles(ctx context.Context, arg GetManySortedArticlesParams) ([]Article, error)
	// GetManySortedArticlesWithAttribute Возвращаем много статей взятых по признаку по
	// какому-то признаку(ам) и отсортированных по другому признаку(ам)
	GetManySortedArticlesWithAttribute(ctx context.Context, arg GetManySortedArticlesWithAttributeParams) ([]Article, error)
	// GetManySortedUsers Возвращаем слайс пользователей отсортированных по какому-то параметру
	// (можно поставить: id_user, и сортировки не будет)
	GetManySortedUsers(ctx context.Context, arg GetManySortedUsersParams) ([]User, error)
	// GetUser Возвращаем пользователя
	GetUser(ctx context.Context, idUser int32) (User, error)
}

var _ Querier = (*Queries)(nil)
