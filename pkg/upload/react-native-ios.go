package upload

import (
	"fmt"
	"os"
	"path/filepath"

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

		// Attempt to parse information from the .xcworkspace file if values aren't provided on the command line
		if bundlePath == "" || (plistPath == "" && (apiKey == "" || versionName == "" || bundleVersion == "")) {

			// Validate workspacePath (if provided) or attempt to find one
			if xcworkspacePath != "" {
				if !utils.FileExists(xcworkspacePath) {
					return errors.New("unable to find the specified .xcworkspace file: " + xcworkspacePath)
				}
			} else {
				// Attempt to locate .xcworkspace in the path/ios/ folder
				if utils.IsDir(filepath.Join(path, "ios")) {
					xcworkspacePath, _ = utils.FindFolderWithSuffix(filepath.Join(path, "ios"), ".xcworkspace")
					if xcworkspacePath != "" {
						log.Info("Found .xcworkspace at: " + xcworkspacePath)
					} else {
						log.Info("No .xcworkspace found at: " + xcworkspacePath)
					}
				}
			}

			// Validate the scheme name (if provided) or attempt to find one in the workspace
			if xcworkspacePath != "" {
				if scheme != "" {
					_, err := ios.IsSchemeInWorkspace(xcworkspacePath, scheme, projectRoot)
					if err != nil {
						return err
					}
				} else {
					// If not, work it out from the xcworkspace file
					possibleSchemeName, err := ios.GetDefaultScheme(xcworkspacePath, projectRoot)
					if err != nil {
						return err
					}

					schemeExists, _ := ios.IsSchemeInWorkspace(xcworkspacePath, possibleSchemeName, projectRoot)
					if schemeExists {
						scheme = possibleSchemeName
						log.Info("Using scheme from .xcworkspace: " + scheme)
					} else {
						log.Info("Unable to determine a scheme from .xcworkspace: " + xcworkspacePath)
					}
				}

				// Pull build settings from the xcworkspace file
				if scheme != "" {
					var err error
					buildSettings, err = ios.GetXcodeBuildSettings(xcworkspacePath, scheme, projectRoot)
					if err != nil {
						log.Warn("Unable to read build settings for scheme " + scheme + " from " + xcworkspacePath)
					}
				}
			}

			// Check to see if we have the Info.Plist path
			if plistPath != "" {
				if !utils.FileExists(plistPath) {
					return errors.New("unable to find specified Info.plist file: " + plistPath)
				}
			} else if buildSettings != nil {
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

		// Check that the bundle file exists and error out if it doesn't
		if bundlePath != "" {
			if !utils.FileExists(bundlePath) {
				return errors.New("unable to find specified bundle file: " + bundlePath)
			}
		} else {
			// Set a bundlePath if it's not defined and check that it exists before proceeding
			if buildSettings != nil {
				possibleBundleFilePath := filepath.Join(buildSettings.ConfigurationBuildDir, "main.jsbundle")
				if utils.FileExists(possibleBundleFilePath) {
					bundlePath = possibleBundleFilePath
					log.Info("Found bundle file at: " + bundlePath)
				} else {
					log.Info("No bundle file found at: " + possibleBundleFilePath)
				}
			}
		}

		// Check that we now have a bundle path
		if bundlePath == "" {
			return errors.New("Could not find a bundle file, please specify the path by using --bundle-path")
		}

		// Check that the source map file exists and error out if it doesn't
		if sourceMapPath != "" {
			if !utils.FileExists(sourceMapPath) {
				return errors.New("unable to find specified source map: " + sourceMapPath)
			}
		} else {
			// Use SOURCEMAP_FILE environment variable, if defined, or use the build directory
			sourceMapDirPath := os.Getenv("SOURCEMAP_FILE")
			if sourceMapDirPath == "" {
				sourceMapDirPath = buildDirPath
			}

			possibleSourceMapPath := filepath.Join(sourceMapDirPath, "sourcemaps", "main.jsbundle.map")
			if utils.FileExists(possibleSourceMapPath) {
				sourceMapPath = possibleSourceMapPath
				log.Info("Found source map at: " + sourceMapPath)
			} else {
				log.Info("No source map found at: " + possibleSourceMapPath)
			}
		}

		// Check that we now have a source map path
		if sourceMapPath == "" {
			return errors.New("Could not find a source map, please specify the path by using --source-map or SOURCEMAP_FILE environment variable")
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

	err = server.ProcessFileRequest(endpoint+"/react-native-source-map", uploadOptions, fileFieldData, timeout, retries, sourceMapPath, dryRun)

	if err != nil {
		return err
	} else {
		log.Success("Uploaded " + filepath.Base(sourceMapPath))
	}

	return nil
}
