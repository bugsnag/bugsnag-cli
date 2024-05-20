package log

type LoggerWrapper struct {
	logger Logger
}

func NewLoggerWrapper(verbose bool, logLevel string) *LoggerWrapper {
	logger := NewLogrusLogger(verbose, logLevel)

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
