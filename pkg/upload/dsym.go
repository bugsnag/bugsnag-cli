package upload

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type Dsym struct {
	VersionName        string      `help:"The version of the application."`
	Scheme             string      `help:"The name of the scheme to use when building the application."`
	Dev                bool        `help:"Indicates whether the application is a debug or release build"`
	XcodeProject       utils.Path  `help:"Path to the dSYM" type:"path"`
	Plist              utils.Path  `help:"Path to the Info.plist file" type:"path"`
	ProjectRoot        string      `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	IgnoreMissingDwarf bool        `help:"Throw warnings instead of errors when a dSYM with missing DWARF data is found"`
	IgnoreEmptyDsym    bool        `help:"Throw warnings instead of errors when a *.dSYM file is found, rather than the expected *.dSYM directory"`
	Path               utils.Paths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
}

func ProcessDsym(
	apiKey string,
	scheme string,
	xcodeProjPath string,
	plistPath string,
	projectRoot string,
	ignoreMissingDwarf bool,
	ignoreEmptyDsym bool,
	paths []string,
	endpoint string,
	timeout int,
	retries int,
	dryRun bool,
) error {

	var buildSettings *ios.XcodeBuildSettings
	var plistData *ios.PlistData
	var uploadOptions map[string]string

	var dwarfInfo []*ios.DwarfInfo
	var tempDirs []string
	var dsymPath string
	var err error
	var tempDir string

	// Performs an automatic cleanup of temporary directories at the end
	defer func() {
		for _, tempDir := range tempDirs {
			_ = os.RemoveAll(tempDir)
		}
	}()

	for _, path := range paths {
		if ios.IsPathAnXcodeProjectOrWorkspace(path) {
			if xcodeProjPath == "" {
				xcodeProjPath = path
			}
		} else {
			dsymPath = path
		}

		if xcodeProjPath != "" {
			projectRoot = ios.GetDefaultProjectRoot(xcodeProjPath, projectRoot)
			log.Info("Defaulting to '" + projectRoot + "' as the project root")

			// Get build settings and dsymPath

			// If scheme is set explicitly, check if it exists
			if scheme != "" {
				_, err := ios.IsSchemeInPath(xcodeProjPath, scheme)
				if err != nil {
					log.Warn(err.Error())
				}
			} else {
				// Otherwise, try to find it
				scheme, err = ios.GetDefaultScheme(xcodeProjPath)
				if err != nil {
					log.Warn(err.Error())
				}
			}

			if scheme != "" {
				buildSettings, err = ios.GetXcodeBuildSettings(xcodeProjPath, scheme)
				if err != nil {
					log.Warn(err.Error())
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
					log.Info("Using dSYM path: " + dsymPath)
				}
			}
		}

		if projectRoot == "" {
			return errors.New("--project-root is required when uploading dSYMs from a directory that is not an Xcode project or workspace")
		}

		if dsymPath == "" {
			return errors.New("No dSYM locations detected. Please provide a valid dSYM path or an Xcode project/workspace path")
		}

		dwarfInfo, tempDir, err = ios.FindDsymsInPath(dsymPath, ignoreEmptyDsym, ignoreMissingDwarf)
		tempDirs = append(tempDirs, tempDir)
		if err != nil {
			return err
		} else if len(dwarfInfo) == 0 {
			return errors.New("No dSYM files found in: " + dsymPath)
		}

		// If the Info.plist path is not defined, we need to build the path to Info.plist from build settings values
		if plistPath == "" && apiKey == "" {
			if buildSettings != nil {
				plistPathExpected := filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.InfoPlistPath)
				if utils.FileExists(plistPathExpected) {
					plistPath = plistPathExpected
					log.Info("Found Info.plist at expected location: " + plistPath)
				} else {
					log.Info("No Info.plist found at expected location: " + plistPathExpected)
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
					log.Info("Using API key from Info.plist: " + apiKey)
				}
			}
		}

		for _, dsym := range dwarfInfo {
			dsymInfo := "(UUID: " + dsym.UUID + ", Name: " + dsym.Name + ", Arch: " + dsym.Arch + ")"
			log.Info("Uploading dSYM " + dsymInfo)

			uploadOptions, err = utils.BuildDsymUploadOptions(apiKey, projectRoot)
			if err != nil {
				return err
			}

			fileFieldData := make(map[string]string)
			fileFieldData["dsym"] = filepath.Join(dsym.Location, dsym.Name)

			err = server.ProcessFileRequest(endpoint+"/dsym", uploadOptions, fileFieldData, timeout, retries, dsym.UUID, dryRun)

			if err != nil {
				if strings.Contains(err.Error(), "404 Not Found") {
					err = server.ProcessFileRequest(endpoint, uploadOptions, fileFieldData, timeout, retries, dsym.UUID, dryRun)
				}
			}

			if err != nil {

				return err
			}
		}
	}

	return nil
}
