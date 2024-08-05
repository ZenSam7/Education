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

func validateCreateCommentRequest(req *pb.CreateCommentRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateString(req.GetText(), 1, 1000); err != nil {
		errorsFields = append(errorsFields, fieldViolation("text", err))
	}
	if err := tools.ValidateNaturalNum(int(req.GetIdArticle())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("id_article", err))
	}
	if err := tools.ValidateNaturalNum(int(req.GetAuthor())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("author", err))
	}

	return wrapFeildErrors(errorsFields)
}

func (server *Server) CreateComment(ctx context.Context, req *pb.CreateCommentRequest) (*pb.CreateCommentResponse, error) {
	if err := validateCreateCommentRequest(req); err != nil {
		return nil, err
	}

	// Может создать только авторизованный пользователь
	info, err := server.extractMetadata(ctx)
	if err != nil {
		return nil, err
	}

	payload, err := server.tokenMaker.VerifyToken(info.AccessToken)
	if err != nil {
		return nil, err
	}

	comment, err := server.querier.CreateComment(ctx, db.CreateCommentParams{
		Author:    payload.IDUser,
		Text:      req.GetText(),
		IDArticle: req.GetIdArticle(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateCommentResponse{Comment: convComment(comment)}, nil
}

func validateGetCommentRequest(req *pb.GetCommentRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateNaturalNum(int(req.GetIdComment())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("id_comment", err))
	}

	return wrapFeildErrors(errorsFields)
}
func (server *Server) GetComment(ctx context.Context, req *pb.GetCommentRequest) (*pb.GetCommentResponse, error) {
	if err := validateGetCommentRequest(req); err != nil {
		return nil, err
	}

	comment, err := server.querier.GetComment(ctx, req.GetIdComment())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "комментария с id %d не существует", req.GetIdComment())
		}

		return nil, err
	}

	return &pb.GetCommentResponse{Comment: convComment(comment)}, nil
}

func validateGetCommentsOfArticleRequest(req *pb.GetCommentsOfArticleRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateNaturalNum(int(req.GetIdArticle())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("id_article", err))
	}
	if err := tools.ValidateNaturalNum(int(req.GetPageNum())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("page_num", err))
	}
	if err := tools.ValidateNaturalNum(int(req.GetPageSize())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("page_size", err))
	}

	return wrapFeildErrors(errorsFields)
}
func (server *Server) GetCommentsOfArticle(ctx context.Context, req *pb.GetCommentsOfArticleRequest) (*pb.GetCommentsOfArticleResponse, error) {
	if err := validateGetCommentsOfArticleRequest(req); err != nil {
		return nil, err
	}

	comments, err := server.querier.GetCommentsOfArticle(ctx, db.GetCommentsOfArticleParams{
		IDArticle: req.GetIdArticle(),
		Offset:    (req.GetPageNum() - 1) * req.GetPageSize(),
		Limit:     req.GetPageSize(),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetCommentsOfArticleResponse{Comments: []*pb.Comment{}}, nil
		}

		return nil, fmt.Errorf("не удалось получить комментарии к статье: %s", err)
	}

	var response pb.GetCommentsOfArticleResponse
	for _, comm := range comments {
		response.Comments = append(response.Comments, convComment(comm))
	}

	return &response, nil
}

func validateDeleteCommentRequest(req *pb.DeleteCommentRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateNaturalNum(int(req.GetIdComment())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("id_comment", err))
	}

	return wrapFeildErrors(errorsFields)
}
func (server *Server) DeleteComment(ctx context.Context, req *pb.DeleteCommentRequest) (*pb.DeleteCommentResponse, error) {
	if err := validateDeleteCommentRequest(req); err != nil {
		return nil, err
	}

	// Только если авторизован автор коммента
	info, err := server.extractMetadata(ctx)
	if err != nil {
		return nil, err
	}

	payload, err := server.tokenMaker.VerifyToken(info.AccessToken)
	if err != nil {
		return nil, err
	}

	if payload.IDUser != req.GetIdComment() {
		return nil, status.Errorf(codes.PermissionDenied, "только автор может удалить комментарий")
	}

	_, err = server.querier.DeleteComment(ctx, req.GetIdComment())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "комментария с id %d не существует", req.GetIdComment())
		}
		return nil, err
	}

	return &pb.DeleteCommentResponse{}, nil
}

func validateEditCommentRequest(req *pb.EditCommentRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateNaturalNum(int(req.GetIdComment())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("id_comment", err))
	}
	if err := tools.ValidateNaturalNum(int(req.GetEvaluation())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("evaluation", err))
	}
	if err := tools.ValidateString(req.GetText(), 1, 1000); err != nil {
		errorsFields = append(errorsFields, fieldViolation("text", err))
	}

	return wrapFeildErrors(errorsFields)
}
func (server *Server) EditComment(ctx context.Context, req *pb.EditCommentRequest) (*pb.EditCommentResponse, error) {
	if err := validateEditCommentRequest(req); err != nil {
		return nil, err
	}

	// Только если авторизован автор коммента
	info, err := server.extractMetadata(ctx)
	if err != nil {
		return nil, err
	}

	payload, err := server.tokenMaker.VerifyToken(info.AccessToken)
	if err != nil {
		return nil, err
	}

	if payload.IDUser != req.GetIdComment() {
		return nil, status.Errorf(codes.PermissionDenied, "только автор может редактировать комментарий")
	}

	_, err = server.querier.EditComment(ctx, db.EditCommentParams{
		IDComment:  req.GetIdComment(),
		Evaluation: req.GetEvaluation(),
		Text:       req.GetText(),
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "комментария с id %d не существует", req.GetIdComment())
		}
		return nil, err
	}

	return &pb.EditCommentResponse{}, nil
}
