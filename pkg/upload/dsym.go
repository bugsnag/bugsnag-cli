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
	var dsyms *[]*ios.DwarfInfo
	var plistData *ios.PlistData
	var uploadOptions map[string]string

	for _, path := range paths {
		uploadInfo, err := ios.ProcessPathValues(path, dsymPath, projectRoot)
		if err != nil {
			return err
		}

		// If projectRoot is not set explicitly, use uploadInfo to set it (where possible)
		if projectRoot == "" && uploadInfo.ProjectRoot != "" {
			projectRoot = uploadInfo.ProjectRoot
		}

		// If scheme is set explicitly, check if it exists
		if scheme != "" && len(uploadInfo.DsymPaths) == 0 {
			_, err = ios.IsSchemeInPath(path, scheme)
			if err != nil {
				return err
			}
			log.Info("Using scheme: " + scheme)
		} else {

			// Only when the dsym path is not set, try and work out the scheme
			if dsymPath == "" && len(uploadInfo.DsymPaths) == 0 {
				// If the scheme is not set explicitly, try to find it
				scheme, err = ios.GetDefaultScheme(path)
				if err != nil {
					return err
				}
				log.Info("Using scheme: " + scheme)
			}
		}

		// If the dsymPath is not fed in via <path>
		if dsymPath == "" && len(uploadInfo.DsymPaths) == 0 {
			buildSettings, err = ios.GetXcodeBuildSettings(path, scheme)
			if err != nil {
				return err
			}

			// Build the dsymPath from build settings
			// Which is built up to look like: /Users/Path/To/Config/Build/Dir/MyApp.app.dSYM/Contents/Resources/DWARF
			dsymPath = filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.DsymName, "Contents", "Resources", "DWARF")

			// Check if dsymPath exists before proceeding
			if utils.Path(dsymPath).Validate() != nil {
				// TODO: This will be downgraded to a warning with --ignore-missing-dwarf in near future
				log.Error("Could not find dSYM with scheme '"+scheme+"' in expected location: "+utils.DisplayBlankIfEmpty(dsymPath)+"\n\n"+
					"Check that the scheme correlates to the above dSYM location, try re-building your project or specify the dSYM path using --dsym-path", 1)
			} else {
				log.Info("Using dSYM path: " + dsymPath)
			}

		}

		dsyms, err = ios.GetDsymsForUpload(uploadInfo.DsymPaths)
		if err != nil {
			return err
		}

		// If the Info.plist path is not defined, we need to build the path to Info.plist from build settings values
		if plistPath == "" && (apiKey == "" || versionName == "") {
			plistPathExpected := filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.InfoPlistPath)
			if utils.FileExists(plistPathExpected) {
				plistPath = plistPathExpected
				log.Info("Found Info.plist at expected location: " + plistPath)
			} else {
				log.Info("No Info.plist found at expected location: " + plistPathExpected)
			}
		}

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

		for _, dsym := range *dsyms {
			dsymInfo := "(UUID: " + dsym.UUID + ", Name: " + dsym.Name + ", Arch: " + dsym.Arch + ")"
			log.Info("Uploading dSYM " + dsymInfo)

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

			if uploadInfo.TempDir != "" {
				utils.RemoveTempDir(uploadInfo.TempDir)
			}
		}
	}

	return nil
}
