package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
)

// ProcessReactNativeIos processes React Native iOS bundle and source map uploads.
//
// It locates the bundle and source map files, resolves Xcode projects, schemes, and plist data,
// builds upload options, and sends the files to the Bugsnag server.
//
// Parameters:
//   - options: CLI options containing upload settings and flags.
//   - logger: Logger instance for debug and error output.
//
// Returns:
//   - error: non-nil if an error occurs during processing or uploading.
func ProcessReactNativeIos(options options.CLI, logger log.Logger) error {
	iosOptions := options.Upload.ReactNativeIos
	var (
		rootDirPath      string
		plistData        *ios.PlistData
		err              error
		xcodeArchivePath string
		buildDirPath     string
	)

	for _, path := range iosOptions.Path {
		if filepath.Ext(path) == ".xcarchive" {
			xcodeArchivePath = path
		} else if utils.IsDir(path) {
			buildDirPath = filepath.Join(path, "ios", "build")
			rootDirPath = path
			if !utils.FileExists(buildDirPath) {
				buildDirPath = filepath.Join(path, "build")
				if utils.FileExists(buildDirPath) {
					rootDirPath = filepath.Join(path, "..")
				} else if iosOptions.ReactNative.Bundle == "" || iosOptions.ReactNative.SourceMap == "" {
					return fmt.Errorf("unable to find bundle files or source maps in within %s", path)
				}
			}

			// Set a default value for options.ProjectRoot if it's not defined
			if iosOptions.ProjectRoot == "" {
				iosOptions.ProjectRoot = rootDirPath
			}

			// Validate workspacePath (if provided) or attempt to find one
			if iosOptions.Ios.XcodeProject != "" {
				if !utils.FileExists(iosOptions.Ios.XcodeProject) {
					return fmt.Errorf("unable to find the specified Xcode project file: %s", iosOptions.Ios.XcodeProject)
				}
			} else {
				if ios.IsPathAnXcodeProjectOrWorkspace(filepath.Join(rootDirPath, "ios")) {
					iosOptions.Ios.XcodeProject = filepath.Join(rootDirPath, "ios")
				}
			}

			logger.Debug(fmt.Sprintf("Found Xcode project at: %s", iosOptions.Ios.XcodeProject))

			// Validate the options.Ios.Scheme name (if provided) or attempt to find one in the workspace
			if iosOptions.Ios.XcodeProject != "" {
				// If options.Ios.Scheme is set explicitly, check if it exists
				if iosOptions.Ios.Scheme != "" {
					_, err := ios.IsSchemeInPath(iosOptions.Ios.XcodeProject, iosOptions.Ios.Scheme)
					if err != nil {
						logger.Warn(err.Error())
					}
				} else {
					// Otherwise, try to find it
					iosOptions.Ios.Scheme, err = ios.GetDefaultScheme(iosOptions.Ios.XcodeProject)
					if err != nil {
						logger.Warn(err.Error())
					}
				}
			} else {
				return fmt.Errorf("could not find an Xcode project file, please specify the path by using --xcode-project")
			}

			logger.Debug(fmt.Sprintf("Found Xcode scheme: %s", iosOptions.Ios.Scheme))

			if iosOptions.Ios.XcarchivePath != "" {
				xcodeArchivePath = string(iosOptions.Ios.XcarchivePath)
			} else {
				xcodeArchivePath, err = ios.GetLatestXcodeArchiveForScheme(iosOptions.Ios.Scheme)

				if err != nil {
					return fmt.Errorf("error locating latest Xcode archive from Xcode project (scheme: %s), please specify the xcarchive path directly using --xcarchive-path: %w", iosOptions.Ios.Scheme, err)
				}
			}

		} else {
			return fmt.Errorf("path should be an Xcode archive or the directory of your React Native project: %s", path)
		}

		logger.Debug(fmt.Sprintf("Found Xcode archive Path: %s", xcodeArchivePath))

		// Attempt to parse information from the .xcworkspace file if values aren't provided on the command line
		if iosOptions.ReactNative.Bundle == "" || (iosOptions.Ios.Plist == "" && (options.ApiKey == "" || iosOptions.ReactNative.VersionName == "" || iosOptions.Ios.BundleVersion == "")) {
			// Check to see if we have the Info.Plist path
			if iosOptions.Ios.Plist != "" {
				if !utils.FileExists(iosOptions.Ios.Plist) {
					return fmt.Errorf("unable to find specified Info.plist file: %s", iosOptions.Ios.Plist)
				}
			} else if xcodeArchivePath != "" {
				// If not, we need to build it from build settings values
				plistPathExpected := filepath.Join(xcodeArchivePath, "Products", "Applications", iosOptions.Ios.Scheme+".app", "Info.plist")
				if utils.FileExists(plistPathExpected) {
					iosOptions.Ios.Plist = plistPathExpected
					logger.Debug(fmt.Sprintf("Found Info.plist at expected location: %s", iosOptions.Ios.Plist))
				} else {
					logger.Debug(fmt.Sprintf("No Info.plist found at: %s", plistPathExpected))
				}
			}
		}
	}

	// Check that the bundle file exists and error out if it doesn't
	if iosOptions.ReactNative.Bundle != "" {
		if !utils.FileExists(iosOptions.ReactNative.Bundle) {
			return fmt.Errorf("unable to find specified bundle file: %s", iosOptions.ReactNative.Bundle)
		}
	} else {
		// Set a options.Bundle if it's not defined and check that it exists before proceeding
		if xcodeArchivePath != "" {
			possibleBundleFilePath := filepath.Join(xcodeArchivePath, "Products", "Applications", iosOptions.Ios.Scheme+".app", "main.jsbundle")
			if utils.FileExists(possibleBundleFilePath) {
				iosOptions.ReactNative.Bundle = possibleBundleFilePath
				logger.Debug(fmt.Sprintf("Found bundle file at: %s", iosOptions.ReactNative.Bundle))
			} else {
				logger.Debug(fmt.Sprintf("No bundle file found at: %s", possibleBundleFilePath))
			}
		}
	}

	// Check that we now have a bundle path
	if iosOptions.ReactNative.Bundle == "" {
		return fmt.Errorf("Could not find a bundle file, please specify the path by using --bundle")
	}

	// Check that the source map file exists and error out if it doesn't
	if iosOptions.ReactNative.SourceMap != "" {
		if !utils.FileExists(iosOptions.ReactNative.SourceMap) {
			return fmt.Errorf("Unable to find specified source map: %s", iosOptions.ReactNative.SourceMap)
		}
	} else {
		// Use SOURCEMAP_FILE environment variable, if defined, or use the build directory
		sourceMapDirPath := os.Getenv("SOURCEMAP_FILE")
		if sourceMapDirPath == "" {
			sourceMapDirPath = buildDirPath
		}

		possibleSourceMapPath := filepath.Join(sourceMapDirPath, "sourcemaps", "main.jsbundle.map")
		if utils.FileExists(possibleSourceMapPath) {
			iosOptions.ReactNative.SourceMap = possibleSourceMapPath
			logger.Debug(fmt.Sprintf("Found source map at: %s", iosOptions.ReactNative.SourceMap))
		} else {
			logger.Debug(fmt.Sprintf("No source map found at: %s", possibleSourceMapPath))
		}
	}

	// Check that we now have a source map path
	if iosOptions.ReactNative.SourceMap == "" {
		return fmt.Errorf("Could not find a source map, please specify the path by using --source-map or SOURCEMAP_FILE environment variable")
	}

	if iosOptions.Ios.Plist != "" && (options.ApiKey == "" || iosOptions.ReactNative.VersionName == "" || iosOptions.Ios.BundleVersion == "") {
		// Read data from the plist
		plistData, err = ios.GetPlistData(iosOptions.Ios.Plist)
		if err != nil {
			return err
		}

		// If we've not passed --code-bundle-id, proceed to populate versionName and versionCode from the plist
		if iosOptions.ReactNative.CodeBundleId == "" {
			// Check if the variables are empty, set if they are abd log that we are using setting from the plist file
			if iosOptions.Ios.BundleVersion == "" {
				iosOptions.Ios.BundleVersion = plistData.BundleVersion
				logger.Debug(fmt.Sprintf("Using bundle version from Info.plist: %s", iosOptions.Ios.BundleVersion))
			}

			if iosOptions.ReactNative.VersionName == "" {
				iosOptions.ReactNative.VersionName = plistData.VersionName
				logger.Debug(fmt.Sprintf("Using version name from Info.plist: %s", iosOptions.ReactNative.VersionName))

			}
		}

		if options.ApiKey == "" {
			options.ApiKey = plistData.BugsnagProjectDetails.ApiKey
			logger.Debug(fmt.Sprintf("Using API key from Info.plist: %s", options.ApiKey))
		}

	}

	uploadOptions, err := utils.BuildReactNativeUploadOptions(iosOptions.ReactNative.VersionName, iosOptions.Ios.BundleVersion, iosOptions.ReactNative.CodeBundleId, iosOptions.ReactNative.Dev, iosOptions.ProjectRoot, options.Upload.Overwrite, "ios")

	if err != nil {
		return err
	}

	fileFieldData := make(map[string]server.FileField)
	fileFieldData["sourceMap"] = server.LocalFile(iosOptions.ReactNative.SourceMap)
	fileFieldData["bundle"] = server.LocalFile(iosOptions.ReactNative.Bundle)

	err = server.ProcessFileRequest(options.ApiKey, "/react-native-source-map", uploadOptions, fileFieldData, iosOptions.ReactNative.SourceMap, options, logger)

	if err != nil {

		return err
	}

	return nil
}
