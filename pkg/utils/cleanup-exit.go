package utils

import (
	"os"
	log "unknwon.dev/clog/v2"
)

// CleanupAndExit - Performs cleanup tasks before exiting with a status code
func CleanupAndExit(statusCode int)  {
	// Stop the logger
	log.Stop()

	os.Args = append(os.Args, "--help")

	// Exit with the desired status code.
	os.Exit(statusCode)
}
