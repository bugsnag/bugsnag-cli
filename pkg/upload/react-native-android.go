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
		if bundlePath == "" {
			bundlePath = filepath.Join(path, "android", "app", "build", "ASSETS", "createBundleReleaseJsAndAssets", "index.android.bundle")
			if !utils.FileExists(bundlePath) {
				bundlePath = filepath.Join(path, "app", "build", "ASSETS", "createBundleReleaseJsAndAssets", "index.android.bundle")
				if !utils.FileExists(bundlePath) {
					return fmt.Errorf("unable to find index.android.bundle within " + path)
				}
			}
		}

		if sourceMapPath == "" {
			sourceMapPath = filepath.Join(path, "android", "app", "build", "generated", "sourcemaps", "react")
			if !utils.IsDir(sourceMapPath) {
				sourceMapPath = filepath.Join(path, "app", "build", "generated", "sourcemaps", "react")
				if !utils.IsDir(sourceMapPath) {
					return fmt.Errorf("unable to find index.android.bundle within " + path)
				}
			}

			if variant == "" {
				variant, err = android.GetVariant(sourceMapPath)

				if err != nil {
					return err
				}
			}

			sourceMapPath = filepath.Join(sourceMapPath, variant, "index.android.bundle.map")

			if !utils.FileExists(sourceMapPath) {
				return fmt.Errorf("unable to find index.android.bundle within " + path)
			}
		}

		if projectRoot == "" {
			projectRoot = path
		}

		if apiKey == "" || version == "" || versionCode == "" {
			if appManifestPath == "" {
				appManifestPath = filepath.Join(path, "android", "app", "build", "intermediates", "merged_manifests")
				if !utils.IsDir(appManifestPath) {
					appManifestPath = filepath.Join(path, "app", "build", "intermediates", "merged_manifests")
					if !utils.IsDir(appManifestPath) {
						return fmt.Errorf("unable to find AndroidManfiest.xml within " + path)
					}
				}

				if variant == "" {
					variant, err = android.GetVariant(appManifestPath)

					if err != nil {
						return fmt.Errorf(err.Error())
					}
				}

				appManifestPath = filepath.Join(appManifestPath, variant, "AndroidManifest.xml")
			}

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
	}

	return nil
}

//./bin/arm64-macos/bugsnag-cli upload react-native-android --source-map=test/testdata/react-native/android/app/build/generated/sourcemaps/react/release/index.android.bundle.map --api-key=676341688653c170796009e05430fc60 --overwrite
