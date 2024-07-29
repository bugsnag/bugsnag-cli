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

type ReactNativeIos struct {
	Path          utils.Paths `arg:"" name:"path" help:"The path to the root of the React Native project to upload files from" type:"path" default:"."`
	Bundle        string      `help:"The path to the bundled JavaScript file to upload" type:"path"`
	BundleVersion string      `help:"The bundle version of this build of the application (Apple platforms only)"`
	CodeBundleID  string      `help:"A unique identifier for the JavaScript bundle"`
	Dev           bool        `help:"Indicates whether this is a debug or release build"`
	Plist         string      `help:"The path to a .plist file from which to obtain build information" type:"path"`
	ProjectRoot   string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`
	Scheme        string      `help:"The name of the Xcode options.Scheme used to build the application"`
	SourceMap     string      `help:"The path to the source map file to upload" type:"path"`
	VersionName   string      `help:"The version of the application"`
	XcodeProject  string      `help:"The path to an Xcode project, workspace or containing directory from which to obtain build information" type:"path"`
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

			} else if options.Bundle == "" || options.SourceMap == "" {
				return fmt.Errorf("unable to find bundle files or source maps in within %s", path)
			}
		}

		// Set a default value for options.ProjectRoot if it's not defined
		if options.ProjectRoot == "" {
			options.ProjectRoot = rootDirPath
		}

		// Attempt to parse information from the .xcworkspace file if values aren't provided on the command line
		if options.Bundle == "" || (options.Plist == "" && (apiKey == "" || options.VersionName == "" || options.BundleVersion == "")) {

			// Validate workspacePath (if provided) or attempt to find one
			if options.XcodeProject != "" {
				if !utils.FileExists(options.XcodeProject) {
					return fmt.Errorf("unable to find the specified Xcode project file: %s", options.XcodeProject)
				}
			} else {
				if ios.IsPathAnXcodeProjectOrWorkspace(filepath.Join(rootDirPath, "ios")) {
					options.XcodeProject = filepath.Join(rootDirPath, "ios")
				}
			}

			// Validate the options.Scheme name (if provided) or attempt to find one in the workspace
			if options.XcodeProject != "" {
				// If options.Scheme is set explicitly, check if it exists
				if options.Scheme != "" {
					_, err := ios.IsSchemeInPath(options.XcodeProject, options.Scheme)
					if err != nil {
						logger.Warn(err.Error())
					}
				} else {
					// Otherwise, try to find it
					options.Scheme, err = ios.GetDefaultScheme(options.XcodeProject)
					if err != nil {
						logger.Warn(err.Error())
					}
				}

				if options.Scheme != "" {
					buildSettings, err = ios.GetXcodeBuildSettings(options.XcodeProject, options.Scheme)
					if err != nil {
						logger.Warn(err.Error())
					}
				}
			} else {
				return fmt.Errorf("Could not find an Xcode project file, please specify the path by using --xcode-proj-path")
			}

			// Check to see if we have the Info.Plist path
			if options.Plist != "" {
				if !utils.FileExists(options.Plist) {
					return fmt.Errorf("unable to find specified Info.plist file: %s", options.Plist)
				}
			} else if buildSettings != nil {
				// If not, we need to build it from build settings values
				plistPathExpected := filepath.Join(buildSettings.ConfigurationBuildDir, buildSettings.InfoPlistPath)
				if utils.FileExists(plistPathExpected) {
					options.Plist = plistPathExpected
					logger.Debug(fmt.Sprintf("Found Info.plist at expected location: %s", options.Plist))
				} else {
					plistPathExpected = filepath.Join(buildSettings.ProjectTempRoot, "ArchiveIntermediates", options.Scheme, "BuildProductsPath", filepath.Base(buildSettings.BuiltProductsDir), buildSettings.InfoPlistPath)
					if utils.FileExists(plistPathExpected) {
						options.Plist = plistPathExpected
						logger.Debug(fmt.Sprintf("Found Info.plist at: %s", options.Plist))
					} else {
						logger.Debug(fmt.Sprintf("No Info.plist found at: %s", plistPathExpected))
					}
				}
			}

		}

		// Check that the bundle file exists and error out if it doesn't
		if options.Bundle != "" {
			if !utils.FileExists(options.Bundle) {
				return fmt.Errorf("unable to find specified bundle file: %s", options.Bundle)
			}
		} else {
			// Set a options.Bundle if it's not defined and check that it exists before proceeding
			if buildSettings != nil {
				possibleBundleFilePath := filepath.Join(buildSettings.ConfigurationBuildDir, "main.jsbundle")
				if utils.FileExists(possibleBundleFilePath) {
					options.Bundle = possibleBundleFilePath
					logger.Debug(fmt.Sprintf("Found bundle file at: %s", options.Bundle))
				} else {
					possibleBundleFilePath = filepath.Join(buildSettings.ProjectTempRoot, "ArchiveIntermediates", options.Scheme, "BuildProductsPath", filepath.Base(buildSettings.BuiltProductsDir), "main.jsbundle")
					if utils.FileExists(possibleBundleFilePath) {
						options.Bundle = possibleBundleFilePath
						logger.Debug(fmt.Sprintf("Found bundle file at: %s", options.Bundle))
					} else {
						logger.Debug(fmt.Sprintf("No bundle file found at: %s", possibleBundleFilePath))
					}
				}
			}
		}

		// Check that we now have a bundle path
		if options.Bundle == "" {
			return fmt.Errorf("Could not find a bundle file, please specify the path by using --bundle-path")
		}

		// Check that the source map file exists and error out if it doesn't
		if options.SourceMap != "" {
			if !utils.FileExists(options.SourceMap) {
				return fmt.Errorf("unable to find specified source map: %s", options.SourceMap)
			}
		} else {
			// Use SOURCEMAP_FILE environment variable, if defined, or use the build directory
			sourceMapDirPath := os.Getenv("SOURCEMAP_FILE")
			if sourceMapDirPath == "" {
				sourceMapDirPath = buildDirPath
			}

			possibleSourceMapPath := filepath.Join(sourceMapDirPath, "sourcemaps", "main.jsbundle.map")
			if utils.FileExists(possibleSourceMapPath) {
				options.SourceMap = possibleSourceMapPath
				logger.Debug(fmt.Sprintf("Found source map at: %s", options.SourceMap))
			} else {
				logger.Debug(fmt.Sprintf("No source map found at: %s", possibleSourceMapPath))
			}
		}

		// Check that we now have a source map path
		if options.SourceMap == "" {
			return fmt.Errorf("Could not find a source map, please specify the path by using --source-map or SOURCEMAP_FILE environment variable")
		}

		if options.Plist != "" && (apiKey == "" || options.VersionName == "" || options.BundleVersion == "") {
			// Read data from the plist
			plistData, err := ios.GetPlistData(options.Plist)
			if err != nil {
				return err
			}

			// Check if the variables are empty, set if they are abd log that we are using setting from the plist file
			if options.BundleVersion == "" {
				options.BundleVersion = plistData.BundleVersion
				logger.Debug(fmt.Sprintf("Using bundle version from Info.plist: %s", options.BundleVersion))
			}

			if options.VersionName == "" {
				options.VersionName = plistData.VersionName
				logger.Debug(fmt.Sprintf("Using version name from Info.plist: %s", options.VersionName))

			}

			if apiKey == "" {
				apiKey = plistData.BugsnagProjectDetails.ApiKey
				logger.Debug(fmt.Sprintf("Using API key from Info.plist: %s", apiKey))
			}

		}

	}

	uploadOptions, err := utils.BuildReactNativeUploadOptions(apiKey, options.VersionName, options.BundleVersion, options.CodeBundleID, options.Dev, options.ProjectRoot, overwrite, "ios")

	if err != nil {
		return err
	}

	fileFieldData := make(map[string]server.FileField)
	fileFieldData["sourceMap"] = server.LocalFile(options.SourceMap)
	fileFieldData["bundle"] = server.LocalFile(options.Bundle)

	err = server.ProcessFileRequest(endpoint+"/react-native-source-map", uploadOptions, fileFieldData, timeout, retries, options.SourceMap, dryRun, logger)

	if err != nil {

		return err
	}

	return nil
}
