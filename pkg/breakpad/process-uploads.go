package breakpad

import (
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"path/filepath"
)

func UploadBreakpadSymbols(
	fileList []string,
	apiKey string,
	projectRoot string,
	endpoint string,
	options options.CLI,
	logger log.Logger,
) error {
	fileFieldData := make(map[string]server.FileField)

	numberOfFiles := len(fileList)
	if numberOfFiles < 1 {
		logger.Info("No breakpad .sym files found to process")
		return nil
	}

	for _, file := range fileList {
		uploadOptions, err := utils.BuildBreakpadUploadOptions(apiKey, projectRoot, filepath.Base(file), options.Upload.Overwrite)

		if err != nil {
			return err
		}

		fileFieldData["symbol_file"] = server.LocalFile(file)
		apiKey := uploadOptions["api_key"]
		err = server.ProcessFileRequest(endpoint+"/breakpad-symbol"+"?api_key="+apiKey, uploadOptions, fileFieldData, file, options, logger)

		if err != nil {
			return err
		}
	}
	return nil
}
