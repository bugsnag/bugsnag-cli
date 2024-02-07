package upload

import (
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type Dsym struct {
	VersionName        string      `help:"The version of the application"`
	Scheme             string      `help:"The name of the scheme to use when building the application"`
	Dev                bool        `help:"Indicates whether the application is a debug or release build"`
	DsymPath           string      `help:"Path to the dSYM" type:"path"`
	Plist              string      `help:"Path to the Info.plist file" type:"path"`
	ProjectRoot        string      `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	IgnoreMissingDwarf bool        `help:"Throw warnings instead of errors when a dSYM with missing DWARF data is found"`
	IgnoreEmptyDsym    bool        `help:"Throw warnings instead of errors when a *.dSYM file is found, rather than the expected *.dSYM directory"`
	FailOnUpload       bool        `help:"Whether to stop any further uploads if a file fails to upload successfully. By default the command attempts to upload all files"`
	Path               utils.Paths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
}

func ProcessDsym(
	apiKey string,
	versionName string,
	scheme string,
	dev bool,
	dsymPath string,
	plistPath string,
	projectRoot string,
	ignoreMissingDwarf bool,
	ignoreEmptyDsym bool,
	failOnUpload bool,
	paths []string,
	endpoint string,
	timeout int,
	retries int,
	overwrite bool,
	dryRun bool,
) error {

	var buildSettings *ios.XcodeBuildSettings
	var dsyms *[]*ios.DsymFile
	var plistData *ios.PlistData
	var uploadOptions map[string]string

	for _, path := range paths {
		uploadInfo, err := ios.ProcessPathValue(path, projectRoot)
		if err != nil {
			return err
		}

		// If dsymPath is not set explicitly, use uploadInfo to set it (if available)
		// uploadInfo.DsymPath is set when <path> is recognised as a dsym path
		if dsymPath == "" && uploadInfo.DsymPath != "" {
			dsymPath = uploadInfo.DsymPath
		}

		// If projectRoot is not set explicitly, use uploadInfo to set it
		if projectRoot == "" {
			projectRoot = uploadInfo.ProjectRoot
		}

		// If scheme is set explicitly, check if it exists
		if scheme != "" {
			_, err = ios.IsSchemeInPath(path, scheme, projectRoot)
			if err != nil {
				return err
			}
		} else {

			// Only when the dsym path is not set, try and work out the scheme
			if dsymPath == "" {
				// If the scheme is not set explicitly, try to find it
				scheme, err = ios.GetDefaultScheme(path, projectRoot)
				if err != nil {
					return err
				}
			}
		}
		log.Info("Using scheme: " + utils.DisplayBlankIfEmpty(scheme))

		// If the dsymPath is not fed in via <path>
		if dsymPath == "" {
			buildSettings, err = ios.GetXcodeBuildSettings(path, scheme, projectRoot)
			if err != nil {
				return err
			}

			// Build the dsymPath from build settings
			// Which is built up to look like: /Users/Path/To/Config/Build/Dir/MyApp.app.dSYM/Contents/Resources/DWARF
			dsymPath = filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.DsymName, "Contents", "Resources", "DWARF")

			// Check if dsymPath exists, if not, try alternative path instead
			if utils.Path(dsymPath).Validate() != nil {

				if ignoreMissingDwarf {
					log.Warn("Could not find dSYM in expected location: " + utils.DisplayBlankIfEmpty(dsymPath))
				} else {
					log.Error("Could not find dSYM in expected location: "+utils.DisplayBlankIfEmpty(dsymPath), 1)
				}

				// Try alternative path which is built up to look like: /Users/Path/To/Config/Build/Dir/MyApp.app
				altDsympath := filepath.Join(buildSettings.ConfigurationBuildDir, strings.TrimSuffix(buildSettings.DsymName, ".dSYM"))

				if utils.Path(dsymPath).Validate() != nil {

					if ignoreMissingDwarf {
						log.Warn("Could not find dSYM in alternative location: " + utils.DisplayBlankIfEmpty(dsymPath))
					} else {
						log.Error("Could not find dSYM in alternative location: "+utils.DisplayBlankIfEmpty(dsymPath), 1)
					}

				} else {
					log.Info("Using alternative dSYM path: " + altDsympath)
					dsymPath = altDsympath
				}

			} else {
				log.Info("Using dSYM path: " + dsymPath)
			}

		}

		dsyms, err = ios.GetDsymsForUpload(dsymPath, ignoreEmptyDsym)
		if err != nil {
			if ignoreMissingDwarf {
				log.Warn(err.Error())
				return nil
			} else {
				return err
			}
		}

		for _, dsym := range *dsyms {
			dsymInfo := "(UUID: " + dsym.UUID + ", Name: " + dsym.Name + ", Arch: " + dsym.Arch + ")"
			log.Info("Uploading dSYM " + dsymInfo)
			// If the Info.plist path is defined and we still don't know the apiKey or verionName, try to extract them from it
			if plistPath != "" && (apiKey == "" || versionName == "") {
				// Read data from the plist
				plistData, err = ios.GetPlistData(plistPath)
				if err != nil {
					return err
				}

				// Check if the variables are empty, set if they are abd log that we are using setting from the plist file
				if versionName == "" {
					versionName = plistData.VersionName
					log.Info("Using version name from Info.plist: " + utils.DisplayBlankIfEmpty(versionName))

				}

				if apiKey == "" {
					apiKey = plistData.BugsnagProjectDetails.ApiKey
					log.Info("Using API key from Info.plist: " + utils.DisplayBlankIfEmpty(apiKey))
				}

			} else if plistPath == "" && (apiKey == "" || versionName == "") {
				// If not, we need to build the path to Info.plist from build settings values
				plistPathExpected := filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.InfoPlistPath)
				if utils.FileExists(plistPathExpected) {
					plistPath = plistPathExpected
					log.Info("Found Info.plist at expected location: " + plistPath)
				} else {
					log.Info("No Info.plist found at expected location: " + plistPathExpected)
				}
			}

			uploadOptions, err = utils.BuildDsymUploadOptions(apiKey, versionName, dev, projectRoot, overwrite)
			if err != nil {
				return err
			}

			fileFieldData := make(map[string]string)
			fileFieldData["dsym"] = filepath.Join(dsymPath, dsym.Name)

			err = server.ProcessFileRequest(filepath.Join(endpoint, "dsym"), uploadOptions, fileFieldData, timeout, retries, dsym.UUID, dryRun)

			if err != nil {
				if failOnUpload {
					return err
				} else {
					log.Warn(err.Error())
				}
			} else {
				log.Success("Uploaded dSYM: " + utils.DisplayBlankIfEmpty(dsym.Name))
			}
		}
	}

	return nil
}
