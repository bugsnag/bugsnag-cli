package upload

import (
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type DiscoverAndUploadAny struct {
	Path          utils.UploadPaths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path"`
	UploadOptions map[string]string `help:"additional arguments to pass to the upload request" mapsep:","`
}

func All(paths []string, options map[string]string, endpoint string, timeout int, retries int, overwrite bool,
	apiKey string, failOnUploadError bool) error {

	var fileFieldName string

	// Build the file list from the path(s)
	log.Info("building file list...")

	fileList, err := utils.BuildFileList(paths)
	numberOfFiles := len(fileList)

	if err != nil {
		log.Error(" error building file list", 1)
	}

	log.Info("File list built..")

	// Build UploadOptions list
	uploadOptions := make(map[string]string)

	uploadOptions["apiKey"] = apiKey

	if overwrite {
		uploadOptions["overwrite"] = "true"
	}

	for key, value := range options {
		uploadOptions[key] = value
	}

	if uploadOptions["fileNameField"] != "" {
		fileFieldName = uploadOptions["fileNameField"]
		delete(uploadOptions, "fileNameField")
	} else {
		fileFieldName = "file"
	}

	for _, file := range fileList {
		requestStatus := server.ProcessRequest(endpoint, uploadOptions, fileFieldName, file, timeout)

		if requestStatus != nil {
			if numberOfFiles > 1 && failOnUploadError {
				return requestStatus
			} else {
				log.Warn(requestStatus.Error())
			}
		} else {
			log.Success(file)
		}
	}

	return nil
}
