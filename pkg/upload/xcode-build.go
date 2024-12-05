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

// ProcessXcodeBuild processes an Xcode build, locates necessary dSYM files, and uploads them
// to a Bugsnag server using the provided Xcode project or workspace configuration.
//
// Parameters:
// - options: CLI options provided by the user, including Xcode build settings.
// - endpoint: The server endpoint for uploading dSYM files.
// - logger: Logger instance for logging messages during processing.
//
// Returns:
// - An error if any part of the process fails, otherwise nil.
func ProcessXcodeBuild(options options.CLI, endpoint string, logger log.Logger) error {
	dsymOptions := options.Upload.XcodeBuild
	var (
		buildSettings *ios.XcodeBuildSettings
		dwarfInfo     []*ios.DwarfInfo
		tempDirs      []string
		dsymPath      string
		tempDir       string
		err           error
	)
	xcodeProjPath := string(dsymOptions.XcodeProject)
	plistPath := string(dsymOptions.Plist)

	// Cleanup temporary directories on exit
	defer func() {
		for _, tempDir := range tempDirs {
			_ = os.RemoveAll(tempDir)
		}
	}()

	// Process paths provided in the CLI options
	for _, path := range dsymOptions.Path {
		if filepath.Ext(path) == ".xcarchive" {
			logger.Warn(fmt.Sprintf("The specified path %s is an .xcarchive. Please use the `xcode-archive` command instead as this functionality will be deprecated in future releases.", path))
		}

		if ios.IsPathAnXcodeProjectOrWorkspace(path) {
			// Use the first valid Xcode project/workspace path
			if xcodeProjPath == "" {
				xcodeProjPath = path
			}
		} else {
			// Assume the path is a dSYM file location
			dsymPath = path
		}

		if xcodeProjPath != "" {
			// Determine project root if not provided
			if dsymOptions.ProjectRoot == "" {
				dsymOptions.ProjectRoot = ios.GetDefaultProjectRoot(xcodeProjPath, dsymOptions.ProjectRoot)
				logger.Info(fmt.Sprintf("Setting `--project-root` from Xcode project settings: %s", dsymOptions.ProjectRoot))
			}

			// Determine or validate the scheme
			if dsymOptions.Scheme == "" {
				dsymOptions.Scheme, err = ios.GetDefaultScheme(xcodeProjPath)
				if err != nil {
					logger.Warn(fmt.Sprintf("Error determining default scheme: %s", err))
				}
			} else {
				_, err = ios.IsSchemeInPath(xcodeProjPath, dsymOptions.Scheme)
				if err != nil {
					logger.Warn(fmt.Sprintf("Scheme validation error: %s", err))
				}
			}

			// Retrieve build settings for the scheme and configuration
			if dsymOptions.Scheme != "" {
				buildSettings, err = ios.GetXcodeBuildSettings(xcodeProjPath, dsymOptions.Scheme, dsymOptions.Configuration)
				if err != nil {
					logger.Warn(fmt.Sprintf("Error retrieving build settings: %s", err))
				}
			}

			// Construct the dSYM path if not already specified
			if buildSettings != nil && dsymPath == "" {
				possibleDsymPath := filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.DsymName)
				if _, err = os.Stat(possibleDsymPath); err == nil {
					dsymPath = possibleDsymPath
					logger.Debug(fmt.Sprintf("Using dSYM path: %s", dsymPath))
				}
			}
		}

		// Default project root to current directory if not set
		if dsymOptions.ProjectRoot == "" {
			dsymOptions.ProjectRoot, _ = os.Getwd()
			logger.Info(fmt.Sprintf("Setting `--project-root` to current working directory: %s", dsymOptions.ProjectRoot))
		}

		// Validate dSYM path
		if dsymPath == "" {
			return fmt.Errorf("No dSYM locations detected. Provide a valid dSYM path or Xcode project/workspace path")
		}

		// Locate and process dSYM files
		dwarfInfo, tempDir, err = ios.FindDsymsInPath(dsymPath, dsymOptions.IgnoreEmptyDsym, dsymOptions.IgnoreMissingDwarf, logger)
		tempDirs = append(tempDirs, tempDir)
		if err != nil {
			return fmt.Errorf("Error locating dSYM files: %w", err)
		}
		if len(dwarfInfo) == 0 {
			return fmt.Errorf("No dSYM files found in: %s", dsymPath)
		}

		// Locate Info.plist if not already specified
		if plistPath == "" && options.ApiKey == "" && buildSettings != nil {
			plistPathExpected := filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.InfoPlistPath)
			if utils.FileExists(plistPathExpected) {
				plistPath = plistPathExpected
				logger.Debug(fmt.Sprintf("Found Info.plist at: %s", plistPath))
			} else {
				logger.Debug(fmt.Sprintf("Info.plist not found at: %s", plistPathExpected))
			}
		}

		// Upload dSYM files
		err = ios.ProcessDsymUpload(plistPath, endpoint, dsymOptions.ProjectRoot, options, dwarfInfo, logger)
		if err != nil {
			return fmt.Errorf("Error uploading dSYM files: %w", err)
		}
	}

	return nil
}
