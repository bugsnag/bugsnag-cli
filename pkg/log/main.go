package log

import (
	"context"
)

type LoggerWrapper struct {
	logger Logger
}

func NewLoggerWrapper(ctx context.Context, verbose bool, noAnsi bool) *LoggerWrapper {
	var logger Logger

	logger = NewLogrusLogger(ctx, verbose, noAnsi)

	return &LoggerWrapper{logger: logger}
}

func (lw *LoggerWrapper) Debug(msg string) {
	lw.logger.Debug(msg)
}

func (lw *LoggerWrapper) Info(msg string) {
	lw.logger.Info(msg)
}

func (lw *LoggerWrapper) Warn(msg string) {
	lw.logger.Warn(msg)
}

func (lw *LoggerWrapper) Error(msg string) {
	lw.logger.Error(msg)
}

func (lw *LoggerWrapper) Fatal(msg string) {
	lw.logger.Fatal(msg)
}