package upload

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type ReactNativeShared struct {
	Bundle       string `help:"The path to the bundled JavaScript file to upload" type:"path"`
	CodeBundleId string `help:"A unique identifier for the JavaScript bundle"`
	Dev          bool   `help:"Indicates whether this is a debug or release build"`
	SourceMap    string `help:"The path to the source map file to upload" type:"path"`
	VersionName  string `help:"The version of the application"`
}

type ReactNativeIosSpecific struct {
	BundleVersion string `help:"The bundle version of this build of the application (Apple platforms only)"`
	Plist         string `help:"The path to a .plist file from which to obtain build information" type:"path"`
	Scheme        string `help:"The name of the Xcode options.Ios.Scheme used to build the application"`
	XcodeProject  string `help:"The path to an Xcode project, workspace or containing directory from which to obtain build information" type:"path"`
}

type ReactNativeIos struct {
	Path        utils.Paths `arg:"" name:"path" help:"The path to the root of the React Native project to upload files from" type:"path" default:"."`
	ProjectRoot string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`

	ReactNative ReactNativeShared      `embed:""`
	Ios         ReactNativeIosSpecific `embed:""`
}

func ProcessReactNativeIos(
	apiKey string,
	options ReactNativeIos,
	endpoint string,
	timeout int,
	retries int,
	overwrite bool,
	dryRun bool,
	logger log.Logger,
) error {

	var rootDirPath string
	var buildSettings *ios.XcodeBuildSettings
	var err error

	for _, path := range options.Path {
		// Check/Set the build folder
		buildDirPath := filepath.Join(path, "ios", "build")
		rootDirPath = path
		if !utils.FileExists(buildDirPath) {
			buildDirPath = filepath.Join(path, "build")

			if utils.FileExists(buildDirPath) {
				rootDirPath = filepath.Join(path, "..")

			} else if options.ReactNative.Bundle == "" || options.ReactNative.SourceMap == "" {
				return fmt.Errorf("unable to find bundle files or source maps in within %s", path)
			}
		}

		// Set a default value for options.ProjectRoot if it's not defined
		if options.ProjectRoot == "" {
			options.ProjectRoot = rootDirPath
		}

		// Attempt to parse information from the .xcworkspace file if values aren't provided on the command line
		if options.ReactNative.Bundle == "" || (options.Ios.Plist == "" && (apiKey == "" || options.ReactNative.VersionName == "" || options.Ios.BundleVersion == "")) {

			// Validate workspacePath (if provided) or attempt to find one
			if options.Ios.XcodeProject != "" {
				if !utils.FileExists(options.Ios.XcodeProject) {
					return fmt.Errorf("unable to find the specified Xcode project file: %s", options.Ios.XcodeProject)
				}
			} else {
				if ios.IsPathAnXcodeProjectOrWorkspace(filepath.Join(rootDirPath, "ios")) {
					options.Ios.XcodeProject = filepath.Join(rootDirPath, "ios")
				}
			}

			// Validate the options.Ios.Scheme name (if provided) or attempt to find one in the workspace
			if options.Ios.XcodeProject != "" {
				// If options.Ios.Scheme is set explicitly, check if it exists
				if options.Ios.Scheme != "" {
					_, err := ios.IsSchemeInPath(options.Ios.XcodeProject, options.Ios.Scheme)
					if err != nil {
						logger.Warn(err.Error())
					}
				} else {
					// Otherwise, try to find it
					options.Ios.Scheme, err = ios.GetDefaultScheme(options.Ios.XcodeProject)
					if err != nil {
						logger.Warn(err.Error())
					}
				}

				if options.Ios.Scheme != "" {
					buildSettings, err = ios.GetXcodeBuildSettings(options.Ios.XcodeProject, options.Ios.Scheme)
					if err != nil {
						logger.Warn(err.Error())
					}
				}
			} else {
				return fmt.Errorf("Could not find an Xcode project file, please specify the path by using --xcode-proj-path")
			}

			// Check to see if we have the Info.Plist path
			if options.Ios.Plist != "" {
				if !utils.FileExists(options.Ios.Plist) {
					return fmt.Errorf("unable to find specified Info.plist file: %s", options.Ios.Plist)
				}
			} else if buildSettings != nil {
				// If not, we need to build it from build settings values
				plistPathExpected := filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.InfoPlistPath)
				if utils.FileExists(plistPathExpected) {
					options.Ios.Plist = plistPathExpected
					logger.Debug(fmt.Sprintf("Found Info.plist at expected location: %s", options.Ios.Plist))
				} else {
					plistPathExpected = filepath.Join(buildSettings.ProjectTempRoot, "ArchiveIntermediates", options.Ios.Scheme, "BuildProductsPath", filepath.Base(buildSettings.BuiltProductsDir), buildSettings.InfoPlistPath)
					if utils.FileExists(plistPathExpected) {
						options.Ios.Plist = plistPathExpected
						logger.Debug(fmt.Sprintf("Found Info.plist at: %s", options.Ios.Plist))
					} else {
						logger.Debug(fmt.Sprintf("No Info.plist found at: %s", plistPathExpected))
					}
				}
			}

		}

		// Check that the bundle file exists and error out if it doesn't
		if options.ReactNative.Bundle != "" {
			if !utils.FileExists(options.ReactNative.Bundle) {
				return fmt.Errorf("unable to find specified bundle file: %s", options.ReactNative.Bundle)
			}
		} else {
			// Set a options.Bundle if it's not defined and check that it exists before proceeding
			if buildSettings != nil {
				possibleBundleFilePath := filepath.Join(buildSettings.ConfigurationBuildDir, "main.jsbundle")
				if utils.FileExists(possibleBundleFilePath) {
					options.ReactNative.Bundle = possibleBundleFilePath
					logger.Debug(fmt.Sprintf("Found bundle file at: %s", options.ReactNative.Bundle))
				} else {
					possibleBundleFilePath = filepath.Join(buildSettings.ProjectTempRoot, "ArchiveIntermediates", options.Ios.Scheme, "BuildProductsPath", filepath.Base(buildSettings.BuiltProductsDir), "main.jsbundle")
					if utils.FileExists(possibleBundleFilePath) {
						options.ReactNative.Bundle = possibleBundleFilePath
						logger.Debug(fmt.Sprintf("Found bundle file at: %s", options.ReactNative.Bundle))
					} else {
						logger.Debug(fmt.Sprintf("No bundle file found at: %s", possibleBundleFilePath))
					}
				}
			}
		}

		// Check that we now have a bundle path
		if options.ReactNative.Bundle == "" {
			return fmt.Errorf("Could not find a bundle file, please specify the path by using --bundle-path")
		}

		// Check that the source map file exists and error out if it doesn't
		if options.ReactNative.SourceMap != "" {
			if !utils.FileExists(options.ReactNative.SourceMap) {
				return fmt.Errorf("unable to find specified source map: %s", options.ReactNative.SourceMap)
			}
		} else {
			// Use SOURCEMAP_FILE environment variable, if defined, or use the build directory
			sourceMapDirPath := os.Getenv("SOURCEMAP_FILE")
			if sourceMapDirPath == "" {
				sourceMapDirPath = buildDirPath
			}

			possibleSourceMapPath := filepath.Join(sourceMapDirPath, "sourcemaps", "main.jsbundle.map")
			if utils.FileExists(possibleSourceMapPath) {
				options.ReactNative.SourceMap = possibleSourceMapPath
				logger.Debug(fmt.Sprintf("Found source map at: %s", options.ReactNative.SourceMap))
			} else {
				logger.Debug(fmt.Sprintf("No source map found at: %s", possibleSourceMapPath))
			}
		}

		// Check that we now have a source map path
		if options.ReactNative.SourceMap == "" {
			return fmt.Errorf("Could not find a source map, please specify the path by using --source-map or SOURCEMAP_FILE environment variable")
		}

		if options.Ios.Plist != "" && (apiKey == "" || options.ReactNative.VersionName == "" || options.Ios.BundleVersion == "") {
			// Read data from the plist
			plistData, err := ios.GetPlistData(options.Ios.Plist)
			if err != nil {
				return err
			}

			// Check if the variables are empty, set if they are abd log that we are using setting from the plist file
			if options.Ios.BundleVersion == "" {
				options.Ios.BundleVersion = plistData.BundleVersion
				logger.Debug(fmt.Sprintf("Using bundle version from Info.plist: %s", options.Ios.BundleVersion))
			}

			if options.ReactNative.VersionName == "" {
				options.ReactNative.VersionName = plistData.VersionName
				logger.Debug(fmt.Sprintf("Using version name from Info.plist: %s", options.ReactNative.VersionName))

			}

			if apiKey == "" {
				apiKey = plistData.BugsnagProjectDetails.ApiKey
				logger.Debug(fmt.Sprintf("Using API key from Info.plist: %s", apiKey))
			}

		}

	}

	uploadOptions, err := utils.BuildReactNativeUploadOptions(apiKey, options.ReactNative.VersionName, options.Ios.BundleVersion, options.ReactNative.CodeBundleId, options.ReactNative.Dev, options.ProjectRoot, overwrite, "ios")

	if err != nil {
		return err
	}

	fileFieldData := make(map[string]server.FileField)
	fileFieldData["sourceMap"] = server.LocalFile(options.ReactNative.SourceMap)
	fileFieldData["bundle"] = server.LocalFile(options.ReactNative.Bundle)

	err = server.ProcessFileRequest(endpoint+"/react-native-source-map", uploadOptions, fileFieldData, timeout, retries, options.ReactNative.SourceMap, dryRun, logger)

	if err != nil {

		return err
	}

	return nil
}
