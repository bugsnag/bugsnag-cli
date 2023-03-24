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
	Path            utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path" default:"."`
	AppManifestPath string            `help:"(required) Path to directory or file to upload" type:"path"`
	CodeBundleId    string            `help:"A unique identifier to identify a code bundle release when using tools like CodePush"`
	Dev             bool              `help:"Indicates whether the application is a debug or release build"`
	SourceMapPath   string            `help:"Path to the source map file" type:"path"`
	BundlePath      string            `help:"Path to the bundle file" type:"path"`
	ProjectRoot     string            `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
}

func ProcessReactNativeAndroid(paths []string, appManifestPath string, appVersion string, appVersionCode string, codeBundleId string, dev bool, sourceMapPath string, bundlePath string, projectRoot string, endpoint string, timeout int, retries int, overwrite bool, apiKey string) error {

	for _, path := range paths {
		log.Info(path)

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

			variants, err := android.BuildVariantsList(appManifestPath)

			if err != nil {
				return err
			}

			if len(variants) > 1 {
				return fmt.Errorf("more than one variant found. Please specify using `--app-manifest-path`")
			}

			appManifestPath = filepath.Join(appManifestPath, variants[0], "AndroidManifest.xml")

		}

		if !utils.FileExists(appManifestPath) {
			return fmt.Errorf(appManifestPath + " doesn't exist on the system")
		}

		androidManifestData, err := android.ParseAndroidManifestXML(appManifestPath)

		if err != nil {
			return err
		}

		if appVersion == "" {
			log.Info("Setting app version from " + appManifestPath)
			appVersion = androidManifestData.VersionName
		}

		if appVersionCode == "" {
			log.Info("Setting app version code from " + appManifestPath)
			appVersionCode = androidManifestData.VersionCode
		}

		if sourceMapPath == "" {
			if utils.FileExists(filepath.Join(path, "android", "app", "build", "generated", "sourcemaps", "react")) {
				sourceMapPath = filepath.Join(path, "android", "app", "build", "generated", "sourcemaps", "react")
			} else if utils.FileExists(filepath.Join(path, "app", "build", "generated", "sourcemaps", "react")) {
				sourceMapPath = filepath.Join(path, "app", "build", "generated", "sourcemaps", "react")
			} else {
				return fmt.Errorf("unable to find the source map path. Please specify using `--source-map-path`")
			}

			variants, err := android.BuildVariantsList(sourceMapPath)

			if err != nil {
				return err
			}

			if len(variants) > 1 {
				return fmt.Errorf("more than one variant found. Please specify using `--source-map-path`")
			}

			sourceMapPath = filepath.Join(sourceMapPath, variants[0], "index.android.bundle.map")

		}

		if !utils.FileExists(sourceMapPath) {
			return fmt.Errorf(sourceMapPath + " doesn't exist on the system")
		}

		if bundlePath == "" {
			if utils.FileExists(filepath.Join(path, "android", "app", "build", "ASSETS", "createBundleReleaseJsAndAssets", "index.android.bundle")) {
				bundlePath = filepath.Join(path, "android", "app", "build", "ASSETS", "createBundleReleaseJsAndAssets", "index.android.bundle")
			} else if utils.FileExists(filepath.Join(path, "app", "build", "ASSETS", "createBundleReleaseJsAndAssets", "index.android.bundle")) {
				bundlePath = filepath.Join(path, "app", "build", "ASSETS", "createBundleReleaseJsAndAssets", "index.android.bundle")
			} else {
				return fmt.Errorf("unable to find the bundle path. Please specify using `--bundle-path`")
			}
		}

		if !utils.FileExists(bundlePath) {
			return fmt.Errorf(bundlePath + " doesn't exist on the system")
		}
	}

	log.Info("Uploading debug information for React Native Android")

	uploadOptions := utils.BuildReactNativeAndroidUploadOptions(apiKey, appVersion, appVersionCode, codeBundleId, dev, projectRoot, overwrite)

	fileFieldData := make(map[string]string)
	fileFieldData["sourceMap"] = sourceMapPath
	fileFieldData["bundle"] = bundlePath

	requestStatus := server.ProcessRequest(endpoint, uploadOptions, fileFieldData, timeout)

	if requestStatus != nil {
		return requestStatus

	}

	return nil
}
