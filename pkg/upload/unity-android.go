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

func ProcessUnityAndroid(globalOptions options.CLI, endpoint string, logger log.Logger) error {
	unityOptions := globalOptions.Upload.UnityAndroid
	var err error
	var zipPath string
	var archList []string
	var symbolFileList []string
	var manifestData map[string]string
	var aabPath = string(unityOptions.AabPath)
	var aabUploaded bool

	for _, path := range unityOptions.Path {
		if utils.IsDir(path) {
			zipPath, err = utils.FindLatestFileWithSuffix(path, ".symbols.zip")

			if err != nil {
				return err
			}

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
		err = ProcessAndroidAab(globalOptions, endpoint, logger)

		aabUploaded = true

		if err != nil {
			if strings.Contains(err.Error(), "No NDK (.so) or Proguard (mapping.txt) files detected for upload.") {
				aabUploaded = false
			} else {
				return err
			}
		}
	}

	if zipPath == "" && !aabUploaded {
		return fmt.Errorf("No .symbols.zip or .aab files detected for upload")
	} else if zipPath != "" {

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
			endpoint,
			globalOptions,
			logger,
		)

		if err != nil {
			return err
		}
	}
	return nil
}
