package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
)

// ProcessDsym processes a xcarchive, locating its dSYM files and uploading them
// to a Bugsnag server.
//
// Parameters:
// - globalOptions: CLI options provided by the user, including xcarchive settings.
// - endpoint: The server endpoint for uploading dSYM files.
// - logger: Logger instance for logging messages during processing.
//
// Returns:
// - An error if any part of the process fails, otherwise nil.
func ProcessDsym(globalOptions options.CLI, endpoint string, logger log.Logger) error {
	var (
		err              error
		numFilesUploaded int
		xcarchivePath    string
	)

	// Extract dSYM upload options from global options
	dsymOptions := globalOptions.Upload.Dsym

	// Configure globalOptions to use xcarchive as the upload source
	globalOptions.Upload.XcodeArchive = options.XcodeArchive(dsymOptions)

	// Locate the Xcode archive (.xcarchive) path based on the provided options
	xcarchivePath, _ = ios.FindXcarchivePath(globalOptions, logger)

	// If no Xcode archive is found, return an error
	if xcarchivePath != "" {
		// Log the located Xcode archive path
		logger.Info(fmt.Sprintf("Found Xcode archive at %s", xcarchivePath))

		// Process and upload the dSYM files extracted from the Xcode archive
		numFilesUploaded, err = ProcessDsymUpload(xcarchivePath, endpoint, globalOptions, logger)
		if err != nil {
			return err // Return error if the upload fails
		}

		// If files were successfully uploaded from the xcarchive, return
		if numFilesUploaded > 0 {
			return nil
		}
	}

	// Configure globalOptions to use Xcode build directory as the upload source
	globalOptions.Upload.XcodeBuild = options.XcodeBuild(dsymOptions)

	// Attempt to process and upload dSYM files from Xcode build directory
	err = ProcessXcodeBuild(globalOptions, endpoint, logger)
	if err != nil {
		return err // Return immediately if an error occurs
	}

	return nil
}
