package api_grpc

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

var sendLog zerolog.Logger

// SettingLogger Настраиваем и создаём логгер
func SettingLogger() {
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

	sendLog = zerolog.New(output).With().Timestamp().Logger()
}

func GrpcLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp interface{}, err error) {
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

	msgType := sendLog.Info()
	if err != nil {
		msgType = sendLog.Error().Err(err)
	}

	msgType.
		Int("duration", int(duration.Milliseconds())).
		Int("code", int(statusCode)).
		Str("status", statusCode.String()).
		Str("protocol", "grpc").
		Msgf("%s |", calledFunc)

	return
}

// Response откладываем внутреннюю информацию ответа в кастомную структуру, чтобы её можно было достать
type Response struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

// WriteHeader перехватываем хедер и записываем к себе код статуса
func (r *Response) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *Response) Write(body []byte) (int, error) {
	r.Body = body
	return r.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Засекаем время и откладываем код статуса к себе в кармашек (в Response.StatusCode)
		startTime := time.Now()
		recorder := &Response{ResponseWriter: res, StatusCode: http.StatusOK}
		handler.ServeHTTP(recorder, req)
		duration := time.Since(startTime)

		fullMethod := strings.Split(req.RequestURI, "/")
		calledFunc := fullMethod[len(fullMethod)-1]

		msgType := sendLog.Info()
		if recorder.StatusCode != http.StatusOK {
			msgType = sendLog.Error().Bytes("body", recorder.Body)
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
