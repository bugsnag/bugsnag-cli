package server

import (
	"bytes"
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
)

// BuildFileRequest - Create a multi-part form request adding a file as a parameter
func BuildFileRequest(url string, fieldData map[string]string, fileFieldData map[string]string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, value := range fileFieldData {
		file, err := os.Open(value)

		if err != nil {
			return nil, err
		}

		part, err := writer.CreateFormFile(key, filepath.Base(file.Name()))

		if err != nil {
			return nil, err
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return nil, err
		}
	}

	for key, value := range fieldData {
		err := writer.WriteField(key, value)
		if err != nil {
			return nil, err
		}
	}

	writer.Close()

	request, err := http.NewRequest("POST", url, body)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

	return request, nil
}

// SendRequest - Sends request
func SendRequest(request *http.Request, timeout int) (*http.Response, error) {

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ProcessFileRequest performs the handling of a file upload request to a specified endpoint.
// It builds an HTTP request using the provided options and file field data, then sends the request.
//
// Parameters:
//   - endpoint: The target URL for the file upload.
//   - uploadOptions: A map containing options for building the file request.
//   - fileFieldData: A map containing data associated with the file field.
//   - timeout: The maximum time allowed for the HTTP request.
//   - fileName: The name of the file to be uploaded.
//   - dryRun: If true, the function performs a dry run without actually sending the file.
//
// Returns:
//   - error: An error if any step of the file processing fails. Nil if the process is successful.
func ProcessFileRequest(endpoint string, uploadOptions map[string]string, fileFieldData map[string]string, timeout int, fileName string, dryRun bool) error {
	req, err := BuildFileRequest(endpoint, uploadOptions, fileFieldData)
	if err != nil {
		return fmt.Errorf("error building file request: %w", err)
	}

	if !dryRun {
		log.Info("Uploading " + filepath.Base(fileName) + " to " + endpoint)

		res, err := SendRequest(req, timeout)
		if err != nil {
			return fmt.Errorf("error sending file request: %w", err)
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading body from response: %w", err)
		}

		statusOK := res.StatusCode >= 200 && res.StatusCode < 300
		if !statusOK {
			return fmt.Errorf("%s : %s", res.Status, string(b))
		}

	} else {
		log.Info("(dryrun) Skipping upload of " + filepath.Base(fileName) + " to " + endpoint)
	}

	return nil
}

// ProcessRequest sends an HTTP POST request to a specified endpoint with the given payload.
// It allows for a dry run mode, where the request is not actually sent but is logged instead.
//
// Parameters:
//   - endpoint: The target URL for the HTTP POST request.
//   - payload: The payload to be sent in the request body.
//   - timeout: The maximum time allowed for the HTTP request.
//   - dryRun: If true, the function performs a dry run without actually sending the request.
//
// Returns:
//   - error: An error if any step of the request processing fails. Nil if the process is successful.
func ProcessRequest(endpoint string, payload []byte, timeout int, dryRun bool) error {
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")

	if !dryRun {
		res, err := SendRequest(req, timeout)
		if err != nil {
			return fmt.Errorf("error sending file request: %w", err)
		}

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("error reading body from response: %w", err)
		}

		warnings, err := utils.CheckResponseWarnings(responseBody)
		if err != nil {
			return err
		}

		for _, warning := range warnings {
			log.Info(warning.(string))
		}

		if res.StatusCode != 200 {
			return fmt.Errorf("%s : %s", res.Status, string(responseBody))
		}

	} else {
		log.Info("(dryrun) Skipping sending build information to " + endpoint)
	}

	return nil
}
