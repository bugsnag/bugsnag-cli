package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"io"
	"strconv"
)

type DiscoverAndUploadAny struct {
	Path             utils.UploadPaths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path"`
	UploadOptions map[string]string `help:"additional arguments to pass to the upload request" mapsep:","`
}

func All(paths []string, options map[string]string, endpoint string, timeout int, retries int, overwrite bool,
	apiKey string) error {

	var fileFieldName string

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

	uploadOptions["retries"] = strconv.Itoa(retries)

	if uploadOptions["fileNameField"] != "" {
		fileFieldName = uploadOptions["fileNameField"]
		delete(uploadOptions, "fileNameField")
	} else {
		fileFieldName = "file"
	}

	for key, value := range options {
		uploadOptions[key] = value
	}

	for _, file := range fileList {
		req, err := server.BuildFileRequest(endpoint, uploadOptions, fileFieldName, file)

		if err != nil {
			return fmt.Errorf("error building file request: %w", err)
		}

		res, err := server.SendRequest(req, timeout)

		if err != nil {
			return fmt.Errorf("error sending file request: %w", err)
		}

		b, err := io.ReadAll(res.Body)

		if err != nil {
			return fmt.Errorf("error reading body from response: %w", err)
		}

		if res.Status != "200 OK" {
			return fmt.Errorf("%s : %s", res.Status, string(b))
		}
	}

	return nil
}
