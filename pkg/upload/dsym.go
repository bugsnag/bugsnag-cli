package upload

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type Dsym struct {
	VersionName string      `help:"The version of the application."`
	Scheme      string      `help:"The name of the scheme to use when building the application."`
	Dev         bool        `help:"Indicates whether the application is a debug or release build"`
	Plist       string      `help:"Path to the Info.plist file" type:"path"`
	ProjectRoot string      `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	Path        utils.Paths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
}

func ProcessDsym(
	apiKey string,
	versionName string,
	scheme string,
	dev bool,
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
	var possibleSchemeName string
	var schemeExists bool
	var schemeDerivedFrom string

	for _, path := range paths {
		uploadInfo, err := ios.ProcessPathValue(path, projectRoot)
		if err != nil {
			return err
		}

		// If scheme is set explicitly, check if it exists
		if scheme != "" {
			schemeExists, schemeDerivedFrom, err = ios.IsSchemeInPath(path, scheme, projectRoot)
			if err != nil {
				return err
			}
		} else {
			// If the scheme is not set explicitly, try to find it
			possibleSchemeName, schemeDerivedFrom, err = ios.GetDefaultScheme(path, uploadInfo.ProjectRoot)
			if err != nil {
				return err
			}

			schemeExists, schemeDerivedFrom, err = ios.IsSchemeInPath(path, scheme, projectRoot)
			if err != nil {
				return err
			}

			scheme = possibleSchemeName
		}

		if schemeExists {
			log.Info("Using scheme: " + scheme)
		} else {
			log.Info("Unable to determine a scheme using " + schemeDerivedFrom)
		}

		// If the dsymPath is not fed in via <path>
		if uploadInfo.DsymPath == "" {
			buildSettings, err = ios.GetXcodeBuildSettings(path, scheme, projectRoot)
			if err != nil {
				return err
			}

			// Build the dsymPath from build settings
			uploadInfo.DsymPath = buildSettings.ConfigurationBuildDir + "/" + buildSettings.DsymName + "/Contents/Resources/DWARF"

			filesFound, _ := os.ReadDir(uploadInfo.DsymPath)
			switch len(filesFound) {
			case 0:
				return errors.Errorf("No files found in location '%s'", uploadInfo.DsymPath)
			case 1:
				uploadInfo.DsymPath = filepath.Join(uploadInfo.DsymPath, filesFound[0].Name())
			default:
				return errors.Errorf("Multiple files found in location '%s'", uploadInfo.DsymPath)
			}

		}

		if plistPath != "" && (apiKey == "" || versionName == "") {
			// Read data from the plist
			plistData, err := ios.GetPlistData(plistPath)
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

		uploadOptions, err := utils.BuildDsymUploadOptions(apiKey, versionName, dev, uploadInfo.ProjectRoot, overwrite)
		if err != nil {
			return err
		}

		fileFieldData := make(map[string]string)
		fileFieldData["dsym"] = uploadInfo.DsymPath

		err = server.ProcessFileRequest(endpoint+"/dsym", uploadOptions, fileFieldData, timeout, retries, uploadInfo.DsymPath, dryRun)

		if err != nil {
			return err
		} else {
			log.Success("Uploaded " + filepath.Base(uploadInfo.DsymPath))
		}
	}

	return nil
}
