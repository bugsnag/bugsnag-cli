package android

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

func ProcessNdk(apiKey string, configuration string, outputPath string, appId string, versionCode string, versionName string, objCopyPath string, projectRoot string, overwrite bool, timeout int, endpoint string, failOnUploadError bool) error {

	log.Info("Building file list for configuration: " + configuration)

	symbolPath := []string{filepath.Join(outputPath, "BUNDLE-METADATA", "com.android.tools.build.debugsymbols")}

	fileList, err := utils.BuildFileList(symbolPath)

	if err != nil {
		return fmt.Errorf("error building file list for configuration: " + configuration)
	}

	log.Info("Processing NDK files for configuration: " + configuration)

	numberOfFiles := len(fileList)

	if numberOfFiles < 1 {
		log.Info("No files to process for configuration: " + configuration)
		return nil
	}

	for _, file := range fileList {

		if strings.HasSuffix(file, ".so.sym") {

			log.Info("Extracting debug info from " + filepath.Base(file) + " using objcopy")
			outputFile, err := Objcopy(objCopyPath, file)

			if err != nil {
				return fmt.Errorf("failed to process file, " + file + " using objcopy. " + err.Error())
			}

			log.Info(outputFile)

			log.Info("Uploading debug information for " + filepath.Base(file))

			uploadOptions := utils.BuildAndroidNDKUploadOptions(apiKey, appId, versionName, versionCode, projectRoot, filepath.Base(file), overwrite)

			fileFieldData := make(map[string]string)
			fileFieldData["soFile"] = outputFile

			requestStatus := server.ProcessRequest(endpoint+"/ndk-symbol", uploadOptions, fileFieldData, timeout)

			if requestStatus != nil {
				if numberOfFiles > 1 && failOnUploadError {
					return requestStatus
				} else {
					log.Warn(requestStatus.Error())
				}
			} else {
				log.Success(filepath.Base(file) + " uploaded")
			}
		}

	}
	return nil
}
