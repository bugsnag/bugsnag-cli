package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// ProcessXcodeArchive processes an Xcode archive, locating the archive, its associated dSYM files,
// and uploading them to a BugSnag server.
//
// Parameters:
// - options (options.CLI): The CLI options provided by the user, including Xcode archive settings.
// - endpoint (string): The server endpoint where the dSYM files will be uploaded.
// - logger (log.Logger): The logger used for logging messages during processing.
//
// Returns:
// - error: An error if the process fails at any point, otherwise nil.
func ProcessXcodeArchive(options options.CLI, endpoint string, logger log.Logger) error {
	xcarchiveOptions := options.Upload.XcodeArchive
	var xcarchivePath, plistPath string
	var err error
	var plistData *ios.PlistData
	var dwarfInfo []*ios.DwarfInfo
	var uploadOptions map[string]string
	var tempDirs []string
	var tempDir string

	// Initialize plistPath from shared options if provided
	plistPath = string(xcarchiveOptions.Plist)

	// Ensure temporary directories are cleaned up after execution
	defer func() {
		for _, tempDir := range tempDirs {
			_ = os.RemoveAll(tempDir)
		}
	}()

	// Search for an .xcarchive in the specified paths
	for _, path := range xcarchiveOptions.Path {
		if filepath.Ext(path) == ".xcarchive" {
			// If the path is directly an .xcarchive file, use it
			xcarchivePath = path
		} else if utils.IsDir(path) {
			// If the path is a directory, explore it for an .xcarchive or an Xcode project/workspace
			logger.Info(fmt.Sprintf("Searching for Xcode Archives in %s", path))

			// Check if the directory contains an Xcode project or workspace
			if ios.IsPathAnXcodeProjectOrWorkspace(path) {
				// Set the project root based on Xcode project settings
				xcarchiveOptions.ProjectRoot = ios.GetDefaultProjectRoot(path, xcarchiveOptions.ProjectRoot)
				logger.Info(fmt.Sprintf("Setting `--project-root` from Xcode project settings: %s", xcarchiveOptions.ProjectRoot))

				// Determine the default scheme for the project if not already provided
				if xcarchiveOptions.Scheme == "" {
					xcarchiveOptions.Scheme, err = ios.GetDefaultScheme(path)
					if err != nil {
						return err
					}
				}

				// Attempt to locate the latest .xcarchive associated with the project
				xcarchiveLocation, err := ios.GetXcodeArchiveLocation()
				if err != nil {
					logger.Warn(fmt.Sprintf("Failed to get Xcode archive location: %s", err))
					return err
				}
				xcarchivePath, err = ios.GetLatestXcodeArchive(xcarchiveLocation, xcarchiveOptions.Scheme)
				if err != nil {
					return err
				}

			} else {
				return fmt.Errorf("No xcarchive found in %s", path)
			}
		}
	}

	// If no .xcarchive was found, return an error
	if xcarchivePath == "" {
		return fmt.Errorf("No xcarchive found in specified paths")
	}
	logger.Info(fmt.Sprintf("Found xcarchive at %s", xcarchivePath))

	// Locate and process dSYM files within the .xcarchive
	dwarfInfo, tempDir, err = ios.FindDsymsInPath(
		xcarchivePath,
		xcarchiveOptions.IgnoreEmptyDsym,
		xcarchiveOptions.IgnoreMissingDwarf,
		logger,
	)
	tempDirs = append(tempDirs, tempDir)
	if err != nil {
		return err
	}

	if len(dwarfInfo) == 0 {
		return fmt.Errorf("No dSYM files found in: %s", xcarchivePath)
	}
	logger.Info(fmt.Sprintf("Found %d dSYM files in %s", len(dwarfInfo), xcarchivePath))

	// Extract API key from Info.plist if available and not already set in options
	plistPath = filepath.Join(xcarchivePath, "Info.plist")

	if utils.FileExists(plistPath) && options.ApiKey == "" {
		plistData, err = ios.GetPlistData(plistPath)
		if err != nil {
			return err
		}
		options.ApiKey = plistData.BugsnagProjectDetails.ApiKey
		if options.ApiKey != "" {
			logger.Debug(fmt.Sprintf("Using API key from Info.plist: %s", options.ApiKey))
		}
	}

	// Upload each dSYM file
	for _, dsym := range dwarfInfo {
		dsymInfo := fmt.Sprintf("(UUID: %s, Name: %s, Arch: %s)", dsym.UUID, dsym.Name, dsym.Arch)
		logger.Debug(fmt.Sprintf("Processing dSYM %s", dsymInfo))

		// Build upload options for each dSYM file
		uploadOptions, err = utils.BuildDsymUploadOptions(options.ApiKey, xcarchiveOptions.ProjectRoot)
		if err != nil {
			return err
		}

		// Prepare the file data for uploading
		fileFieldData := map[string]server.FileField{
			"dsym": server.LocalFile(filepath.Join(dsym.Location, dsym.Name)),
		}

		// Attempt to upload the dSYM file to the endpoint
		err = server.ProcessFileRequest(endpoint+"/dsym", uploadOptions, fileFieldData, dsym.UUID, options, logger)
		if err != nil && strings.Contains(err.Error(), "404 Not Found") {
			// If the first upload fails due to 404, retry uploading to the base endpoint
			err = server.ProcessFileRequest(endpoint, uploadOptions, fileFieldData, dsym.UUID, options, logger)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
