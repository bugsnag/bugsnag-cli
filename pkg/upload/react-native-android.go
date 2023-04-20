package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"path/filepath"
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

	if dryRun {
		log.Info("Performing dry run - no files will be uploaded")
	}

	for _, path := range paths {

		if projectRoot == "" {
			projectRoot = path
		}

		if appManifestPath == "" {
			log.Info("Locating Android manifest")

			if utils.FileExists(filepath.Join(path, "android", "app", "build", "intermediates", "merged_manifests")) {
				appManifestPath = filepath.Join(path, "android", "app", "build", "intermediates", "merged_manifests")
			} else if utils.FileExists(filepath.Join(path, "app", "build", "intermediates", "merged_manifests")) {
				appManifestPath = filepath.Join(path, "app", "build", "intermediates", "merged_manifests")
			} else {
				return fmt.Errorf("unable to find AndroidManifest.xml. Please specify using `--app-manifest-path` ")
			}

			if variant == "" {
				variant, err = android.GetVariant(appManifestPath)

				if err != nil {
					return err
				}
			}

			appManifestPath = filepath.Join(appManifestPath, variant, "AndroidManifest.xml")

		}

		// Check to see if we need to read the manifest file due to missing options
		if apiKey == "" || versionCode == "" || version == "" {

			log.Info("Reading data from AndroidManifest.xml")
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

			if versionCode == "" {
				versionCode = manifestData.VersionCode
				log.Info("Using " + versionCode + " as version code from AndroidManifest.xml")
			}

			if version == "" {
				version = manifestData.VersionName
				log.Info("Using " + version + " as version name from AndroidManifest.xml")
			}
		}

		if sourceMapPath == "" {
			if utils.FileExists(filepath.Join(path, "android", "app", "build", "generated", "sourcemaps", "react")) {
				sourceMapPath = filepath.Join(path, "android", "app", "build", "generated", "sourcemaps", "react")
			} else if utils.FileExists(filepath.Join(path, "app", "build", "generated", "sourcemaps", "react")) {
				sourceMapPath = filepath.Join(path, "app", "build", "generated", "sourcemaps", "react")
			} else {
				return fmt.Errorf("unable to find the source map path. Please specify using `--source-map-path`")
			}

			sourceMapPath = filepath.Join(sourceMapPath, variant, "index.android.bundlePath.map")
		}

		if !utils.FileExists(sourceMapPath) {
			return fmt.Errorf(sourceMapPath + " doesn't exist on the system")
		}

		if bundlePath == "" {
			if utils.FileExists(filepath.Join(path, "android", "app", "build", "ASSETS", "createBundleReleaseJsAndAssets", "index.android.bundlePath")) {
				bundlePath = filepath.Join(path, "android", "app", "build", "ASSETS", "createBundleReleaseJsAndAssets", "index.android.bundlePath")
			} else if utils.FileExists(filepath.Join(path, "app", "build", "ASSETS", "createBundleReleaseJsAndAssets", "index.android.bundlePath")) {
				bundlePath = filepath.Join(path, "app", "build", "ASSETS", "createBundleReleaseJsAndAssets", "index.android.bundlePath")
			} else {
				return fmt.Errorf("unable to find the bundlePath path. Please specify using `--bundlePath-path`")
			}
		}

		if !utils.FileExists(bundlePath) {
			return fmt.Errorf(bundlePath + " doesn't exist on the system")
		}
	}

	log.Info("Uploading debug information for React Native Android")

	uploadOptions := utils.BuildReactNativeAndroidUploadOptions(apiKey, version, versionCode, codeBundleId, dev, projectRoot, overwrite)

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

	return nil
}
