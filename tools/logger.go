package tools

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"os"
	"strings"
	"time"
)

var Log zerolog.Logger

// MakeLogger Настраиваем и создаём логгер
func MakeLogger() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC822}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %s |", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	Log = zerolog.New(output).With().Timestamp().Logger()
}

func GrpcLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// Засекаем время
	startTime := time.Now()
	resp, err = handler(ctx, req)
	duration := time.Since(startTime)

	fullMethod := strings.Split(info.FullMethod, "/")
	calledFunc := fullMethod[len(fullMethod)-1]

	statusCode := codes.Unknown
	if sc, ok := status.FromError(err); ok {
		statusCode = sc.Code()
	}

	msgType := Log.Info()
	if err != nil {
		msgType = Log.Error().Err(err)
	}

	msgType.
		Int("duration", int(duration.Milliseconds())).
		Int("code", int(statusCode)).
		Str("status", statusCode.String()).
		Str("protocol", "grpc").
		Msgf("%s |", calledFunc)

	return
}

// response откладываем внутреннюю информацию ответа в кастомную структуру, чтобы её можно было достать
type response struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

// WriteHeader перехватываем хедер и записываем к себе код статуса
func (r *response) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *response) Write(body []byte) (int, error) {
	r.Body = body
	return r.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Засекаем время и откладываем код статуса к себе в кармашек (в response.StatusCode)
		startTime := time.Now()
		recorder := &response{ResponseWriter: res, StatusCode: http.StatusOK}
		handler.ServeHTTP(recorder, req)
		duration := time.Since(startTime)

		fullMethod := strings.Split(req.RequestURI, "/")
		calledFunc := fullMethod[len(fullMethod)-1]

		msgType := Log.Info()
		if recorder.StatusCode != http.StatusOK {
			msgType = Log.Error().Bytes("body", recorder.Body)
		}

		msgType.
			Int("duration", int(duration.Milliseconds())).
			Str("method", req.Method).
			Str("protocol", "http").
			Int("code", recorder.StatusCode).
			Str("status", http.StatusText(recorder.StatusCode)).
			Msgf("%s |", calledFunc)
	})
}
