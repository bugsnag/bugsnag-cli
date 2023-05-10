package upload

import (
	"fmt"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type ReactNativeAndroid struct {
	AppManifest  string            `help:"(required) Path to directory or file to upload" type:"path"`
	Bundle       string            `help:"Path to the bundle file" type:"path"`
	CodeBundleId string            `help:"A unique identifier to identify a code bundle release when using tools like CodePush"`
	Dev          bool              `help:"Indicates whether the application is a debug or release build"`
	Path         utils.UploadPaths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
	ProjectRoot  string            `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	SourceMap    string            `help:"Path to the source map file" type:"path"`
	Variant      string            `help:"Build type, like 'debug' or 'release'"`
	Version      string            `help:"The version of the application."`
	VersionCode  string            `help:"The version code for the application (Android only)."`
}

func ProcessReactNativeAndroid(apiKey string, appManifestPath string, bundlePath string, codeBundleId string, dev bool, paths []string, projectRoot string, variant string, version string, versionCode string, sourceMapPath string, endpoint string, timeout int, retries int, overwrite bool, dryRun bool) error {

	var err error
	var uploadOptions map[string]string
	var rootDirPath string

	if dryRun {
		log.Info("Performing dry run - no files will be uploaded")
	}

	for _, path := range paths {

		buildDirPath := filepath.Join(path, "android", "app", "build")
		rootDirPath = path
		if !utils.FileExists(buildDirPath) {
			buildDirPath = filepath.Join(path, "app", "build")
			if utils.FileExists(buildDirPath) {
				rootDirPath = filepath.Join(path, "..")
			} else if bundlePath == "" || sourceMapPath == "" {
				return fmt.Errorf("unable to find bundle files or source maps in within " + path)
			}
		}

		if projectRoot == "" {
			projectRoot = rootDirPath
		}

		if bundlePath == "" {
			bundlePath = filepath.Join(buildDirPath, "ASSETS", "createBundleReleaseJsAndAssets", "index.android.bundle")
		}

		if !utils.FileExists(bundlePath) {
			return fmt.Errorf("unable to find index.android.bundle at " + bundlePath)
		}

		if sourceMapPath == "" {
			sourceMapDirPath := filepath.Join(buildDirPath, "generated", "sourcemaps", "react")

			if variant == "" {
				variant, err = android.GetVariant(sourceMapDirPath)
			}

			sourceMapPath = filepath.Join(sourceMapDirPath, variant, "index.android.bundle.map")
		} else {
			if variant == "" {
				//	Set the variant based off the source map file location
				sourceMapDirPath := filepath.Join(sourceMapPath, "..", "..")

				if filepath.Base(sourceMapDirPath) == "react" {
					variant, err = android.GetVariant(sourceMapDirPath)
				}
			}
		}

		if !utils.FileExists(sourceMapPath) {
			return fmt.Errorf("unable to find index.android.bundle at " + sourceMapPath)
		}

		if appManifestPath == "" {
			appManifestPathExpected := filepath.Join(buildDirPath, "intermediates", "merged_manifests", variant, "AndroidManifest.xml")
			if utils.FileExists(appManifestPathExpected) {
				appManifestPath = appManifestPathExpected
				log.Info("Found app manifest at: " + appManifestPath)
			}
		}

		if appManifestPath != "" && (apiKey == "" || version == "" || versionCode == "") {

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

				log.Info("Using " + apiKey + " as API key from AndroidManifest.xml")
			}

			if version == "" {
				version = manifestData.VersionName
				log.Info("Using " + version + " as version code from AndroidManifest.xml")
			}

			if versionCode == "" {
				versionCode = manifestData.VersionCode
				log.Info("Using " + versionCode + " as version name from AndroidManifest.xml")
			}
		}

		log.Info("Uploading debug information for React Native Android")

		uploadOptions, err = utils.BuildReactNativeAndroidUploadOptions(apiKey, version, versionCode, codeBundleId, dev, projectRoot, overwrite)

		if err != nil {
			return err
		}

		fileFieldData := make(map[string]string)
		fileFieldData["sourceMap"] = sourceMapPath
		fileFieldData["bundle"] = bundlePath

		if dryRun {
			err = nil
		} else {
			err = server.ProcessRequest(endpoint+"/react-native-source-map", uploadOptions, fileFieldData, timeout)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
