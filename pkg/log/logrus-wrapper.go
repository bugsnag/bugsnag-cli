package log

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type LogrusLogger struct {
	logger *logrus.Logger
	ctx    context.Context
}

// CustomFormatter is a custom logrus formatter
type CustomFormatter struct{}
type NoAnsiCustomFormatter struct{}

// Format formats the log entry
func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Define colors for different log levels
	var levelColor string
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = "\x1b[37;1m" // White
	case logrus.InfoLevel:
		levelColor = "\x1b[32;1m" // Green
	case logrus.WarnLevel:
		levelColor = "\x1b[33;1m" // Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = "\x1b[31;1m" // Red
	default:
		levelColor = "\x1b[0m" // Reset
	}

	// Customize the log format here
	logMessage := "[" + levelColor + strings.ToUpper(entry.Level.String()) + "\x1b[0m] " + entry.Message + "\n"

	return []byte(logMessage), nil
}

// Format formats the log entry
func (f *NoAnsiCustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Customize the log format here
	logMessage := "[" + strings.ToUpper(entry.Level.String()) + "] " + entry.Message + "\n"

	return []byte(logMessage), nil
}

func NewLogrusLogger(ctx context.Context, verbose bool, noAnsi bool) *LogrusLogger {
	logger := logrus.New()
	logger.Out = os.Stdout

	if noAnsi {
		logger.Formatter = &NoAnsiCustomFormatter{}
	} else {
		logger.Formatter = &CustomFormatter{}

	}

	// Set the log level to debug if verbose is true
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	}

	return &LogrusLogger{logger: logger, ctx: ctx}
}

func (l *LogrusLogger) Debug(msg string) {

	l.logger.WithFields(nil).Debug(msg)
}

func (l *LogrusLogger) Info(msg string) {
	l.addContextCommonFields(nil)

	l.logger.WithFields(nil).Info(msg)
}

func (l *LogrusLogger) Warn(msg string) {
	l.logger.WithFields(nil).Warn(msg)
}

func (l *LogrusLogger) Error(msg string) {
	l.logger.WithFields(nil).Error(msg)
}

func (l *LogrusLogger) Fatal(msg string) {
	l.logger.WithFields(nil).Fatal(msg)
}

func (l *LogrusLogger) addContextCommonFields(fields map[string]interface{}) {
	if l.ctx != nil {
		for k, v := range l.ctx.Value("commonFields").(map[string]interface{}) {
			if _, ok := fields[k]; !ok {
				fields[k] = v
			}
		}
	}
}
