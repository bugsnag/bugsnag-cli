package upload

import (
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"path/filepath"
)

type DiscoverAndUploadAny struct {
	Path          utils.Paths       `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	UploadOptions map[string]string `help:"Additional arguments to pass to the upload request" mapsep:","`
}

func All(
	paths []string,
	options map[string]string,
	endpoint string,
	timeout int,
	retries int,
	overwrite bool,
	apiKey string,
	dryRun bool,
) error {

	// Build the file list from the path(s)
	log.Info("building file list...")

	fileList, err := utils.BuildFileList(paths)

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

	for _, file := range fileList {

		fileFieldData := make(map[string]string)

		if uploadOptions["fileNameField"] != "" {
			fileFieldData[uploadOptions["fileNameField"]] = file
			delete(uploadOptions, "fileNameField")
		} else {
			fileFieldData["file"] = file
		}

		requestStatus := server.ProcessFileRequest(endpoint, uploadOptions, fileFieldData, timeout, retries, file, dryRun)

		if requestStatus != nil {
			return err
		} else {
			log.Success("Uploaded " + filepath.Base(file))
		}
	}

	return nil
}
