package api

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/ZenSam7/Education/db/sqlc"
	pb "github.com/ZenSam7/Education/protobuf"
	"github.com/ZenSam7/Education/tools"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func isAuthorArticle(IDArticle, IDUser int32, server *Server) bool {
	targetArticle, _ := server.querier.GetArticle(context.Background(), IDArticle)
	for _, authorID := range targetArticle.Authors {
		if authorID == IDUser {
			return true
		}
	}
	return false
}

func validateCreateArticleRequest(req *pb.CreateArticleRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateString(req.GetText(), 1, 99999); err != nil {
		errorsFields = append(errorsFields, fieldViolation("text", err))
	}
	if err := tools.ValidateString(req.GetTitle(), 1, 100); err != nil {
		errorsFields = append(errorsFields, fieldViolation("title", err))
	}

	return wrapFeildErrors(errorsFields)
}

func (server *Server) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleResponse, error) {
	if err := validateCreateArticleRequest(req); err != nil {
		return nil, err
	}

	// Только для авторизованных пользователей
	mtdt, err := server.extractMetadata(ctx)
	if err != nil {
		return nil, err
	}
	payload, err := server.tokenMaker.VerifyToken(mtdt.AccessToken)
	if err != nil {

	}

	article, err := server.querier.CreateArticle(ctx, db.CreateArticleParams{
		Title:   req.GetTitle(),
		Text:    req.GetText(),
		Authors: []int32{payload.IDUser},
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при создании статьи: %s", err)
	}

	return &pb.CreateArticleResponse{Article: convArticle(article)}, nil
}

func validateGetArticleRequest(req *pb.GetArticleRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateNaturalNum(int(req.GetIdArticle())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("id_article", err))
	}

	return wrapFeildErrors(errorsFields)
}
func (server *Server) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.GetArticleResponse, error) {
	// Попытка получения данных из кеша
	cacheKey := fmt.Sprintf("article:%d", req.GetIdArticle())
	var cachedResponse pb.GetArticleResponse
	if err := server.cacher.GetCache(ctx, cacheKey, &cachedResponse); err == nil {
		return &cachedResponse, nil
	}

	// Проверка валидности запроса
	if err := validateGetArticleRequest(req); err != nil {
		return nil, err
	}

	article, err := server.querier.GetArticle(ctx, req.GetIdArticle())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "статья не найдена")
		}
		return nil, status.Errorf(codes.Internal, "ошибка при получении статьи: %s", err)
	}

	response := &pb.GetArticleResponse{
		Article: convArticle(article),
	}

	if err := server.cacher.SetCache(ctx, cacheKey, response); err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось сохранить данные в кеш: %s", err)
	}

	return response, nil
}

func validateEditArticleRequest(req *pb.EditArticleRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateString(req.GetText(), 1, 9999); err != nil {
		errorsFields = append(errorsFields, fieldViolation("text", err))
	}
	if err := tools.ValidateString(req.GetTitle(), 1, 100); err != nil {
		errorsFields = append(errorsFields, fieldViolation("title", err))
	}
	if err := tools.ValidateNotEmpty(req.GetComments()); err != nil {
		errorsFields = append(errorsFields, fieldViolation("comments", err))
	}
	if err := tools.ValidateNotEmpty(req.GetAuthors()); err != nil {
		errorsFields = append(errorsFields, fieldViolation("authors", err))
	}

	return wrapFeildErrors(errorsFields)
}
func (server *Server) EditArticle(ctx context.Context, req *pb.EditArticleRequest) (*pb.EditArticleResponse, error) {
	if err := validateEditArticleRequest(req); err != nil {
		return nil, err
	}

	// Только для авторов
	mtdt, err := server.extractMetadata(ctx)
	if err != nil {
		return nil, err
	}
	payload, err := server.tokenMaker.VerifyToken(mtdt.AccessToken)
	if err != nil {

	}

	if !isAuthorArticle(payload.IDUser, payload.IDUser, server) {
		return nil, status.Errorf(codes.PermissionDenied, "вы не автор статьи")
	}

	editedArticle, err := server.querier.EditArticle(ctx, db.EditArticleParams{
		IDArticle: req.GetIdArticle(),
		Title:     req.GetTitle(),
		Text:      req.GetText(),
		Comments:  req.GetComments(),
		Authors:   req.GetAuthors(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при изменении статьи: %s", err)
	}

	return &pb.EditArticleResponse{Article: convArticle(editedArticle)}, nil
}

func validateDeleteArticleRequest(req *pb.DeleteArticleRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateNaturalNum(int(req.GetIdArticle())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("id_article", err))
	}

	return wrapFeildErrors(errorsFields)
}
func (server *Server) DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*pb.DeleteArticleResponse, error) {
	if err := validateDeleteArticleRequest(req); err != nil {
		return nil, err
	}

	// Только для авторов
	mtdt, err := server.extractMetadata(ctx)
	if err != nil {
		return nil, err
	}
	payload, err := server.tokenMaker.VerifyToken(mtdt.AccessToken)
	if err != nil {
		return nil, err
	}

	if !isAuthorArticle(payload.IDUser, payload.IDUser, server) {
		return nil, status.Errorf(codes.PermissionDenied, "вы не автор статьи")
	}

	_, err = server.querier.DeleteArticle(ctx, req.GetIdArticle())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка при удалении статьи: %s", err)
	}

	return &pb.DeleteArticleResponse{}, nil
}

func validateGetManySortedArticlesRequest(req *pb.GetManySortedArticlesRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateNaturalNum(int(req.GetPageNum())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("page_num", err))
	}
	if err := tools.ValidateNaturalNum(int(req.GetPageSize())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("page_size", err))
	}

	return wrapFeildErrors(errorsFields)
}
func (server *Server) GetManySortedArticles(ctx context.Context, req *pb.GetManySortedArticlesRequest) (*pb.GetManySortedArticlesResponse, error) {
	if err := validateGetManySortedArticlesRequest(req); err != nil {
		return nil, err
	}

	articles, err := server.querier.GetManySortedArticles(ctx, db.GetManySortedArticlesParams{
		IDArticle:  req.GetIdArticle(),
		Title:      req.GetTitle(),
		Text:       req.GetText(),
		Comments:   req.GetComments(),
		Authors:    req.GetAuthors(),
		Evaluation: req.GetEvaluation(),
		EditedAt:   req.GetEditedAt(),
		CreatedAt:  req.GetCreatedAt(),
		Limit:      req.GetPageSize(),
		Offset:     (req.GetPageNum() - 1) * req.GetPageSize(),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "статьи не найдены")
		}
		return nil, status.Errorf(codes.Internal, "не удалось получить статьи: %s", err)
	}

	var pbArticles []*pb.Article
	for _, a := range articles {
		pbArticles = append(pbArticles, convArticle(a))
	}

	response := &pb.GetManySortedArticlesResponse{
		Articles: pbArticles,
	}
	return response, nil
}

func validateGetArticlesWithAttributeRequest(req *pb.GetArticlesWithAttributeRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateNaturalNum(int(req.GetPageNum())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("page_num", err))
	}
	if err := tools.ValidateNaturalNum(int(req.GetPageSize())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("page_size", err))
	}

	return wrapFeildErrors(errorsFields)
}

func (server *Server) GetArticlesWithAttribute(ctx context.Context, req *pb.GetArticlesWithAttributeRequest) (*pb.GetArticlesWithAttributeResponse, error) {
	if err := validateGetArticlesWithAttributeRequest(req); err != nil {
		return nil, err
	}

	articles, err := server.querier.GetArticlesWithAttribute(ctx, db.GetArticlesWithAttributeParams{
		Title:      req.GetTitle(),
		Text:       req.GetText(),
		Authors:    req.GetAuthors(),
		Evaluation: req.GetEvaluation(),
		Limit:      req.GetPageSize(),
		Offset:     (req.GetPageNum() - 1) * req.GetPageSize(),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "статьи не найдены")
		}
		return nil, status.Errorf(codes.Internal, "не удалось получить статьи: %s", err)
	}

	var pbArticles []*pb.Article
	for _, a := range articles {
		pbArticles = append(pbArticles, convArticle(a))
	}

	response := &pb.GetArticlesWithAttributeResponse{
		Articles: pbArticles,
	}
	return response, nil
}

func validateGetManySortedArticlesWithAttributeRequest(req *pb.GetManySortedArticlesWithAttributeRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateNaturalNum(int(req.GetPageNum())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("page_num", err))
	}
	if err := tools.ValidateNaturalNum(int(req.GetPageSize())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("page_size", err))
	}

	return wrapFeildErrors(errorsFields)
}

func (server *Server) GetManySortedArticlesWithAttribute(ctx context.Context, req *pb.GetManySortedArticlesWithAttributeRequest) (*pb.GetManySortedArticlesWithAttributeResponse, error) {
	if err := validateGetManySortedArticlesWithAttributeRequest(req); err != nil {
		return nil, err
	}

	articles, err := server.querier.GetManySortedArticlesWithAttribute(ctx, db.GetManySortedArticlesWithAttributeParams{
		SelectTitle:      req.GetSelectTitle(),
		SelectText:       req.GetSelectText(),
		SelectEvaluation: req.GetSelectEvaluation(),
		SelectAuthors:    req.GetSelectAuthors(),
		SortedIDArticle:  req.GetSortedIdArticle(),
		SortedEvaluation: req.GetSortedEvaluation(),
		SortedComments:   req.GetSortedComments(),
		SortedAuthors:    req.GetSortedAuthors(),
		SortedTitle:      req.GetSortedTitle(),
		SortedText:       req.GetSortedText(),
		SortedEditedAt:   req.GetSortedEditedAt(),
		SortedCreatedAt:  req.GetSortedCreatedAt(),
		Offset:           (req.GetPageNum() - 1) * req.GetPageSize(),
		Limit:            req.GetPageSize(),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "статьи не найдены")
		}
		return nil, status.Errorf(codes.Internal, "не удалось получить статьи: %s", err)
	}

	var pbArticles []*pb.Article
	for _, a := range articles {
		pbArticles = append(pbArticles, convArticle(a))
	}

	response := &pb.GetManySortedArticlesWithAttributeResponse{
		Articles: pbArticles,
	}
	return response, nil
}
