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

type ReactNativeAndroid struct {
	Path         utils.Paths `arg:"" name:"path" help:"The path to the root of the React Native project to upload files from" type:"path" default:"."`
	AppManifest  string      `help:"The path to a manifest file (AndroidManifest.xml) from which to obtain build information" type:"path"`
	Bundle       string      `help:"The path to the bundled JavaScript file to upload" type:"path"`
	CodeBundleId string      `help:"A unique identifier for the JavaScript bundle"`
	Dev          bool        `help:"Indicates whether this is a debug or release build"`
	ProjectRoot  string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`
	SourceMap    string      `help:"The path to the source map file to upload" type:"path"`
	Variant      string      `help:"The build type/flavor (e.g. debug, release) used to disambiguate the between built files when searching the project directory"`
	VersionCode  string      `help:"The version code of this build of the application"`
	VersionName  string      `help:"The version of the application"`
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
			} else if options.Bundle == "" || options.SourceMap == "" {
				return fmt.Errorf("unable to find bundle files or source maps in within %s", path)
			}
		}

		if options.ProjectRoot == "" {
			options.ProjectRoot = rootDirPath
		}

		if options.Bundle == "" {
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
				if options.Variant == "" {
					variantDirName, err = android.GetVariantDirectory(bundleDirPath)
					if err != nil {
						return err
					}
				} else {
					if variantFileFormat != "" {
						variantDirName = fmt.Sprintf(variantFileFormat,
							cases.Title(language.Und, cases.NoLower).String(options.Variant))

					} else {
						variantDirName = options.Variant
					}
				}
				options.Bundle = filepath.Join(bundleDirPath, variantDirName, "index.android.bundle")
			}
		}

		if !utils.FileExists(options.Bundle) {
			return fmt.Errorf("unable to find index.android.bundle at %s", options.Bundle)
		}

		if options.SourceMap == "" {
			sourceMapDirPath := filepath.Join(buildDirPath, "generated", "sourcemaps", "react")

			if options.Variant == "" {
				options.Variant, err = android.GetVariantDirectory(sourceMapDirPath)
				if err != nil {
					return err
				}
			}

			options.SourceMap = filepath.Join(sourceMapDirPath, options.Variant, "index.android.bundle.map")
		} else {
			if options.Variant == "" {
				//	Set the options.Variant based off the source map file location
				sourceMapDirPath := filepath.Join(options.SourceMap, "..", "..")

				if filepath.Base(sourceMapDirPath) == "react" {
					options.Variant, err = android.GetVariantDirectory(sourceMapDirPath)
					if err != nil {
						return err
					}
				}
			}
		}

		if !utils.FileExists(options.SourceMap) {
			return fmt.Errorf("unable to find index.android.bundle at %s", options.SourceMap)
		}

		if options.AppManifest == "" {
			appManifestPathExpected := filepath.Join(buildDirPath, "intermediates", "merged_manifests", options.Variant, "AndroidManifest.xml")
			if utils.FileExists(appManifestPathExpected) {
				options.AppManifest = appManifestPathExpected
				logger.Debug(fmt.Sprintf("Found app manifest at: %s", options.AppManifest))
			} else {
				logger.Debug(fmt.Sprintf("No app manifest found at: %s", appManifestPathExpected))
			}
		}

		if options.AppManifest != "" && (apiKey == "" || options.VersionName == "" || options.VersionCode == "") {

			manifestData, err := android.ParseAndroidManifestXML(options.AppManifest)

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

			if options.VersionName == "" {
				options.VersionName = manifestData.VersionName
				logger.Debug(fmt.Sprintf("Using %s as version name from AndroidManifest.xml", options.VersionName))
			}

			if options.VersionCode == "" {
				options.VersionCode = manifestData.VersionCode
				logger.Debug(fmt.Sprintf("Using %s as version code from AndroidManifest.xml", options.VersionCode))
			}
		}

		uploadOptions, err = utils.BuildReactNativeUploadOptions(apiKey, options.VersionName, options.VersionCode, options.CodeBundleId, options.Dev, options.ProjectRoot, overwrite, "android")

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
	}

	return nil
}
