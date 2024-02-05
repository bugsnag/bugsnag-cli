package upload

import (
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
		log.Info("Using scheme: " + scheme)

		// If the dsymPath is not fed in via <path>
		if dsymPath == "" {
			buildSettings, err = ios.GetXcodeBuildSettings(path, scheme, projectRoot)
			if err != nil {
				return err
			}

			// Build the dsymPath from build settings
			dsymPath = buildSettings.ConfigurationBuildDir + "/" + buildSettings.DsymName + "/Contents/Resources/DWARF"

		} else {
			// Use the dsymPath from the command line, check if it's zipped and unzip it if necessary
			var extractedLocation string
			extractedLocation, err = utils.ExtractFile(dsymPath, "dsym")
			if err != nil {
				log.Warn(dsymPath + " is not a zip file or directory")
				return err
			}

			if extractedLocation != "" {
				log.Info("Unzipped " + dsymPath + " to " + extractedLocation + " for uploading")
				dsymPath = extractedLocation
			}
		}

		dsyms, err = ios.GetDsymsForUpload(dsymPath)
		if err != nil {
			return err
		}

		for _, dsym := range *dsyms {
			dsymInfo := "(UUID: " + dsym.UUID + ", Name: " + dsym.Name + ", Arch: " + dsym.Arch + ")"
			log.Info("Uploading dSYM " + dsymInfo)
			if plistPath != "" && (apiKey == "" || versionName == "") {
				// Read data from the plist
				plistData, err = ios.GetPlistData(plistPath)
				if err != nil {
					return err
				}

				// Check if the variables are empty, set if they are abd log that we are using setting from the plist file
				if versionName == "" {
					versionName = plistData.VersionName
					log.Info("Using version name from Info.plist: " + versionName)

				}

				if apiKey == "" {
					apiKey = plistData.BugsnagProjectDetails.ApiKey
					log.Info("Using API key from Info.plist: " + apiKey)
				}

			}

			uploadOptions, err = utils.BuildDsymUploadOptions(apiKey, versionName, dev, projectRoot, overwrite)
			if err != nil {
				return err
			}

			fileFieldData := make(map[string]string)
			fileFieldData["dsym"] = dsymPath + "/" + dsym.Name

			err = server.ProcessFileRequest(endpoint+"/dsym", uploadOptions, fileFieldData, timeout, retries, dsym.UUID, dryRun)

			if err != nil {
				return err
			} else {
				log.Success("Uploaded dSYM: " + dsymInfo)
			}
		}
	}

	return nil
}
