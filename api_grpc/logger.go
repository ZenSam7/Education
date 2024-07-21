package api_grpc

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	methodName := fullMethod[len(fullMethod)-1]

	statusCode := codes.Unknown
	if sc, ok := status.FromError(err); ok {
		statusCode = sc.Code()
	}

	var msgType *zerolog.Event
	if err != nil {
		msgType = sendLog.Error().Err(err)
	} else {
		msgType = sendLog.Info()
	}

	msgType.Dur("duration", duration).
		Int("code", int(statusCode)).
		Str("status", statusCode.String()).
		Msgf("%s |", methodName)

	return
}
