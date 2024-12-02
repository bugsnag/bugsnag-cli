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

func ProcessDsym(options options.CLI, endpoint string, logger log.Logger) error {
	dsymOptions := options.Upload.XcodeBuild
	var buildSettings *ios.XcodeBuildSettings
	var plistData *ios.PlistData
	var uploadOptions map[string]string

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

	for _, path := range dsymOptions.Path {
		if ios.IsPathAnXcodeProjectOrWorkspace(path) {
			if xcodeProjPath == "" {
				xcodeProjPath = path
			}
		} else {
			dsymPath = path
		}

		if xcodeProjPath != "" {
			if dsymOptions.ProjectRoot == "" {
				dsymOptions.ProjectRoot = ios.GetDefaultProjectRoot(xcodeProjPath, dsymOptions.ProjectRoot)
				logger.Info(fmt.Sprintf("Setting `--project-root` from Xcode project settings: %s", dsymOptions.ProjectRoot))
			}

			// Get build settings and dsymPath
			// If options.Scheme is set explicitly, check if it exists
			if dsymOptions.Scheme != "" {
				_, err := ios.IsSchemeInPath(xcodeProjPath, dsymOptions.Scheme)
				if err != nil {
					logger.Warn(err.Error())
				}
			} else {
				// Otherwise, try to find it
				dsymOptions.Scheme, err = ios.GetDefaultScheme(xcodeProjPath)
				if err != nil {
					logger.Warn(err.Error())
				}
			}

			if dsymOptions.Scheme != "" {
				buildSettings, err = ios.GetXcodeBuildSettings(xcodeProjPath, dsymOptions.Scheme, options.Upload.XcodeBuild.Configuration)
				if err != nil {
					logger.Warn(err.Error())
				}
			}

			if buildSettings != nil && dsymPath == "" {
				// Build the dsymPath from build settings
				// Which is built up to look like: /Users/Path/To/Config/Build/Dir/MyApp.app.dSYM
				possibleDsymPath := filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.DsymName)

				// Check if dsymPath exists before proceeding
				_, err := os.Stat(possibleDsymPath)
				if err == nil {
					dsymPath = possibleDsymPath
					logger.Debug(fmt.Sprintf("Using dSYM path: %s", dsymPath))
				}
			}
		}

		if dsymOptions.ProjectRoot == "" {
			dsymOptions.ProjectRoot, _ = os.Getwd()
			logger.Info(fmt.Sprintf("Setting `--project-root` to current working directory: %s", dsymOptions.ProjectRoot))
		}

		if dsymPath == "" {
			return fmt.Errorf("No dSYM locations detected. Please provide a valid dSYM path or an Xcode project/workspace path")
		}

		dwarfInfo, tempDir, err = ios.FindDsymsInPath(dsymPath, dsymOptions.IgnoreEmptyDsym, dsymOptions.IgnoreMissingDwarf, logger)
		tempDirs = append(tempDirs, tempDir)
		if err != nil {
			return err
		} else if len(dwarfInfo) == 0 {
			return fmt.Errorf("No dSYM files found in: %s", dsymPath)
		}

		// If the Info.plist path is not defined, we need to build the path to Info.plist from build settings values
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

		// If the Info.plist path is defined and we still don't know the apiKey try to extract them from it
		if plistPath != "" && options.ApiKey == "" {
			// Read data from the plist
			plistData, err = ios.GetPlistData(plistPath)
			if err != nil {
				return err
			}

			if options.ApiKey == "" {
				options.ApiKey = plistData.BugsnagProjectDetails.ApiKey
				if options.ApiKey != "" {
					logger.Debug(fmt.Sprintf("Using API key from Info.plist: %s", options.ApiKey))
				}
			}
		}

		for _, dsym := range dwarfInfo {
			dsymInfo := fmt.Sprintf("(UUID: %s, Name: %s, Arch: %s)", dsym.UUID, dsym.Name, dsym.Arch)
			logger.Debug(fmt.Sprintf("Processing dSYM %s", dsymInfo))

			uploadOptions, err = utils.BuildDsymUploadOptions(options.ApiKey, dsymOptions.ProjectRoot)
			if err != nil {
				return err
			}

			fileFieldData := make(map[string]server.FileField)
			fileFieldData["dsym"] = server.LocalFile(filepath.Join(dsym.Location, dsym.Name))

			err = server.ProcessFileRequest(endpoint+"/dsym", uploadOptions, fileFieldData, dsym.UUID, options, logger)

			if err != nil {
				if strings.Contains(err.Error(), "404 Not Found") {
					err = server.ProcessFileRequest(endpoint, uploadOptions, fileFieldData, dsym.UUID, options, logger)
				}
			}

			if err != nil {

				return err
			}
		}
	}

	return nil
}
