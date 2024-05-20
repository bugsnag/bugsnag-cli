package log

import (
	"github.com/mattn/go-isatty"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

type LogrusLogger struct {
	logger *logrus.Logger
}

// CustomFormatter is a custom logrus formatter
type CustomFormatter struct{}
type NoAnsiCustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := ""

	// Check if output is a TTY
	isTerminal := isatty.IsTerminal(os.Stdout.Fd())

	output += "["

	// Add level with color if output is a TTY
	if isTerminal {
		switch entry.Level {
		case logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel:
			output += "\x1b[31m" // red
		case logrus.WarnLevel:
			output += "\x1b[33m" // yellow
		case logrus.InfoLevel:
			output += "\x1b[36m" // cyan
		case logrus.DebugLevel, logrus.TraceLevel:
			output += "\x1b[32m" // green
		}
	}

	output += strings.ToUpper(entry.Level.String())

	// Reset color if output is a TTY
	if isTerminal {
		output += "\x1b[0m"
	}

	output += "] " + entry.Message + "\n"

	return []byte(output), nil
}

func NewLogrusLogger(verbose bool, logLevel string) *LogrusLogger {
	logger := logrus.New()
	logger.Out = os.Stdout
	logger.Formatter = &CustomFormatter{}

	// Set the log level to debug if verbose is true or logLevel is set
	if verbose {
		logger.SetLevel(logrus.DebugLevel)
	} else if logLevel != "" {
		level, err := logrus.ParseLevel(logLevel)
		if err != nil {
			logger.Fatal("Invalid log level")
		}
		logger.SetLevel(level)
	}

	return &LogrusLogger{logger: logger}
}

func (l *LogrusLogger) Debug(msg string) {

	l.logger.WithFields(nil).Debug(msg)
}

func (l *LogrusLogger) Info(msg string) {
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
