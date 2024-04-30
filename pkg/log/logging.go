package log

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"os"
)

const Reset = "\033[0m"
const Red = "\033[31m"
const Green = "\033[32m"
const Yellow = "\033[33m"
const White = "\033[37m"

func LogMessage(message string, status string, color string) {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Printf("[%s%s%s] %s", color, status, Reset, message)
	} else {
		fmt.Printf("[%s] %s", status, message)
	}
}

// Error - Displays error message and exits with a status code
func Error(message string, statusCode int) {
	LogMessage(message, "ERROR", Red)
	os.Exit(statusCode)
}

// Warn - Displays warn message
func Warn(message string) {
	LogMessage(message, "WARN", Yellow)
}

// Info - Displays info message
func Info(message string) {
	LogMessage(message, "INFO", White)
}

// Success - Displays success message
func Success(message string) {
	LogMessage(message, "SUCCESS", Green)
}
