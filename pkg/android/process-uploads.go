package android

import (
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"path/filepath"
	"fmt"
)

func UploadAndroidNdk(
	fileList []string,
	apiKey string,
	applicationId string,
	versionName string,
	versionCode string,
	projectRoot string,
	overwrite bool,
	endpoint string,
	timeout int,
	retries int,
	dryRun bool,
) error {
	fileFieldData := make(map[string]string)

	numberOfFiles := len(fileList)

	if numberOfFiles < 1 {
		log.Info("No NDK files found to process")
		return nil
	}

	for _, file := range fileList {
		uploadOptions, err := utils.BuildAndroidNDKUploadOptions(apiKey, applicationId, versionName, versionCode, projectRoot, filepath.Base(file), overwrite)

		if err != nil {
			return err
		}

		fileFieldData["soFile"] = file

		err = server.ProcessFileRequest(fmt.Sprintf("%s/ndk-symbol", endpoint), uploadOptions, fileFieldData, timeout, retries, file, dryRun)

		if err != nil {
			return err
		}
	}

	return nil
}
