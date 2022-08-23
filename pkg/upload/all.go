package upload

import (
	"errors"
	"io"
	"github.com/bugsnag/bugsnag-cli/pkg/server"

)

func All(file string, uploadOptions map[string]string, uploadUrl string, timeout int) (string, error) {
	var fileFieldName string
	fileFieldName = "file"

	if uploadOptions["fileNameField"] != "" {
		fileFieldName = uploadOptions["fileNameField"]
		delete(uploadOptions, "fileNameField")
	}

	req, err := server.BuildFileRequest(uploadUrl, uploadOptions, fileFieldName, file)

	if err != nil {
		return "error building file request", err
	}

	res, err := server.SendRequest(req, timeout)

	if err != nil {
		return "error sending file request", err
	}

	b, err := io.ReadAll(res.Body)

	if err != nil {
		return "error reading body from response", err
	}

	if res.Status != "200 OK" {
		err := errors.New(res.Status)
		return res.Status + " " + string(b), err
	}
	return "OK", nil
}
