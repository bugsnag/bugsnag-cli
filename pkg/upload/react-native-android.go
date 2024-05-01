package upload

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type ReactNativeAndroid struct {
	AppManifest  string      `help:"(required) Path to directory or file to upload" type:"path"`
	Bundle       string      `help:"Path to the bundle file" type:"path"`
	CodeBundleId string      `help:"A unique identifier to identify a code bundle release when using tools like CodePush"`
	Dev          bool        `help:"Indicates whether the application is a debug or release build"`
	Path         utils.Paths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
	ProjectRoot  string      `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	SourceMap    string      `help:"Path to the source map file" type:"path"`
	Variant      string      `help:"Build type, like 'debug' or 'release'"`
	VersionName  string      `help:"The version name of the application."`
	VersionCode  string      `help:"The version code for the application (Android only)."`
}

func ProcessReactNativeAndroid(
	apiKey string,
	appManifestPath string,
	bundlePath string,
	codeBundleId string,
	dev bool,
	paths []string,
	projectRoot string,
	variant string,
	versionName string,
	versionCode string,
	sourceMapPath string,
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

	for _, path := range paths {

		buildDirPath := filepath.Join(path, "android", "app", "build")
		rootDirPath = path
		if !utils.FileExists(buildDirPath) {
			buildDirPath = filepath.Join(path, "app", "build")
			if utils.FileExists(buildDirPath) {
				rootDirPath = filepath.Join(path, "..")
			} else if bundlePath == "" || sourceMapPath == "" {
				return fmt.Errorf("unable to find bundle files or source maps in within %s", path)
			}
		}

		if projectRoot == "" {
			projectRoot = rootDirPath
		}

		if bundlePath == "" {
			if utils.IsDir(filepath.Join(buildDirPath, "generated", "assets", "react")) {
				// RN version < 0.70 - generated/assets/react/<variant>/index.android.bundle
				bundleDirPath = filepath.Join(buildDirPath, "generated", "assets", "react")
			} else if utils.IsDir(filepath.Join(buildDirPath, "ASSETS")) {
				// RN versions < 0.72 - ASSETS/createBundle<variant>JsAndAssets/index.android.bundle
				bundleDirPath = filepath.Join(buildDirPath, "ASSETS")
				variantFileFormat = "createBundle%sJsAndAssets"
			} else if utils.IsDir(filepath.Join(buildDirPath, "generated", "assets")) {
				// RN versions >= 0.72 - generated/assets/<variant>/index.android.bundle
				bundleDirPath = filepath.Join(buildDirPath, "generated", "assets")
				variantFileFormat = "createBundle%sJsAndAssets"
			} else {
				return fmt.Errorf("unable to find index.android.bundle in your project, please specify the path using --bundle-path")
			}

			if bundleDirPath != "" {
				if variant == "" {
					variantDirName, err = android.GetVariantDirectory(bundleDirPath)
					if err != nil {
						return err
					}
				} else {
					if variantFileFormat != "" {
						variantDirName = fmt.Sprintf(variantFileFormat,
							cases.Title(language.Und, cases.NoLower).String(variant))

					} else {
						variantDirName = variant
					}
				}
				bundlePath = filepath.Join(bundleDirPath, variantDirName, "index.android.bundle")
			}
		}

		if !utils.FileExists(bundlePath) {
			return fmt.Errorf("unable to find index.android.bundle at %s", bundlePath)
		}

		if sourceMapPath == "" {
			sourceMapDirPath := filepath.Join(buildDirPath, "generated", "sourcemaps", "react")

			if variant == "" {
				variant, err = android.GetVariantDirectory(sourceMapDirPath)
				if err != nil {
					return err
				}
			}

			sourceMapPath = filepath.Join(sourceMapDirPath, variant, "index.android.bundle.map")
		} else {
			if variant == "" {
				//	Set the variant based off the source map file location
				sourceMapDirPath := filepath.Join(sourceMapPath, "..", "..")

				if filepath.Base(sourceMapDirPath) == "react" {
					variant, err = android.GetVariantDirectory(sourceMapDirPath)
					if err != nil {
						return err
					}
				}
			}
		}

		if !utils.FileExists(sourceMapPath) {
			return fmt.Errorf("unable to find index.android.bundle at %s", sourceMapPath)
		}

		if appManifestPath == "" {
			appManifestPathExpected := filepath.Join(buildDirPath, "intermediates", "merged_manifests", variant, "AndroidManifest.xml")
			if utils.FileExists(appManifestPathExpected) {
				appManifestPath = appManifestPathExpected
				logger.Info(fmt.Sprintf("Found app manifest at: %s", appManifestPath))
			} else {
				logger.Info(fmt.Sprintf("No app manifest found at: %s", appManifestPathExpected))
			}
		}

		if appManifestPath != "" && (apiKey == "" || versionName == "" || versionCode == "") {

			manifestData, err := android.ParseAndroidManifestXML(appManifestPath)

			if err != nil {
				return err
			}

			if apiKey == "" {
				for key, value := range manifestData.Application.MetaData.Name {
					if value == "com.bugsnag.android.API_KEY" {
						apiKey = manifestData.Application.MetaData.Value[key]
					}
				}
				logger.Info(fmt.Sprintf("Using %s as API key from AndroidManifest.xml", apiKey))
			}

			if versionName == "" {
				versionName = manifestData.VersionName
				logger.Info(fmt.Sprintf("Using %s as version name from AndroidManifest.xml", versionName))
			}

			if versionCode == "" {
				versionCode = manifestData.VersionCode
				logger.Info(fmt.Sprintf("Using %s as version code from AndroidManifest.xml", versionCode))
			}
		}

		uploadOptions, err = utils.BuildReactNativeUploadOptions(apiKey, versionName, versionCode, codeBundleId, dev, projectRoot, overwrite, "android")

		if err != nil {
			return err
		}

		fileFieldData := make(map[string]string)
		fileFieldData["sourceMap"] = sourceMapPath
		fileFieldData["bundle"] = bundlePath

		err = server.ProcessFileRequest(endpoint+"/react-native-source-map", uploadOptions, fileFieldData, timeout, retries, sourceMapPath, dryRun, logger)

		if err != nil {

			return err
		}
	}

	return nil
}
