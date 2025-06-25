package android

import (
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

func UploadAndroidNdk(
	file string,
	apiKey string,
	applicationId string,
	versionName string,
	versionCode string,
	projectRoot string,
	endpoint string,
	options options.CLI,
	overwrite bool,
	logger log.Logger,
) error {
	fileFieldData := make(map[string]server.FileField)

	uploadOptions, err := utils.BuildAndroidNDKUploadOptions(apiKey, applicationId, versionName, versionCode, projectRoot, filepath.Base(file), overwrite)

	if err != nil {
		return err
	}

	fileFieldData["soFile"] = server.LocalFile(file)

	err = server.ProcessFileRequest(endpoint+"/ndk-symbol", uploadOptions, fileFieldData, file, options, logger)

	if err != nil {
		return err
	}

	return nil
}
