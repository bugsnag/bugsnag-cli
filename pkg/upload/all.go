package upload

import (
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// All processes and uploads all files specified by the upload options.
//
// It builds a list of files from the given path, applies any specified upload options
// (such as overwriting existing files), and uploads each file individually.
// The field name for the file upload can be customized via the "fileNameField" option.
//
// Parameters:
//   - options: CLI options containing upload settings and API key.
//   - logger: logger instance for logging messages.
//
// Returns:
//   - error: non-nil if file list building or any file upload fails.
func All(options options.CLI, logger log.Logger) error {
	allOptions := options.Upload.All
	fileList, err := utils.BuildFileList(allOptions.Path)
	if err != nil {
		logger.Fatal("Error building file list")
	}

	// Build UploadOptions map from CLI options
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

		err := server.ProcessFileRequest(options.ApiKey, "", uploadOptions, fileFieldData, file, options, logger)
		if err != nil {
			return err
		}
	}

	return nil
}
