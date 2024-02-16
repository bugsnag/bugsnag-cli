package upload

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type Dsym struct {
	VersionName string      `help:"The version of the application."`
	Scheme      string      `help:"The name of the scheme to use when building the application."`
	Dev         bool        `help:"Indicates whether the application is a debug or release build"`
	DsymPath    string      `help:"Path to the dSYM" type:"path"`
	Plist       string      `help:"Path to the Info.plist file" type:"path"`
	ProjectRoot string      `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	Path        utils.Paths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
}

func ProcessDsym(
	apiKey string,
	versionName string,
	scheme string,
	dev bool,
	dsymPath string,
	plistPath string,
	projectRoot string,
	paths []string,
	endpoint string,
	timeout int,
	retries int,
	overwrite bool,
	dryRun bool,
) error {

	var buildSettings *ios.XcodeBuildSettings
	var plistData *ios.PlistData
	var uploadOptions map[string]string

	var dwarfInfo []*ios.DwarfInfo
	var tempDirs []string

	// Performs an automatic cleanup of temporary directories at the end
	defer func() {
		for _, tempDir := range tempDirs {
			_ = os.RemoveAll(tempDir)
		}
	}()

	for _, path := range paths {

		if path == "" {
			// set path to current directory if not set
			path, _ = os.Getwd()
		}

		if ios.IsPathAnXcodeProjectOrWorkspace(path) {
			projectRoot = ios.GetDefaultProjectRoot(path, projectRoot)
			log.Info("Defaulting to '" + projectRoot + "' as the project root")

			// Get build settings and dsymPath

			// If scheme is set explicitly, check if it exists
			if scheme != "" {
				_, err := ios.IsSchemeInPath(path, scheme)
				if err != nil {
					log.Warn(err.Error())
				}

			} else {
				// Otherwise, try to find it
				var err error
				scheme, err = ios.GetDefaultScheme(path)
				if err != nil {
					log.Warn(err.Error())
				}

			}

			var err error
			buildSettings, err = ios.GetXcodeBuildSettings(path, scheme)
			if err != nil {
				return err
			}

			if dsymPath == "" {
				// Build the dsymPath from build settings
				// Which is built up to look like: /Users/Path/To/Config/Build/Dir/MyApp.app.dSYM
				dsymPath = filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.DsymName)

				// Check if dsymPath exists before proceeding
				if utils.Path(dsymPath).Validate() != nil {
					// TODO: This will be toggled between Error and Warn with --ignore-missing-dwarf in near future
					log.Error("Could not find dSYM with scheme '" + scheme + "' in expected location: " + utils.DisplayBlankIfEmpty(dsymPath) + "\n\n" +
						"Check that the scheme correlates to the above dSYM location, try re-building your project or specify the dSYM path using --dsym-path\n", 1)
				} else {
					log.Info("Using dSYM path: " + dsymPath)
				}
			}

		} else if dsymPath == "" {
			log.Info("No Xcode project, workspace or package in '" + path + "'")
			dsymPath = path
		}

		if dsymPath != "" {
			var tempDir string
			dwarfInfo, tempDir, _ = ios.FindDsymsInPath(dsymPath)
			if len(dwarfInfo) > 0 && projectRoot == "" {
				return errors.New("--project-root is required when uploading dSYMs from a directory that is not an Xcode project or workspace")
			}
			tempDirs = append(tempDirs, tempDir)
		}

		if len(dwarfInfo) == 0 {
			return errors.New("No dSYM files found in expected locations '" + dsymPath + "' and '" + path + "'")
		}

		// If the Info.plist path is not defined, we need to build the path to Info.plist from build settings values
		if plistPath == "" && (apiKey == "" || versionName == "") {
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

		// If the Info.plist path is defined and we still don't know the apiKey or verionName, try to extract them from it
		if plistPath != "" && (apiKey == "" || versionName == "") {
			// Read data from the plist
			var err error
			plistData, err = ios.GetPlistData(plistPath)
			if err != nil {
				return err
			}

			// Check if the variables are empty, set if they are and log that we are using setting from the plist file
			if versionName == "" {
				versionName = plistData.VersionName
				if versionName != "" {
					log.Info("Using version name from Info.plist: " + versionName)
				}
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

			var err error
			uploadOptions, err = utils.BuildDsymUploadOptions(apiKey, versionName, dev, projectRoot, overwrite)
			if err != nil {
				return err
			}

			fileFieldData := make(map[string]string)
			fileFieldData["dsym"] = filepath.Join(dsym.Location, dsym.Name)

			err = server.ProcessFileRequest(endpoint+"/dsym", uploadOptions, fileFieldData, timeout, retries, dsym.UUID, dryRun)

			if err != nil {
				return err
			} else {
				log.Success("Uploaded dSYM: " + dsym.Name)
			}

		}
	}

	return nil
}
