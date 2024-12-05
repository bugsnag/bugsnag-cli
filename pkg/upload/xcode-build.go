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

// ProcessXcodeBuild processes an Xcode build, locating the necessary dSYM files and uploading them
// to a Bugsnag server, using the provided Xcode project or workspace configuration.
//
// Parameters:
// - options (options.CLI): The CLI options provided by the user, including Xcode build settings.
// - endpoint (string): The server endpoint where the dSYM files will be uploaded.
// - logger (log.Logger): The logger used for logging messages during processing.
//
// Returns:
// - error: An error if the process fails at any point, otherwise nil.
func ProcessXcodeBuild(options options.CLI, endpoint string, logger log.Logger) error {
	dsymOptions := options.Upload.XcodeBuild
	var buildSettings *ios.XcodeBuildSettings

	var dwarfInfo []*ios.DwarfInfo
	var tempDirs []string
	var dsymPath string
	var err error
	var tempDir string
	xcodeProjPath := string(dsymOptions.XcodeProject)
	plistPath := string(dsymOptions.Plist)

	// Performs an automatic cleanup of temporary directories at the end
	defer func() {
		for _, tempDir := range tempDirs {
			_ = os.RemoveAll(tempDir)
		}
	}()

	// Iterate over the provided paths to locate an Xcode project or workspace
	for _, path := range dsymOptions.Path {
		if ios.IsPathAnXcodeProjectOrWorkspace(path) {
			// If the path is an Xcode project or workspace, use it for further processing
			if xcodeProjPath == "" {
				xcodeProjPath = path
			}
		} else {
			// Otherwise, assume the path is the location of the dSYM file
			dsymPath = path
		}

		// If an Xcode project path is specified, proceed with setting up the project root and build settings
		if xcodeProjPath != "" {
			if dsymOptions.ProjectRoot == "" {
				// Set the project root based on the Xcode project settings if not provided
				dsymOptions.ProjectRoot = ios.GetDefaultProjectRoot(xcodeProjPath, dsymOptions.ProjectRoot)
				logger.Info(fmt.Sprintf("Setting `--project-root` from Xcode project settings: %s", dsymOptions.ProjectRoot))
			}

			// Determine the scheme for the project; check if it's explicitly provided or discover it
			if dsymOptions.Scheme != "" {
				_, err := ios.IsSchemeInPath(xcodeProjPath, dsymOptions.Scheme)
				if err != nil {
					logger.Warn(err.Error())
				}
			} else {
				// If scheme is not provided, try to find the default one
				dsymOptions.Scheme, err = ios.GetDefaultScheme(xcodeProjPath)
				if err != nil {
					logger.Warn(err.Error())
				}
			}

			// Retrieve the Xcode build settings based on the determined scheme and configuration
			if dsymOptions.Scheme != "" {
				buildSettings, err = ios.GetXcodeBuildSettings(xcodeProjPath, dsymOptions.Scheme, options.Upload.XcodeBuild.Configuration)
				if err != nil {
					logger.Warn(err.Error())
				}
			}

			// Build the dSYM path if not already provided using the build settings
			if buildSettings != nil && dsymPath == "" {
				possibleDsymPath := filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.DsymName)

				// Verify if the dSYM file exists before proceeding
				_, err := os.Stat(possibleDsymPath)
				if err == nil {
					dsymPath = possibleDsymPath
					logger.Debug(fmt.Sprintf("Using dSYM path: %s", dsymPath))
				}
			}
		}

		// Set the project root to the current working directory if not already defined
		if dsymOptions.ProjectRoot == "" {
			dsymOptions.ProjectRoot, _ = os.Getwd()
			logger.Info(fmt.Sprintf("Setting `--project-root` to current working directory: %s", dsymOptions.ProjectRoot))
		}

		// If no valid dSYM path is found, return an error
		if dsymPath == "" {
			return fmt.Errorf("No dSYM locations detected. Please provide a valid dSYM path or an Xcode project/workspace path")
		}

		// Locate and process the dSYM files in the specified path
		dwarfInfo, tempDir, err = ios.FindDsymsInPath(dsymPath, dsymOptions.IgnoreEmptyDsym, dsymOptions.IgnoreMissingDwarf, logger)
		tempDirs = append(tempDirs, tempDir)
		if err != nil {
			return err
		} else if len(dwarfInfo) == 0 {
			return fmt.Errorf("No dSYM files found in: %s", dsymPath)
		}

		// If the Info.plist path is not defined, attempt to find it from the build settings
		if plistPath == "" && options.ApiKey == "" {
			if buildSettings != nil {
				plistPathExpected := filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.InfoPlistPath)
				if utils.FileExists(plistPathExpected) {
					plistPath = plistPathExpected
					logger.Debug(fmt.Sprintf("Found Info.plist at expected location: %s", plistPath))
				} else {
					logger.Debug(fmt.Sprintf("No Info.plist found at expected location: %s", plistPathExpected))
				}
			}
		}

		err = ios.ProcessDsymUpload(plistPath, endpoint, dsymOptions.ProjectRoot, options, dwarfInfo, logger)

		if err != nil {
			return err
		}
	}
	return nil
}
