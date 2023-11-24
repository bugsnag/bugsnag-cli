package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type ReactNativeIos struct {
	VersionName   string      `help:"The version of the application."`
	BundleVersion string      `help:"Bundle version for the application. (iOS only)"`
	Scheme        string      `help:"The name of the scheme to use when building the application."`
	SourceMap     string      `help:"Path to the source map file" type:"path"`
	Bundle        string      `help:"Path to the bundle file" type:"path"`
	Plist         string      `help:"Path to the Info.plist file" type:"path"`
	Xcworkspace   string      `help:"Path to the .xcworkspace file" type:"path"`
	CodeBundleID  string      `help:"A unique identifier to identify a code bundle release when using tools like CodePush"`
	Dev           bool        `help:"Indicates whether the application is a debug or release build"`
	ProjectRoot   string      `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	Path          utils.Paths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
}

func ProcessReactNativeIos(
	apiKey string,
	versionName string,
	bundleVersion string,
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

	var rootDirPath string
	var buildSettings *ios.XcodeBuildSettings

	for _, path := range paths {
		// Check/Set the build folder
		buildDirPath := filepath.Join(path, "ios", "build")
		rootDirPath = path
		if !utils.FileExists(buildDirPath) {
			buildDirPath = filepath.Join(path, "build")

			if utils.FileExists(buildDirPath) {
				rootDirPath = filepath.Join(path, "..")

			} else if bundlePath == "" || sourceMapPath == "" {
				return fmt.Errorf("unable to find bundle files or source maps in within " + path)
			}
		}

		// Set a default value for projectRoot if it's not defined
		if projectRoot == "" {
			projectRoot = rootDirPath
		}

		// Check if we're missing any of these parameters
		if bundlePath == "" || plistPath == "" && (apiKey == "" || versionName == "" || bundleVersion == "") {

			// Check to see if we have the xcworkspacePath
			if xcworkspacePath == "" {
				// If not, attempt to locate it in the path/ios/ folder
				if utils.IsDir(filepath.Join(path, "ios")) {
					workspacePath, err := utils.FindFolderWithSuffix(filepath.Join(path, "ios"), ".xcworkspace")
					if err != nil {
						return err
					}
					xcworkspacePath = workspacePath
				} else {
					return fmt.Errorf("unable to find xcworkspace file in your project, please specify the path using --xcworkspace")
				}

			}

			// Check to see if we have a scheme
			if scheme == "" {
				// If not, work it out from the xcworkspace file
				possibleSchemeName := strings.TrimSuffix(filepath.Base(xcworkspacePath), ".xcworkspace")
				schemeExists, err := ios.IsSchemeInWorkspace(xcworkspacePath, possibleSchemeName)
				if err != nil {
					return err
				}

				if schemeExists {
					scheme = possibleSchemeName
				}
			}

			// Pull build settings from the xcworkspace file
			var err error
			buildSettings, err = ios.GetXcodeBuildSettings(xcworkspacePath, scheme)
			if err != nil {
				return err
			}

			// Check to see if we have the Info.Plist path
			if plistPath == "" || !utils.FileExists(plistPath) {
				// If not, we need to build it from build settings values
				plistPathExpected := filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.InfoPlistPath)
				if utils.FileExists(plistPathExpected) {
					plistPath = plistPathExpected
					log.Info("Found Info.plist at: " + plistPath)
				} else {
					log.Info("No Info.plist found at: " + plistPathExpected)
				}
			}
		}

		// Set a sourceMapPath if it's not defined and check that it exists before proceeding
		if sourceMapPath == "" {
			sourceMapFileEnvVar := os.Getenv("SOURCEMAP_FILE")

			// Depending on the value of the SOURCEMAP_FILE environment variable, we will either use the build directory or the value of the environment variable to locate the source map file
			var sourceMapPathToUse string
			if sourceMapFileEnvVar != "" {
				sourceMapPathToUse = sourceMapFileEnvVar
			} else {
				sourceMapPathToUse = buildDirPath
			}

			sourceMapPath = filepath.Join(sourceMapPathToUse, "sourcemaps", "main.jsbundle.map")
			if !utils.FileExists(sourceMapPath) {
				return errors.New("Could not find a suitable source map file, please specify the path by using --source-map")
			}
		}

		// Set a bundlePath if it's not defined and check that it exists before proceeding
		if bundlePath == "" {
			bundleFilePath := filepath.Join(buildSettings.ConfigurationBuildDir, "main.jsbundle")
			if !utils.FileExists(bundleFilePath) {
				return errors.New("Could not find a suitable bundle file, please specify the path by using --bundlePath")
			}
			bundlePath = bundleFilePath
		}

		if plistPath != "" && (apiKey == "" || versionName == "" || bundleVersion == "") {
			// Read data from the plist
			plistData, err := ios.GetPlistData(plistPath)
			if err != nil {
				return err
			}

			// Check if the variables are empty, set if they are abd log that we are using setting from the plist file
			if bundleVersion == "" {
				bundleVersion = plistData.BundleVersion
				log.Info("Using bundle version from Info.plist: " + bundleVersion)
			}

			if versionName == "" {
				versionName = plistData.VersionName
				log.Info("Using version name from Info.plist: " + versionName)

			}

			if apiKey == "" {
				apiKey = plistData.BugsnagProjectDetails.ApiKey
				log.Info("Using API key from Info.plist: " + apiKey)
			}

		}

	}

	uploadOptions, err := utils.BuildReactNativeUploadOptions(apiKey, versionName, bundleVersion, codeBundleId, dev, projectRoot, overwrite, "ios")

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
