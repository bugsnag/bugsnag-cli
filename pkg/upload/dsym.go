package upload

import (
	"errors"
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

type Dsym struct {
	VersionName        string      `help:"The version of the application."`
	Scheme             string      `help:"The name of the scheme to use when building the application."`
	Dev                bool        `help:"Indicates whether the application is a debug or release build"`
	DsymPath           string      `help:"Path to the dSYM" type:"path"`
	Plist              string      `help:"Path to the Info.plist file" type:"path"`
	ProjectRoot        string      `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	IgnoreMissingDwarf bool        `help:"Throw warnings instead of errors when a dSYM with missing DWARF data is found"`
	IgnoreEmptyDsym    bool        `help:"Throw warnings instead of errors when a *.dSYM file is found, rather than the expected *.dSYM directory"`
	FailOnUpload       bool        `help:"Whether to stop any further uploads if a file fails to upload successfully. By default the command attempts to upload"`
	Path               utils.Paths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
}

func ProcessDsym(
	apiKey string,
	scheme string,
	xcodeProjectPath string,
	plistPath string,
	paths []string,
	failOnUpload bool,
	endpoint string,
	timeout int,
	retries int,
	overwrite bool,
	dryRun bool,
) error {
	var buildSettings *ios.XcodeBuildSettings
	var dsymPath string
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
		if utils.IsDir(path) {
			if xcodeProjectPath == "" {
				xcodeProjectPath, _ = ios.FindProjectOrWorkspaceInPath(path)
			}

			fmt.Println(xcodeProjectPath)

			// If scheme is set explicitly, check if it exists
			if scheme != "" {
				_, err := ios.IsSchemeInPath(xcodeProjectPath, scheme)
				if err != nil {
					log.Warn(err.Error())
				}

			} else {
				// Otherwise, try to find it
				var err error
				scheme, err = ios.GetDefaultScheme(xcodeProjectPath)
				if err != nil {
					log.Warn(err.Error())
				}

			}

			fmt.Println(scheme)

			if scheme != "" {
				var err error
				buildSettings, err = ios.GetXcodeBuildSettings(path, scheme)
				if err != nil {
					return err
				}
			}

			fmt.Println(buildSettings)

			if buildSettings != nil {
				dsymPath = filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.DsymName)

				_, err := os.Stat(dsymPath)
				if err == nil {
					log.Info("Using dSYM path: " + dsymPath)
				}
			}

			fmt.Println(dsymPath)

			if dsymPath != "" {
				var tempDir string
				dwarfInfo, tempDir, _ = ios.FindDsymsInPath(dsymPath, false, false)
				tempDirs = append(tempDirs, tempDir)
			}
			if len(dwarfInfo) == 0 {
				return errors.New("No dSYM files found")
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

			// If the Info.plist path is defined and we still don't know the apiKey or verionName, try to extract them from it
			if plistPath != "" && apiKey == "" {
				// Read data from the plist
				var err error
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

				var err error
				uploadOptions, err = utils.BuildDsymUploadOptions(apiKey, overwrite)
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
					if failOnUpload {
						return err
					} else {
						log.Warn(err.Error())
					}
				} else {
					log.Success("Uploaded dSYM: " + dsym.Name)
				}
			}
		} else {
			return fmt.Errorf("Invalid path: %s", path)
		}
	}

	return nil
}
