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

func validateGetImageRequest(req *pb.GetImageRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateNaturalNum(int(req.GetIdImage())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("id_image", err))
	}

	return wrapFeildErrors(errorsFields)
}

func (server *Server) GetImage(ctx context.Context, req *pb.GetImageRequest) (*pb.GetImageResponse, error) {
	// Попытка получения данных из кеша
	cacheKey := fmt.Sprintf("image:%d", req.GetIdImage())
	var cachedResponse pb.GetImageResponse
	if err := server.cacher.GetCache(ctx, cacheKey, &cachedResponse); err == nil {
		return &cachedResponse, nil
	}

	// Проверка валидности запроса
	if err := validateGetImageRequest(req); err != nil {
		return nil, err
	}

	// Сам запрос
	image, err := server.querier.GetImage(ctx, req.GetIdImage())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "изображение не найдено")
		}
		return nil, status.Errorf(codes.Internal, "не удалось получить изображение: %s", err)
	}

	response := &pb.GetImageResponse{
		Image: convImage(image),
	}

	if err := server.cacher.SetCache(ctx, cacheKey, response); err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось сохранить данные в кеш: %s", err)
	}

	return response, nil
}

// validateLoadImageRequest проверяет корректность запроса на загрузку изображения
func validateLoadImageRequest(req *pb.LoadImageRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	// Проверка имени изображения
	if err := tools.ValidateString(req.GetName(), 1, 0); err != nil {
		errorsFields = append(errorsFields, fieldViolation("name", err))
	}

	// Проверка наличия содержимого изображения
	if len(req.GetContent()) == 0 {
		errorsFields = append(errorsFields, fieldViolation("content", fmt.Errorf("изображение не должно быть пустым")))
	}

	return wrapFeildErrors(errorsFields)
}

// LoadImage загружает изображение в базу данных
func (server *Server) LoadImage(ctx context.Context, req *pb.LoadImageRequest) (*pb.LoadImageResponse, error) {
	// Проверка валидности запроса
	if err := validateLoadImageRequest(req); err != nil {
		return nil, err
	}

	// Аутентификация пользователя
	accessPayload, err := server.authUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "пользователь не авторизовался: %s", err)
	}

	// Подготовка параметров для вставки в базу данных
	arg := db.LoadImageParams{
		Name:    req.GetName(),
		Content: req.GetContent(),
		IDUser:  accessPayload.IDUser,
	}

	// Вставка изображения в базу данных
	image, err := server.querier.LoadImage(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось загрузить изображение: %s", err)
	}

	response := &pb.LoadImageResponse{
		Image: convImage(image),
	}
	return response, nil
}

func validateEditImageRequest(req *pb.EditImageRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateNaturalNum(int(req.GetIdImage())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("id_image", err))
	}

	if req.Content == nil || len(req.GetContent()) == 0 {
		errorsFields = append(errorsFields, fieldViolation("content", fmt.Errorf("изображение не должно быть пустым")))
	}

	return wrapFeildErrors(errorsFields)
}

func (server *Server) EditImage(ctx context.Context, req *pb.EditImageRequest) (*pb.EditImageResponse, error) {
	if err := validateEditImageRequest(req); err != nil {
		return nil, err
	}

	accessPayload, err := server.authUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "пользователь не авторизован: %s", err)
	}

	// Проверка на авторство
	image, err := server.querier.GetImage(ctx, req.GetIdImage())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось получить изображение: %s", err)
	}

	if image.IDUser != accessPayload.IDUser {
		return nil, status.Errorf(codes.PermissionDenied, "у вас нет прав на редактирование этого изображения")
	}

	arg := db.EditImageParams{
		Content: req.GetContent(),
		IDImage: req.GetIdImage(),
	}

	editedImage, err := server.querier.EditImage(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось изменить изображение: %s", err)
	}

	response := &pb.EditImageResponse{
		Image: convImage(editedImage),
	}
	return response, nil
}

func validateDeleteImageRequest(req *pb.DeleteImageRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateNaturalNum(int(req.GetIdImage())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("id_image", err))
	}

	return wrapFeildErrors(errorsFields)
}

func (server *Server) DeleteImage(ctx context.Context, req *pb.DeleteImageRequest) (*pb.DeleteImageResponse, error) {
	if err := validateDeleteImageRequest(req); err != nil {
		return nil, err
	}

	accessPayload, err := server.authUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "пользователь не авторизован: %s", err)
	}

	// Проверка на авторство
	image, err := server.querier.GetImage(ctx, req.GetIdImage())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось получить изображение: %s", err)
	}

	if image.IDUser != accessPayload.IDUser {
		return nil, status.Errorf(codes.PermissionDenied, "у вас нет прав на удаление этого изображения")
	}

	deletedImage, err := server.querier.DeleteImage(ctx, req.GetIdImage())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось удалить изображение: %s", err)
	}

	response := &pb.DeleteImageResponse{
		Image: convImage(deletedImage),
	}
	return response, nil
}

func validateRenameImageRequest(req *pb.RenameImageRequest) error {
	var errorsFields []*errdetails.BadRequest_FieldViolation

	if err := tools.ValidateNaturalNum(int(req.GetIdImage())); err != nil {
		errorsFields = append(errorsFields, fieldViolation("id_image", err))
	}

	if err := tools.ValidateString(req.GetName(), 1, 0); err != nil {
		errorsFields = append(errorsFields, fieldViolation("name", err))
	}

	return wrapFeildErrors(errorsFields)
}

func (server *Server) RenameImage(ctx context.Context, req *pb.RenameImageRequest) (*pb.RenameImageResponse, error) {
	if err := validateRenameImageRequest(req); err != nil {
		return nil, err
	}

	accessPayload, err := server.authUser(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "пользователь не авторизован: %s", err)
	}

	// Проверка на авторство
	image, err := server.querier.GetImage(ctx, req.GetIdImage())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось получить изображение: %s", err)
	}

	if image.IDUser != accessPayload.IDUser {
		return nil, status.Errorf(codes.PermissionDenied, "у вас нет прав на переименование этого изображения")
	}

	arg := db.RenameImageParams{
		Name:    req.GetName(),
		IDImage: req.GetIdImage(),
	}

	renamedImage, err := server.querier.RenameImage(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "не удалось переименовать изображение: %s", err)
	}

	response := &pb.RenameImageResponse{
		Image: convImage(renamedImage),
	}
	return response, nil
}
