package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type Dsym struct {
	Path               utils.Paths `arg:"" name:"path" help:"The path to the directory or file to upload" type:"path" default:"."`
	IgnoreEmptyDsym    bool        `help:"Throw warnings instead of errors when a dSYM file is found, rather than the expected dSYM directory"`
	IgnoreMissingDwarf bool        `help:"Throw warnings instead of errors when a dSYM with missing DWARF data is found"`
	Plist              utils.Path  `help:"The path to a .plist file from which to obtain build information" type:"path"`
	Scheme             string      `help:"The name of the Xcode options.Scheme used to build the application"`
	VersionName        string      `help:"The version of the application"`
	XcodeProject       utils.Path  `help:"The path to an Xcode project, workspace or containing directory from which to obtain build information" type:"path"`
	ProjectRoot        string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`
}

func ProcessDsym(
	apiKey string,
	options Dsym,
	endpoint string,
	timeout int,
	retries int,
	dryRun bool,
	logger log.Logger,
) error {

	var buildSettings *ios.XcodeBuildSettings
	var plistData *ios.PlistData
	var uploadOptions map[string]string

	var dwarfInfo []*ios.DwarfInfo
	var tempDirs []string
	var dsymPath string
	var err error
	var tempDir string
	xcodeProjPath := string(options.XcodeProject)
	plistPath := string(options.Plist)

	// Performs an automatic cleanup of temporary directories at the end
	defer func() {
		for _, tempDir := range tempDirs {
			_ = os.RemoveAll(tempDir)
		}
	}()

	for _, path := range options.Path {
		if ios.IsPathAnXcodeProjectOrWorkspace(path) {
			if xcodeProjPath == "" {
				xcodeProjPath = path
			}
		} else {
			dsymPath = path
		}

		if xcodeProjPath != "" {
			options.ProjectRoot = ios.GetDefaultProjectRoot(xcodeProjPath, options.ProjectRoot)
			logger.Debug(fmt.Sprintf("Defaulting to '%s' as the project root", options.ProjectRoot))

			// Get build settings and dsymPath

			// If options.Scheme is set explicitly, check if it exists
			if options.Scheme != "" {
				_, err := ios.IsSchemeInPath(xcodeProjPath, options.Scheme)
				if err != nil {
					logger.Warn(err.Error())
				}
			} else {
				// Otherwise, try to find it
				options.Scheme, err = ios.GetDefaultScheme(xcodeProjPath)
				if err != nil {
					logger.Warn(err.Error())
				}
			}

			if options.Scheme != "" {
				buildSettings, err = ios.GetXcodeBuildSettings(xcodeProjPath, options.Scheme)
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

		if options.ProjectRoot == "" {
			return fmt.Errorf("--project-root is required when uploading dSYMs from a directory that is not an Xcode project or workspace")
		}

		if dsymPath == "" {
			return fmt.Errorf("No dSYM locations detected. Please provide a valid dSYM path or an Xcode project/workspace path")
		}

		dwarfInfo, tempDir, err = ios.FindDsymsInPath(dsymPath, options.IgnoreEmptyDsym, options.IgnoreMissingDwarf, logger)
		tempDirs = append(tempDirs, tempDir)
		if err != nil {
			return err
		} else if len(dwarfInfo) == 0 {
			return fmt.Errorf("No dSYM files found in: %s", dsymPath)
		}

		// If the Info.plist path is not defined, we need to build the path to Info.plist from build settings values
		if plistPath == "" && apiKey == "" {
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
		if plistPath != "" && apiKey == "" {
			// Read data from the plist
			plistData, err = ios.GetPlistData(plistPath)
			if err != nil {
				return err
			}

			if apiKey == "" {
				apiKey = plistData.BugsnagProjectDetails.ApiKey
				if apiKey != "" {
					logger.Debug(fmt.Sprintf("Using API key from Info.plist: %s", apiKey))
				}
			}
		}

		for _, dsym := range dwarfInfo {
			dsymInfo := fmt.Sprintf("(UUID: %s, Name: %s, Arch: %s)", dsym.UUID, dsym.Name, dsym.Arch)
			logger.Debug(fmt.Sprintf("Processing dSYM %s", dsymInfo))

			uploadOptions, err = utils.BuildDsymUploadOptions(apiKey, options.ProjectRoot)
			if err != nil {
				return err
			}

			fileFieldData := make(map[string]server.FileField)
			fileFieldData["dsym"] = server.LocalFile(filepath.Join(dsym.Location, dsym.Name))

			err = server.ProcessFileRequest(endpoint+"/dsym", uploadOptions, fileFieldData, timeout, retries, dsym.UUID, dryRun, logger)

			if err != nil {
				if strings.Contains(err.Error(), "404 Not Found") {
					err = server.ProcessFileRequest(endpoint, uploadOptions, fileFieldData, timeout, retries, dsym.UUID, dryRun, logger)
				}
			}

			if err != nil {

				return err
			}
		}
	}

	return nil
}
