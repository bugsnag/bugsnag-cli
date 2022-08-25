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

// Error - Displays error message and exits with a status code
func Error(message string, statusCode int)  {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Println("[" + Red + "ERROR" + Reset + "] " + message )

	} else {
		fmt.Println("[ERROR] " + message )
	}
	os.Exit(statusCode)
}

// Warn - Displays warn message
func Warn(message string)  {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Println("[" + Yellow + "WARN" + Reset + "] " + message )
	} else {
		fmt.Println("[WARN] " + message )
	}
}

// Info - Displays info message
func Info(message string)  {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Println("[" + White + "INFO" + Reset + "] " + message)
	} else {
		fmt.Println("[INFO] " + message)
	}
}

// Success - Displays success message
func Success(message string)  {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		fmt.Println("[" + Green + "SUCCESS" + Reset + "] " + message )
	} else {
		fmt.Println("[SUCCESS] " + message )
	}
}
