package android

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

func ProcessNdk(apiKey string, variant string, outputPath string, aabManifestData map[string]string, objCopyPath string, projectRoot string, overwrite bool, timeout int, endpoint string, failOnUploadError bool) error {

	log.Info("Building file list for variant: " + variant)

	symbolPath := []string{filepath.Join(outputPath, "BUNDLE-METADATA", "com.android.tools.build.debugsymbols")}

	fileList, err := utils.BuildFileList(symbolPath)

	if err != nil {
		return fmt.Errorf("error building file list for variant: " + variant)
	}

	log.Info("Processing NDK files for variant: " + variant)

	numberOfFiles := len(fileList)

	if numberOfFiles < 1 {
		log.Info("No files to process for variant: " + variant)
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

			uploadOptions := utils.BuildAndroidNDKUploadOptions(apiKey, aabManifestData["package"], aabManifestData["versionName"], aabManifestData["versionCode"], projectRoot, filepath.Base(file), overwrite)

			requestStatus := server.ProcessRequest(endpoint+"/ndk-symbol", uploadOptions, "soFile", outputFile, timeout)

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
