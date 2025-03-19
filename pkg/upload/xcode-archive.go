package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
)

// ProcessXcodeArchive locates a Xcode archive (xcarchive), extracts its dSYM files,
// and uploads them to a Bugsnag server.
//
// Parameters:
// - options: CLI options provided by the user, including xcarchive settings.
// - endpoint: The server endpoint for uploading dSYM files.
// - logger: Logger instance for logging messages during processing.
//
// Returns:
// - An error if any part of the process fails, otherwise nil.
func ProcessXcodeArchive(options options.CLI, endpoint string, logger log.Logger) error {
	var (
		xcarchivePath string
		err           error
	)

	// Locate the Xcode archive (.xcarchive) path based on the provided options
	xcarchivePath, err = ios.FindXcarchivePath(options, logger)
	if err != nil {
		return err // Return error if the archive path cannot be determined
	}

	// If no Xcode archive is found, return an error
	if xcarchivePath == "" {
		return fmt.Errorf("no Xcode archive found in specified paths")
	}

	// Log the located Xcode archive path
	logger.Info(fmt.Sprintf("Found Xcode archive at %s", xcarchivePath))

	// Process and upload the dSYM files extracted from the Xcode archive
	_, err = ProcessDsymUpload(xcarchivePath, endpoint, options, logger)
	if err != nil {
		return err // Return error if the upload fails
	}

	return nil // Successfully processed and uploaded dSYM files
}
