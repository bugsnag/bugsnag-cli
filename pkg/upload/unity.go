package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
	"path/filepath"
	"strings"
)

type UnityOptions struct {
	ApplicationId string            `help:"Module application identifier"`
	Arch          string            `help:"The architecture of the shared object that the symbols are for (e.g. x86, armeabi-v7a)."`
	Path          utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	ProjectRoot   string            `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	VersionCode   string            `help:"The version code for the application (Android only)."`
	VersionName   string            `help:"The version name of the application."`
}

func ProcessUnity(apiKey string, applicationId string, versionCode string, arch string, versionName string, projectRoot string, paths []string, endpoint string, failOnUploadError bool, retries int, timeout int, overwrite bool, dryRun bool) error {
	var archList []string
	var symbolFileList []string

	for _, path := range paths {

		if projectRoot == "" {
			projectRoot, _ = filepath.Split(path)
		}

		if strings.HasSuffix(path, ".symbols.zip") {
			tempDir, err := os.MkdirTemp("", "bugsnag-cli-unity-unpacking-*")

			if err != nil {
				return fmt.Errorf("error creating temporary working directory " + err.Error())
			}

			defer os.RemoveAll(tempDir)

			log.Info("Extracting " + filepath.Base(path) + " to " + tempDir)

			err = utils.Unzip(path, tempDir)

			if err != nil {
				return err
			}

			if arch == "" {
				archList, err = utils.BuildFolderList([]string{tempDir})

				if err != nil {
					return err
				}
			} else {
				archList = []string{arch}
			}

			for _, arch := range archList {
				soPath := filepath.Join(tempDir, arch)

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
		} else if filepath.Ext(path) == ".so" {
			symbolFileList = append(symbolFileList, path)
		} else {
			return fmt.Errorf(path + " is not an Unity symbols file")
		}
	}

	if apiKey != "" {
		log.Info("Using " + apiKey + " as the API key")
	}

	if applicationId == "" {
		return fmt.Errorf("Application ID not provided, please specify using `--application-id`")
	} else {
		log.Info("Using " + applicationId + " as the application ID")
	}

	if versionName == "" && versionName == "" {
		return fmt.Errorf("Version Code or version name not provided, please specify using `--version-code` or `--version-name`")
	}

	if versionCode != "" {
		log.Info("Using " + versionCode + " as the version code")
	}

	if versionName != "" {
		log.Info("Using " + versionName + " as the version Name")
	}

	if projectRoot != "" {
		log.Info("Using " + projectRoot + " as the project root")
	}

	numberOfFiles := len(symbolFileList)

	if numberOfFiles < 1 {
		log.Info("No symbol files found to process")
		return nil
	}

	for _, file := range symbolFileList {
		uploadOptions, err := utils.BuildAndroidNDKUploadOptions(apiKey, applicationId, versionName, versionCode, projectRoot, filepath.Base(file), overwrite)

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
