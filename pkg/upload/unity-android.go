package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

type UnityAndroidOptions struct {
	AabPath       string            `help:"Path to Android AAB file"`
	ApplicationId string            `help:"Module application identifier"`
	Arch          string            `help:"The architecture of the shared object that the symbols are for (e.g. x86, armeabi-v7a)."`
	Path          utils.UploadPaths `arg:"" name:"path" help:"(required) Path to Unity symbols zip file or directory to upload" type:"path"`
	ProjectRoot   string            `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	VersionCode   string            `help:"Module version code"`
	VersionName   string            `help:"Module version name"`
	BuildUuid     string            `help:"Module Build UUID"`
}

func ProcessUnityAndroid(apiKey string, aabPath string, applicationId string, versionCode string, buildUuid string, arch string, versionName string, projectRoot string, paths []string, endpoint string, failOnUploadError bool, retries int, timeout int, overwrite bool, dryRun bool) error {
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
				aabPath, err = utils.FindLatestFileWithSuffix(path, ".aab")

				if err != nil {
					return err
				}
			}
		} else if strings.HasSuffix(path, ".symbols.zip") {
			zipPath = path

			if aabPath == "" {
				buildDirectory := filepath.Dir(path)

				aabPath, err = utils.FindLatestFileWithSuffix(buildDirectory, ".aab")

				if err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf("unsupported file paths provided. Please specify the Unity `symbols.zip` file or build directory.")
		}
	}

	log.Info("Extracting " + filepath.Base(aabPath) + " into a temporary directory")

	aabDir, err := utils.ExtractFile(aabPath, "aab")

	if err != nil {
		return err
	}

	defer os.RemoveAll(aabDir)

	if applicationId == "" || buildUuid == "" || versionCode == "" || versionName == "" {
		log.Info("Reading data from AndroidManifest.xml")

		manifestData, err = android.GetUploadOptionsFromAabManifest(aabDir, apiKey, applicationId, buildUuid, versionCode, versionName)

		if err != nil {
			return err
		}
	}

	err = ProcessAndroidAab(manifestData["apiKey"], manifestData["applicationId"], manifestData["buildUuid"], []string{aabDir}, projectRoot, manifestData["versionCode"], manifestData["versionName"], endpoint, failOnUploadError, retries, timeout, overwrite, dryRun)

	if err != nil {
		return err
	}

	log.Info("Extracting " + filepath.Base(zipPath) + " into a temporary directory")

	unityDir, err := utils.ExtractFile(zipPath, "unity-android")

	if err != nil {
		return err
	}

	defer os.RemoveAll(unityDir)

	if arch == "" {
		archList, err = utils.BuildFolderList([]string{unityDir})
		if err != nil {
			return err
		}
	} else {
		archList = []string{arch}
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

	numberOfFiles := len(symbolFileList)

	if numberOfFiles < 1 {
		log.Info("No symbol files found in " + zipPath)
		return nil
	}

	for _, file := range symbolFileList {
		uploadOptions, err := utils.BuildAndroidNDKUploadOptions(manifestData["apiKey"], manifestData["applicationId"], manifestData["versionName"], manifestData["versionCode"], projectRoot, filepath.Base(file), overwrite)

		if err != nil {
			return err
		}

		fileFieldData := make(map[string]string)
		fileFieldData["soFile"] = file

		err = server.ProcessRequest(endpoint+"/ndk-symbol", uploadOptions, fileFieldData, timeout, file, dryRun)

		if err != nil {
			if numberOfFiles > 1 && failOnUploadError {
				return err
			} else {
				log.Warn(err.Error())
			}
		} else {
			log.Success("Uploaded " + filepath.Base(file))
		}
	}

	return nil
}
