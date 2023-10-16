package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

type UnityAndroid struct {
	AabPath       utils.Path  `help:"Path to Android AAB file to upload with your Unity symbols"`
	ApplicationId string      `help:"Module application identifier"`
	Path          utils.Paths `arg:"" name:"path" help:"(required) Path to Unity symbols zip file or directory to upload" type:"path"`
	ProjectRoot   string      `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	VersionCode   string      `help:"Module version code"`
	VersionName   string      `help:"Module version name"`
	BuildUuid     string      `help:"Module Build UUID"`
}

func ProcessUnityAndroid(apiKey string, aabPath string, applicationId string, versionCode string, buildUuid string, versionName string, projectRoot string, paths []string, endpoint string, failOnUploadError bool, retries int, timeout int, overwrite bool, dryRun bool) error {
	var err error
	var zipPath string
	var archList []string
	var symbolFileList []string
	var manifestData map[string]string

	for _, path := range paths {
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
			return fmt.Errorf(path + " is not a .symbols.zip file or containing directory")
		}
	}

	if aabPath != "" {
		log.Info("Extracting " + filepath.Base(aabPath) + " into a temporary directory")

		aabDir, err := utils.ExtractFile(aabPath, "aab")

		if err != nil {
			return err
		}

		defer os.RemoveAll(aabDir)

		manifestData, err = android.MergeUploadOptionsFromAabManifest(aabDir, apiKey, applicationId, buildUuid, versionCode, versionName)

		if err != nil {
			return err
		}

		err = ProcessAndroidAab(manifestData["apiKey"], manifestData["applicationId"], manifestData["buildUuid"], []string{aabDir}, projectRoot, manifestData["versionCode"], manifestData["versionName"], endpoint, failOnUploadError, retries, timeout, overwrite, dryRun)

		if err != nil {
			return err
		}
	}

	log.Info("Extracting " + filepath.Base(zipPath) + " into a temporary directory")

	if manifestData == nil {
		manifestData, _ = android.MergeUploadOptionsFromAabManifest("", apiKey, applicationId, buildUuid, versionCode, versionName)
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

	err = android.UploadAndroidNdk(symbolFileList, manifestData["apiKey"], manifestData["applicationId"], manifestData["versionName"], manifestData["versionCode"], projectRoot, overwrite, endpoint, timeout, dryRun, failOnUploadError)

	if err != nil {
		return err
	}

	return nil
}
