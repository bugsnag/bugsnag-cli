package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
	"path/filepath"
)

// ProcessXcodeArchive processes an Xcode archive, locating its dSYM files and uploading them
// to a Bugsnag server.
//
// Parameters:
// - options: CLI options provided by the user, including Xcode archive settings.
// - endpoint: The server endpoint for uploading dSYM files.
// - logger: Logger instance for logging messages during processing.
//
// Returns:
// - An error if any part of the process fails, otherwise nil.
func ProcessXcodeArchive(options options.CLI, endpoint string, logger log.Logger) error {
	xcarchiveOptions := options.Upload.XcodeArchive
	var (
		xcarchivePath string
		plistPath     string
		dwarfInfo     []*ios.DwarfInfo
		tempDirs      []string
		tempDir       string
		err           error
	)

	// Ensure temporary directories are cleaned up after execution
	defer func() {
		for _, tempDir := range tempDirs {
			_ = os.RemoveAll(tempDir)
		}
	}()

	// Search for an xcarchive in the specified paths
	for _, path := range xcarchiveOptions.Path {
		if filepath.Ext(path) == ".xcarchive" {
			xcarchivePath = path
		} else if utils.IsDir(path) {
			logger.Info(fmt.Sprintf("Searching for Xcode Archives in %s", path))
			if ios.IsPathAnXcodeProjectOrWorkspace(path) {
				xcarchiveOptions.ProjectRoot = ios.GetDefaultProjectRoot(path, xcarchiveOptions.ProjectRoot)
				logger.Info(fmt.Sprintf("Setting `--project-root` from Xcode project settings: %s", xcarchiveOptions.ProjectRoot))

				if xcarchiveOptions.Scheme == "" {
					xcarchiveOptions.Scheme, err = ios.GetDefaultScheme(path)
					if err != nil {
						return fmt.Errorf("Error determining default scheme: %w", err)
					}
				}

				xcarchivePath, err = ios.GetLatestXcodeArchiveForScheme(xcarchiveOptions.Scheme)
				if err != nil {
					return fmt.Errorf("Error locating latest xcarchive: %w", err)
				}
			} else {
				return fmt.Errorf("No xcarchive found in %s", path)
			}
		}

		if xcarchivePath == "" {
			return fmt.Errorf("No xcarchive found in specified paths")
		}

		logger.Info(fmt.Sprintf("Found xcarchive at %s", xcarchivePath))

		// Locate and process dSYM files in the xcarchive
		dwarfInfo, tempDir, err = ios.FindDsymsInPath(
			xcarchivePath,
			xcarchiveOptions.IgnoreEmptyDsym,
			xcarchiveOptions.IgnoreMissingDwarf,
			logger,
		)
		tempDirs = append(tempDirs, tempDir)
		if err != nil {
			return fmt.Errorf("Error locating dSYM files: %w", err)
		}
		if len(dwarfInfo) == 0 {
			return fmt.Errorf("No dSYM files found in: %s", xcarchivePath)
		}
		logger.Info(fmt.Sprintf("Found %d dSYM files in %s", len(dwarfInfo), xcarchivePath))

		// Extract API key from Info.plist if available
		plistPath = filepath.Join(xcarchivePath, "Info.plist")
		err = ios.ProcessDsymUpload(plistPath, endpoint, xcarchiveOptions.ProjectRoot, options, dwarfInfo, logger)
		if err != nil {
			return fmt.Errorf("Error uploading dSYM files: %w", err)
		}
	}
	return nil
}
