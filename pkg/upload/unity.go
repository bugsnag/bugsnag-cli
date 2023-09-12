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
	Path        utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	AppId       string            `help:"the Android applicationId for this application."`
	VersionCode string            `help:"the Android versionCode for this application release."`
	Arch        string            `help:"the architecture of the shared object that the symbols are for (e.g. x86, armeabi-v7a)."`
	VersionName string            `help:"the Android versionName for this application release."`
	ProjectRoot string            `help:"path to remove from the beginning of the filenames in the mapping file"`
}

func ProcessUnity(apiKey string, applicationId string, versionCode string, arch string, versionName string, projectRoot string, paths []string, endpoint string, retries int, timeout int, overwrite bool, dryRun bool) error {
	var archList []string

	tempDir, err := os.MkdirTemp("", "bugsnag-cli-unity-unpacking-*")

	if err != nil {
		return fmt.Errorf("error creating temporary working directory " + err.Error())
	}

	defer os.RemoveAll(tempDir)

	// extract symbols files
	for _, path := range paths {
		if strings.HasSuffix(path, ".symbols.zip") {
			log.Info("Extracting " + filepath.Base(path) + " to " + tempDir)
			err = utils.Unzip(path, tempDir)

			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf(path + " is not an Unity symbols file")
		}
	}

	// Check to see if we have been provided an arch type.
	if arch == "" {
		// Build a list of arch types if not.
		archList, err = utils.BuildFolderList([]string{tempDir})

		if err != nil {
			return err
		}
	} else {
		archList = []string{arch}
	}

	if applicationId == "" {
		return fmt.Errorf("Missing Application ID")
	}

	if versionCode == "" {
		return fmt.Errorf("Missing Version Code")
	}

	// process files for each arch type
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

			log.Info("Uploading debug information for " + filepath.Base(file))

			if err != nil {
				return fmt.Errorf("failed to process file, " + file + " using objcopy : " + err.Error())
			}

			uploadOptions, err := utils.BuildAndroidNDKUploadOptions(apiKey, applicationId, versionName, versionCode, projectRoot, filepath.Base(file), overwrite)

			if err != nil {
				return err
			}

			fileFieldData := make(map[string]string)
			fileFieldData["soSymbolFile"] = file
			err = server.ProcessRequest(endpoint+"/ndk-symbol", uploadOptions, fileFieldData, timeout, file, dryRun)
			if err != nil {
				return err
			} else {
				log.Success(filepath.Base(file) + " uploaded")
			}

		}
	}

	return nil
}
