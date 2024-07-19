package android

import (
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"path/filepath"
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
	logger log.Logger,
) error {
	fileFieldData := make(map[string]server.FileField)

	numberOfFiles := len(fileList)

	if numberOfFiles < 1 {
		logger.Info("No NDK files found to process")
		return nil
	}

	for _, file := range fileList {
		uploadOptions, err := utils.BuildAndroidNDKUploadOptions(apiKey, applicationId, versionName, versionCode, projectRoot, filepath.Base(file), overwrite)

		if err != nil {
			return err
		}

		fileFieldData["soFile"] = server.LocalFile(file)

		err = server.ProcessFileRequest(endpoint+"/ndk-symbol", uploadOptions, fileFieldData, timeout, retries, file, dryRun, logger)

		if err != nil {
			return err
		}
	}

	return nil
}
