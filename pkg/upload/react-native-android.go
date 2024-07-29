package upload

import (
	"fmt"
	"path/filepath"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type ReactNativeAndroidSpecific struct {
	AppManifest string `help:"The path to a manifest file (AndroidManifest.xml) from which to obtain build information" type:"path"`
	Variant     string `help:"The build type/flavor (e.g. debug, release) used to disambiguate the between built files when searching the project directory"`
	VersionCode string `help:"The version code of this build of the application"`
}

type ReactNativeAndroid struct {
	Path        utils.Paths `arg:"" name:"path" help:"The path to the root of the React Native project to upload files from" type:"path" default:"."`
	ProjectRoot string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`

	ReactNative ReactNativeShared          `embed:""`
	Android     ReactNativeAndroidSpecific `embed:""`
}

func ProcessReactNativeAndroid(
	apiKey string,
	options ReactNativeAndroid,
	endpoint string,
	timeout int,
	retries int,
	overwrite bool,
	dryRun bool,
	logger log.Logger,
) error {

	var err error
	var uploadOptions map[string]string
	var rootDirPath string
	var variantDirName string
	var bundleDirPath string
	var variantFileFormat string

	for _, path := range options.Path {

		buildDirPath := filepath.Join(path, "android", "app", "build")
		rootDirPath = path
		if !utils.FileExists(buildDirPath) {
			buildDirPath = filepath.Join(path, "app", "build")
			if utils.FileExists(buildDirPath) {
				rootDirPath = filepath.Join(path, "..")
			} else if options.ReactNative.Bundle == "" || options.ReactNative.SourceMap == "" {
				return fmt.Errorf("unable to find bundle files or source maps in within %s", path)
			}
		}

		if options.ProjectRoot == "" {
			options.ProjectRoot = rootDirPath
		}

		if options.ReactNative.Bundle == "" {
			if utils.IsDir(filepath.Join(buildDirPath, "generated", "assets", "react")) {
				// RN version < 0.70 - generated/assets/react/<options.Variant>/index.android.bundle
				bundleDirPath = filepath.Join(buildDirPath, "generated", "assets", "react")
			} else if utils.IsDir(filepath.Join(buildDirPath, "ASSETS")) {
				// RN versions < 0.72 - ASSETS/createBundle<options.Variant>JsAndAssets/index.android.bundle
				bundleDirPath = filepath.Join(buildDirPath, "ASSETS")
				variantFileFormat = "createBundle%sJsAndAssets"
			} else if utils.IsDir(filepath.Join(buildDirPath, "generated", "assets")) {
				// RN versions >= 0.72 - generated/assets/<options.Variant>/index.android.bundle
				bundleDirPath = filepath.Join(buildDirPath, "generated", "assets")
				variantFileFormat = "createBundle%sJsAndAssets"
			} else {
				return fmt.Errorf("unable to find index.android.bundle in your project, please specify the path using --bundle-path")
			}

			if bundleDirPath != "" {
				if options.Android.Variant == "" {
					variantDirName, err = android.GetVariantDirectory(bundleDirPath)
					if err != nil {
						return err
					}
				} else {
					if variantFileFormat != "" {
						variantDirName = fmt.Sprintf(variantFileFormat,
							cases.Title(language.Und, cases.NoLower).String(options.Android.Variant))

					} else {
						variantDirName = options.Android.Variant
					}
				}
				options.ReactNative.Bundle = filepath.Join(bundleDirPath, variantDirName, "index.android.bundle")
			}
		}

		if !utils.FileExists(options.ReactNative.Bundle) {
			return fmt.Errorf("unable to find index.android.bundle at %s", options.ReactNative.Bundle)
		}

		if options.ReactNative.SourceMap == "" {
			sourceMapDirPath := filepath.Join(buildDirPath, "generated", "sourcemaps", "react")

			if options.Android.Variant == "" {
				options.Android.Variant, err = android.GetVariantDirectory(sourceMapDirPath)
				if err != nil {
					return err
				}
			}

			options.ReactNative.SourceMap = filepath.Join(sourceMapDirPath, options.Android.Variant, "index.android.bundle.map")
		} else {
			if options.Android.Variant == "" {
				//	Set the options.Variant based off the source map file location
				sourceMapDirPath := filepath.Join(options.ReactNative.SourceMap, "..", "..")

				if filepath.Base(sourceMapDirPath) == "react" {
					options.Android.Variant, err = android.GetVariantDirectory(sourceMapDirPath)
					if err != nil {
						return err
					}
				}
			}
		}

		if !utils.FileExists(options.ReactNative.SourceMap) {
			return fmt.Errorf("unable to find index.android.bundle at %s", options.ReactNative.SourceMap)
		}

		if options.Android.AppManifest == "" {
			appManifestPathExpected := filepath.Join(buildDirPath, "intermediates", "merged_manifests", options.Android.Variant, "AndroidManifest.xml")
			if utils.FileExists(appManifestPathExpected) {
				options.Android.AppManifest = appManifestPathExpected
				logger.Debug(fmt.Sprintf("Found app manifest at: %s", options.Android.AppManifest))
			} else {
				logger.Debug(fmt.Sprintf("No app manifest found at: %s", appManifestPathExpected))
			}
		}

		if options.Android.AppManifest != "" && (apiKey == "" || options.ReactNative.VersionName == "" || options.Android.VersionCode == "") {

			manifestData, err := android.ParseAndroidManifestXML(options.Android.AppManifest)

			if err != nil {
				return err
			}

			if apiKey == "" {
				for key, value := range manifestData.Application.MetaData.Name {
					if value == "com.bugsnag.android.API_KEY" {
						apiKey = manifestData.Application.MetaData.Value[key]
					}
				}
				logger.Debug(fmt.Sprintf("Using %s as API key from AndroidManifest.xml", apiKey))
			}

			if options.ReactNative.VersionName == "" {
				options.ReactNative.VersionName = manifestData.VersionName
				logger.Debug(fmt.Sprintf("Using %s as version name from AndroidManifest.xml", options.ReactNative.VersionName))
			}

			if options.Android.VersionCode == "" {
				options.Android.VersionCode = manifestData.VersionCode
				logger.Debug(fmt.Sprintf("Using %s as version code from AndroidManifest.xml", options.Android.VersionCode))
			}
		}

		uploadOptions, err = utils.BuildReactNativeUploadOptions(apiKey, options.ReactNative.VersionName, options.Android.VersionCode, options.ReactNative.CodeBundleId, options.ReactNative.Dev, options.ProjectRoot, overwrite, "android")

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
	}

	return nil
}
