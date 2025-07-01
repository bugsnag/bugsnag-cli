package upload

import (
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

func All(options options.CLI, endpoint string, logger log.Logger,
) error {
	allOptions := options.Upload.All
	fileList, err := utils.BuildFileList(allOptions.Path)

	if err != nil {
		logger.Fatal("Error building file list")
	}

	// Build UploadOptions list
	uploadOptions := make(map[string]string)

	if options.Upload.Overwrite {
		uploadOptions["overwrite"] = "true"
	}

	for key, value := range allOptions.UploadOptions {
		uploadOptions[key] = value
	}

	for _, file := range fileList {

		fileFieldData := make(map[string]server.FileField)

		if uploadOptions["fileNameField"] != "" {
			fileFieldData[uploadOptions["fileNameField"]] = server.LocalFile(file)
			delete(uploadOptions, "fileNameField")
		} else {
			fileFieldData["file"] = server.LocalFile(file)
		}

		err := server.ProcessFileRequest(options.ApiKey, endpoint, uploadOptions, fileFieldData, file, options, logger)

		if err != nil {

			return err
		}
	}

	return nil
}
