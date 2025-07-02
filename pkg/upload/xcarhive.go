package upload

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
)

// ProcessDsymUpload locates and uploads dSYM files from an Xcode archive.
// It searches for dSYM files in the specified Xcode archive, processes them and uploads them.
//
// Parameters:
// - xcarchivePath: The path to the Xcode archive (.xcarchive) containing dSYM files.
// - endpoint: The server endpoint for uploading dSYM files.
// - opts: CLI options containing upload configuration.
// - logger: Logger instance for logging messages during processing.
//
// Returns:
// - The number of dSYM files found and uploaded.
// - An error if any part of the process fails, otherwise nil.
func ProcessDsymUpload(xcarchivePath string, opts options.CLI, logger log.Logger) (int, error) {
	// Locate dSYM files within the specified Xcode archive
	dwarfInfo, tempDir, err := ios.FindDsymsInPath(
		xcarchivePath,
		opts.Upload.XcodeArchive.Shared.IgnoreEmptyDsym,
		opts.Upload.XcodeArchive.Shared.IgnoreMissingDwarf,
		logger,
	)
	// Ensure temporary directory is removed after execution
	defer os.RemoveAll(tempDir)

	if err != nil {
		return 0, fmt.Errorf("error locating dSYM files: %w", err)
	}

	// If no dSYM files are found, return without an error
	if len(dwarfInfo) == 0 {
		return 0, nil
	}

	// Log the number of dSYM files found
	logger.Info(fmt.Sprintf("Found %d dSYM files in %s", len(dwarfInfo), xcarchivePath))

	// Set the project root if not already specified
	opts.Upload.XcodeArchive.Shared.ProjectRoot = ios.GetDefaultProjectRoot(opts.Upload.XcodeArchive.Path[0], opts.Upload.XcodeArchive.Shared.ProjectRoot)
	logger.Info(fmt.Sprintf("Setting `--project-root`: %s", opts.Upload.XcodeArchive.Shared.ProjectRoot))

	// Process and upload the located dSYM files
	err = ios.ProcessDsymUpload(
		filepath.Join(xcarchivePath, "Info.plist"),
		opts.Upload.XcodeArchive.Shared.ProjectRoot,
		opts,
		dwarfInfo,
		logger,
	)
	if err != nil {
		return len(dwarfInfo), fmt.Errorf("error uploading dSYM files: %w", err)
	}

	// Return the number of successfully uploaded dSYM files
	return len(dwarfInfo), nil
}
