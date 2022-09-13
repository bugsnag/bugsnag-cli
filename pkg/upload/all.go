package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"io"
)

func All(file string, uploadOptions map[string]string, uploadUrl string, timeout int) error {
	var fileFieldName string
	fileFieldName = "file"

	if uploadOptions["fileNameField"] != "" {
		fileFieldName = uploadOptions["fileNameField"]
		delete(uploadOptions, "fileNameField")
	}

	req, err := server.BuildFileRequest(uploadUrl, uploadOptions, fileFieldName, file)

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
	return nil
}
