package upload

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// ProcessUnityAndroid processes Unity Android symbols and AAB files.
//
// This function searches for Unity Android symbols.zip files and AAB files in the specified paths,
// extracts the necessary data, and uploads the symbols.
// It handles both the symbols.zip and AAB files, extracting architecture-specific symbols
// and merging metadata from the AAB manifest if available.
//
// Parameters:
//   - globalOptions: CLI options containing Unity Android upload settings.
//   - logger: Logger instance for debug and error output.
//
// Returns:
//   - error: non-nil if an error occurs during processing or uploading.
func ProcessUnityAndroid(globalOptions options.CLI, logger log.Logger) error {
	unityOptions := globalOptions.Upload.UnityAndroid
	var zipPath string
	var archList []string
	var symbolFileList []string
	var manifestData map[string]string
	var aabPath = string(unityOptions.AabPath)

	for _, path := range unityOptions.Path {
		if utils.IsDir(path) {
			zipPath, _ = utils.FindLatestFileWithSuffix(path, ".symbols.zip")

			if aabPath == "" {
				aabPath, _ = utils.FindLatestFileWithSuffix(path, ".aab")
			}
		} else if strings.HasSuffix(path, ".symbols.zip") {
			zipPath = path

			if aabPath == "" {
				buildDirectory := filepath.Dir(path)

				aabPath, _ = utils.FindLatestFileWithSuffix(buildDirectory, ".aab")
			}
		} else {
			return fmt.Errorf("%s is not a .symbols.zip file or containing directory", path)
		}
	}

	if aabPath != "" {

		logger.Debug(fmt.Sprintf("Extracting %s into a temporary directory", filepath.Base(aabPath)))

		aabDir, err := utils.ExtractFile(aabPath, "aab")

		if err != nil {
			return err
		}

		defer os.RemoveAll(aabDir)

		manifestData, err = android.MergeUploadOptionsFromAabManifest(aabDir, globalOptions.ApiKey, unityOptions.ApplicationId, unityOptions.BuildUuid, unityOptions.NoBuildUuid, unityOptions.VersionCode, unityOptions.VersionName, logger)

		if err != nil {
			return err
		}

		globalOptions.ApiKey = manifestData["apiKey"]
		globalOptions.Upload.AndroidAab = options.AndroidAabMapping{
			ApplicationId: manifestData["applicationId"],
			BuildUuid:     manifestData["buildUuid"],
			NoBuildUuid:   unityOptions.NoBuildUuid,
			Path:          []string{aabDir},
			ProjectRoot:   unityOptions.ProjectRoot,
			VersionCode:   manifestData["versionCode"],
			VersionName:   manifestData["versionName"],
		}
		err = ProcessAndroidAab(globalOptions, logger)

		if err != nil {
			return err
		}
	}

	if zipPath == "" {
		logger.Info("No Unity Android symbols.zip file found, skipping")
	} else {
		logger.Debug(fmt.Sprintf("Extracting %s into a temporary directory", filepath.Base(zipPath)))

		if manifestData == nil {
			manifestData, _ = android.MergeUploadOptionsFromAabManifest("", globalOptions.ApiKey, unityOptions.ApplicationId, unityOptions.BuildUuid, unityOptions.NoBuildUuid, unityOptions.VersionCode, unityOptions.VersionName, logger)
		}

		unityDir, err := utils.ExtractFile(zipPath, "unity-android")

		if err != nil {
			return err
		}

		defer os.RemoveAll(unityDir)

		archList, err = utils.BuildDirectoryList([]string{unityDir})

		if err != nil {
			return err
		}

		for _, arch := range archList {
			soPath := filepath.Join(unityDir, arch)
			fileList, err := utils.BuildFileList([]string{soPath})
			if err != nil {
				return err
			}
			for _, file := range fileList {
				if filepath.Base(file) == "libil2cpp.sym.so" && utils.ContainsString(fileList, "libil2cpp.dbg.so") {
					continue
				}
				symbolFileList = append(symbolFileList, file)
			}
		}

		err = android.UploadAndroidNdk(
			symbolFileList,
			manifestData["apiKey"],
			manifestData["applicationId"],
			manifestData["versionName"],
			manifestData["versionCode"],
			unityOptions.ProjectRoot,
			globalOptions,
			logger,
		)

		if err != nil {
			return err
		}
	}
	return nil
}
