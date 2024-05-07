package log

import (
	"context"
)

type LoggerWrapper struct {
	logger Logger
}

const commonFieldsKey string = "commonFields"

func NewLoggerWrapper(verbose bool) *LoggerWrapper {
	// Init Logger
	loggerCommonFields := map[string]interface{}{}
	ctx := context.WithValue(context.Background(), commonFieldsKey, loggerCommonFields)

	logger := NewLogrusLogger(ctx, verbose)

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
