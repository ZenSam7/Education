package worker

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct {
	asynq.Logger
}

func NewWorkerLogger() *Logger {
	return &Logger{}
}

func (l *Logger) PrintLog(level zerolog.Level, args ...interface{}) {
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (l *Logger) Debug(args ...interface{}) {
	l.PrintLog(zerolog.DebugLevel, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.PrintLog(zerolog.InfoLevel, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.PrintLog(zerolog.WarnLevel, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.PrintLog(zerolog.ErrorLevel, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.PrintLog(zerolog.FatalLevel, args...)
}
