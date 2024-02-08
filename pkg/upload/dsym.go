package upload

import (
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

			// Check if dsymPath exists before proceeding
			if utils.Path(dsymPath).Validate() != nil {
				// TODO: This will be downgraded to a warning with --ignore-missing-dwarf in near future
				log.Error("Could not find dSYM in alternative location: "+utils.DisplayBlankIfEmpty(dsymPath), 1)
			} else {
				log.Info("Using dSYM path: " + dsymPath)
			}

		}

		dsyms, err = ios.GetDsymsForUpload(dsymPath)
		if err != nil {
			return err
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

				// Check if the variables are empty, set if they are and log that we are using setting from the plist file
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
				return err
			} else {
				log.Success("Uploaded dSYM: " + utils.DisplayBlankIfEmpty(dsym.Name))
			}
		}
	}

	return nil
}
