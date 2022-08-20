package log

import (
	"fmt"
	"os"
)

const Reset = "\033[0m"
const Red = "\033[31m"
const Green = "\033[32m"
const Yellow = "\033[33m"
const White = "\033[37m"

// Error - Displays error message and exits with a status code
func Error(message string, statusCode int)  {
	fmt.Println("[" + Red + "ERROR" + Reset + "] " + message )
	os.Exit(statusCode)
}

// Warn - Displays warn message
func Warn(message string)  {
	fmt.Println("[" + Yellow + "WARN" + Reset + "] " + message )
}

// Info - Displays info message
func Info(message string)  {
	fmt.Println("[" + White + "INFO" + Reset + "] " + message )
}

// Success - Displays success message
func Success(message string)  {
	fmt.Println("[" + Green + "SUCCESS" + Reset + "] " + message )
}
