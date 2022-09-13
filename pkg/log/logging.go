package log

import (
	"fmt"
	"os"
	"github.com/mattn/go-isatty"
)

const Reset = "\033[0m"
const Red = "\033[31m"
const Green = "\033[32m"
const Yellow = "\033[33m"
const White = "\033[37m"


func logMessage(message string, status string, color string) {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Println("[" + color + status + Reset + "] " + message )
	} else {
		fmt.Println("["+ status +"] " + message )
	}
}

// Error - Displays error message and exits with a status code
func Error(message string, statusCode int)  {
	logMessage(message,"ERROR", Red)
	os.Exit(statusCode)
}

// Warn - Displays warn message
func Warn(message string)  {
	logMessage(message,"WARN", Yellow)
}

// Info - Displays info message
func Info(message string)  {
	logMessage(message,"INFO", White)
}

// Success - Displays success message
func Success(message string)  {
	logMessage(message,"SUCCESS", Green)
}
