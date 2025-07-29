package upload

import (
	"fmt"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
)

// ProcessDsym locates dSYM files in an Xcode archive or build directory,
// then uploads them to a Bugsnag server.
//
// Parameters:
// - globalOptions: CLI options including dSYM upload settings.
// - logger: Logger instance for logging progress and errors.
//
// Returns:
// - error if any step fails; otherwise nil.
func ProcessDsym(globalOptions options.CLI, logger log.Logger) error {
	var (
		err              error
		numFilesUploaded int
		xcarchivePath    string
	)

	// Extract dSYM-related upload options
	dsymOptions := globalOptions.Upload.Dsym

	// Set globalOptions to use xcarchive path for dSYM upload
	globalOptions.Upload.XcodeArchive = options.XcodeArchive(dsymOptions)

	// Try to find the .xcarchive path using the provided options
	xcarchivePath, _ = ios.FindXcarchivePath(globalOptions, logger)

	if xcarchivePath != "" {
		logger.Info(fmt.Sprintf("Found Xcode archive at %s", xcarchivePath))

		// Process and upload dSYM files from the archive
		numFilesUploaded, err = ProcessDsymUpload(xcarchivePath, globalOptions, logger)
		if err != nil {
			return err
		}

		// If files were uploaded, no need to continue
		if numFilesUploaded > 0 {
			return nil
		}
	}

	// If no archive found or no files uploaded, fallback to Xcode build directory
	globalOptions.Upload.XcodeBuild = options.XcodeBuild(dsymOptions)

	// Process and upload dSYM files from Xcode build directory
	err = ProcessXcodeBuild(globalOptions, logger)
	if err != nil {
		return err
	}

	return nil
}
