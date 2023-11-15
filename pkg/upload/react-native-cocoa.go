package upload

import (
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/bugsnag/bugsnag-cli/pkg/cocoa"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type ReactNativeCocoa struct {
	AppVersion       string      `help:"The version of the application."`
	AppBundleVersion string      `help:"Bundle version for the application. (iOS only)"`
	Scheme           string      `help:"The name of the scheme to use when building the application."`
	SourceMap        string      `help:"Path to the source map file" type:"path"`
	Bundle           string      `help:"Path to the bundle file" type:"path"`
	Plist            string      `help:"Path to the Info.plist file" type:"path"`
	Xcworkspace      string      `help:"Path to the .xcworkspace file" type:"path"`
	CodeBundleID     string      `help:"A unique identifier to identify a code bundle release when using tools like CodePush"`
	Dev              bool        `help:"Indicates whether the application is a debug or release build"`
	ProjectRoot      string      `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	Path             utils.Paths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
}

func ProcessReactNativeCocoa(
	apiKey string,
	appVersion string,
	appBundleVersion string,
	scheme string,
	sourceMapPath string,
	bundlePath string,
	plistPath string,
	xcworkspacePath string,
	codeBundleId string,
	dev bool,
	projectRoot string,
	paths []string,
	endpoint string,
	timeout int,
	retries int,
	overwrite bool,
	dryRun bool,
) error {

	var buildSettings *cocoa.XcodeBuildSettings
	var plistData *cocoa.PlistData

	for _, path := range paths {

		// Set a default value for projectRoot if it's not defined
		if projectRoot == "" {
			projectRoot = path
		}

		// Set a default value for xcworkspacePath if it's not defined
		if xcworkspacePath == "" {
			iosPath := filepath.Join(path, "ios")
			workspacePath, err := utils.FindFolderWithSuffix(iosPath, ".xcworkspace")
			if err != nil {
				return err
			}
			xcworkspacePath = workspacePath
		}

		// Set a sourceMapPath if it's not defined and check that it exists before proceeding
		if sourceMapPath == "" {
			sourceMapPath = filepath.Join(projectRoot, "ios", "build", "sourcemaps", "main.jsbundle.map")
			if !utils.FileExists(sourceMapPath) {
				return errors.New("Could not find a suitable source map file, " +
					"please specify the path by using `--source-map`")
			}
		}

		// If the scheme is not defined, work out what the possible name is and retrieve all xcode schemes based on the xcworkspacePath
		if scheme == "" {
			possibleSchemeName := strings.TrimSuffix(filepath.Base(xcworkspacePath), ".xcworkspace")
			schemeExists, err := cocoa.IsSchemeInWorkspace(xcworkspacePath, possibleSchemeName)
			if err != nil {
				return err
			}

			if schemeExists {
				// We can deduce that possibleSchemeName is the scheme name at this point, and can default to using it's value
				scheme = possibleSchemeName
				buildSettings, err = cocoa.GetXcodeBuildSettings(xcworkspacePath, scheme)
				if err != nil {
					return err
				}

				// Set a default value for bundlePath if it's not defined and check that it exists before proceeding
				if bundlePath == "" {
					bundleFilePath := filepath.Join(buildSettings.ConfigurationBuildDir, "main.jsbundle")
					if !utils.FileExists(bundleFilePath) {
						return errors.New("Could not find a suitable bundle file, " +
							"please specify the path by using `--bundlePath`")
					}
					bundlePath = bundleFilePath
				}

				// Set a default value for plistPath if it's not defined
				if plistPath == "" {
					plistPath = filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.InfoPlistPath)
				}

				// Fetch the plist data with the provided plistPath
				plistData, err = cocoa.GetPlistData(plistPath)
				if err != nil {
					return err
				}

				// Set a default value for the relevant plist data needed for creating the upload options if they aren't already defined
				if appBundleVersion == "" {
					appBundleVersion = plistData.BundleVersion
				}

				if appVersion == "" {
					appVersion = plistData.AppVersion
				}

				if apiKey == "" {
					apiKey = plistData.BugsnagProjectDetails.ApiKey
				}

			} else {
				return errors.New("Could not find a suitable scheme, please specify the scheme by using `--scheme`")
			}

			if err != nil {
				return err
			}
		}

	}

	uploadOptions, err := utils.BuildReactNativeUploadOptions(apiKey, appVersion, appBundleVersion, codeBundleId, dev, projectRoot, overwrite, "ios")

	if err != nil {
		return err
	}

	fileFieldData := make(map[string]string)
	fileFieldData["sourceMap"] = sourceMapPath
	fileFieldData["bundle"] = bundlePath

	err = server.ProcessRequest(endpoint+"/react-native-source-map", uploadOptions, fileFieldData, timeout, sourceMapPath, dryRun)

	if err != nil {
		return err
	} else {
		log.Success("Uploaded " + filepath.Base(sourceMapPath))
	}

	return nil
}
