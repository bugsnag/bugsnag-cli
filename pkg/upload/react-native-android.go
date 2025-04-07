package upload

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

func ProcessReactNativeAndroid(options options.CLI, endpoint string, logger log.Logger) error {
	androidOptions := options.Upload.ReactNativeAndroid
	var err error
	var uploadOptions map[string]string
	var rootDirPath string
	var variantDirName string
	var bundleDirPath string
	var variantFileFormat string
	var appManifestPathExpected string

	for _, path := range androidOptions.Path {

		buildDirPath := filepath.Join(path, "android", "app", "build")
		rootDirPath = path
		if !utils.FileExists(buildDirPath) {
			buildDirPath = filepath.Join(path, "app", "build")
			if utils.FileExists(buildDirPath) {
				rootDirPath = filepath.Join(path, "..")
			} else if androidOptions.ReactNative.Bundle == "" || androidOptions.ReactNative.SourceMap == "" {
				return fmt.Errorf("unable to find bundle files or source maps in within %s", path)
			}
		}

		if androidOptions.ProjectRoot == "" {
			androidOptions.ProjectRoot = rootDirPath
		}

		if androidOptions.ReactNative.Bundle == "" {
			if utils.IsDir(filepath.Join(buildDirPath, "generated", "assets", "react")) {
				// RN version < 0.70 - generated/assets/react/<androidOptions.Android.Variant>/index.android.bundle
				bundleDirPath = filepath.Join(buildDirPath, "generated", "assets", "react")
			} else if utils.IsDir(filepath.Join(buildDirPath, "ASSETS")) {
				// RN versions < 0.72 - ASSETS/createBundle<androidOptions.Android.Variant>JsAndAssets/index.android.bundle
				bundleDirPath = filepath.Join(buildDirPath, "ASSETS")
				variantFileFormat = "createBundle%sJsAndAssets"
			} else if utils.IsDir(filepath.Join(buildDirPath, "generated", "assets")) {
				// RN versions >= 0.72 - generated/assets/<androidOptions.Android.Variant>/index.android.bundle
				bundleDirPath = filepath.Join(buildDirPath, "generated", "assets")
				variantFileFormat = "createBundle%sJsAndAssets"
			} else {
				return fmt.Errorf("unable to find index.android.bundle in your project, please specify the path using --bundle-path")
			}

			if bundleDirPath != "" {
				if androidOptions.Android.Variant == "" {
					variantDirName, err = android.GetVariantDirectory(bundleDirPath)
					if err != nil {
						return err
					}
				} else {
					if variantFileFormat != "" {
						variantDirName = fmt.Sprintf(variantFileFormat,
							cases.Title(language.Und, cases.NoLower).String(androidOptions.Android.Variant))

					} else {
						variantDirName = androidOptions.Android.Variant
					}
				}
				androidOptions.ReactNative.Bundle = filepath.Join(bundleDirPath, variantDirName, "index.android.bundle")
			}
		}

		if !utils.FileExists(androidOptions.ReactNative.Bundle) {
			return fmt.Errorf("unable to find index.android.bundle at %s", androidOptions.ReactNative.Bundle)
		}

		if androidOptions.ReactNative.SourceMap == "" {
			sourceMapDirPath := filepath.Join(buildDirPath, "generated", "sourcemaps", "react")

			if androidOptions.Android.Variant == "" {
				androidOptions.Android.Variant, err = android.GetVariantDirectory(sourceMapDirPath)
				if err != nil {
					return err
				}
			}

			androidOptions.ReactNative.SourceMap = filepath.Join(sourceMapDirPath, androidOptions.Android.Variant, "index.android.bundle.map")
		} else {
			if androidOptions.Android.Variant == "" {
				// Set androidOptions.Android.Variant based off the source map file location
				sourceMapDirPath := filepath.Join(androidOptions.ReactNative.SourceMap, "..", "..")

				if filepath.Base(sourceMapDirPath) == "react" {
					androidOptions.Android.Variant, err = android.GetVariantDirectory(sourceMapDirPath)
					if err != nil {
						return err
					}
				}
			}
		}

		if !utils.FileExists(androidOptions.ReactNative.SourceMap) {
			return fmt.Errorf("unable to find index.android.bundle at %s", androidOptions.ReactNative.SourceMap)
		}

		if androidOptions.Android.AppManifest == "" {
			// RN versions <= 0.74 intermediates/merged_manifests/<androidOptions.Android.Variant>/AndroidManifest.xml"
			appManifestPathExpected = filepath.Join(buildDirPath, "intermediates", "merged_manifests", androidOptions.Android.Variant, "AndroidManifest.xml")
			if utils.FileExists(appManifestPathExpected) {
				androidOptions.Android.AppManifest = appManifestPathExpected
				logger.Debug(fmt.Sprintf("Found app manifest at: %s", androidOptions.Android.AppManifest))
			} else {
				// RN versions > 0.74 "intermediates/merged_manifests/androidOptions.Android.Variant, <androidOptions.Android.Variant>/AndroidManifest.xml"
				appManifestPathExpected = filepath.Join(buildDirPath, "intermediates", "merged_manifests", androidOptions.Android.Variant, "process"+cases.Title(language.English).String(androidOptions.Android.Variant)+"Manifest", "AndroidManifest.xml")
				if utils.FileExists(appManifestPathExpected) {
					androidOptions.Android.AppManifest = appManifestPathExpected
					logger.Debug(fmt.Sprintf("Found app manifest at: %s", androidOptions.Android.AppManifest))
				} else {
					logger.Debug(fmt.Sprintf("No app manifest found at: %s", appManifestPathExpected))
				}
			}
		}

		if androidOptions.Android.AppManifest != "" && (options.ApiKey == "" || androidOptions.ReactNative.VersionName == "" || androidOptions.Android.VersionCode == "") {

			manifestData, err := android.ParseAndroidManifestXML(androidOptions.Android.AppManifest)

			if err != nil {
				return err
			}

			if options.ApiKey == "" {
				for key, value := range manifestData.Application.MetaData.Name {
					if value == "com.bugsnag.android.API_KEY" {
						options.ApiKey = manifestData.Application.MetaData.Value[key]
					}
				}
				logger.Debug(fmt.Sprintf("Using %s as API key from AndroidManifest.xml", options.ApiKey))
			}

			// If we've not passed --code-bundle-id, proceed to populate versionName and versionCode from AndroidManifest.xml
			if androidOptions.ReactNative.CodeBundleId == "" {
				if androidOptions.ReactNative.VersionName == "" {
					androidOptions.ReactNative.VersionName = manifestData.VersionName
					logger.Debug(fmt.Sprintf("Using %s as version name from AndroidManifest.xml", androidOptions.ReactNative.VersionName))
				}

				if androidOptions.Android.VersionCode == "" {
					androidOptions.Android.VersionCode = manifestData.VersionCode
					logger.Debug(fmt.Sprintf("Using %s as version code from AndroidManifest.xml", androidOptions.Android.VersionCode))
				}
			}
		}

		uploadOptions, err = utils.BuildReactNativeUploadOptions(options.ApiKey, androidOptions.ReactNative.VersionName, androidOptions.Android.VersionCode, androidOptions.ReactNative.CodeBundleId, androidOptions.ReactNative.Dev, androidOptions.ProjectRoot, options.Upload.Overwrite, "android")

		if err != nil {
			return err
		}

		fileFieldData := make(map[string]server.FileField)
		fileFieldData["sourceMap"] = server.LocalFile(androidOptions.ReactNative.SourceMap)
		fileFieldData["bundle"] = server.LocalFile(androidOptions.ReactNative.Bundle)

		err = server.ProcessFileRequest(endpoint+"/react-native-source-map", uploadOptions, fileFieldData, androidOptions.ReactNative.SourceMap, options, logger)

		if err != nil {

			return err
		}
	}

	return nil
}
