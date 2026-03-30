package server

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/bugsnag/bugsnag-cli/pkg/endpoints"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type FileField interface {
	writeToForm(writer *multipart.Writer, key string) error
}

// requestBuilder is a function type that builds a fresh HTTP request for each attempt
type requestBuilder func() (*http.Request, error)

type LocalFile string

func (localFile LocalFile) writeToForm(writer *multipart.Writer, key string) error {
	file, err := os.Open(string(localFile))
	if err != nil {
		return err
	}

	part, err := writer.CreateFormFile(key, filepath.Base(file.Name()))
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	file.Close()
	return nil
}

type InMemoryFile struct {
	Path string
	Data []byte
}

func (inMemoryFile InMemoryFile) writeToForm(writer *multipart.Writer, key string) error {
	part, err := writer.CreateFormFile(key, inMemoryFile.Path)
	if err != nil {
		return err
	}

	_, err = part.Write(inMemoryFile.Data)
	if err != nil {
		return err
	}

	return nil
}

// buildFileRequest constructs an HTTP request for file upload with specified field data.
//
// Parameters:
//   - url: The target URL for the file upload request.
//   - fieldData: A map containing additional form fields for the request.
//   - fileFieldData: A map containing file field names and their corresponding file paths.
//
// Returns:
//   - *http.Request: The constructed HTTP request.
//   - error: An error if any step of the request construction fails.
func buildFileRequest(url string, fieldData map[string]string, fileFieldData map[string]FileField) (*http.Request, error) {
	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)

	for key, value := range fileFieldData {
		err := value.writeToForm(writer, key)
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

// ProcessFileRequest processes a file upload request by building an HTTP request,
// uploading the specified file to the endpoint, and logging information based on the dryRun flag.
// It handles the API key, constructs the endpoint URL, and manages retries in case of failures.
//
// Parameters:
//   - apiKey: The project API key.
//   - endpointPath: The path to the upload endpoint, which can be empty for the default endpoint.
//   - uploadOptions: A map containing options for building the file request.
//   - fileFieldData: A map containing data associated with the file field.
//   - fileName: The name of the file to be uploaded.
//   - options: used to determine dry run, timeout, and retries.
//
// Returns:
//   - error: An error if any step of the file processing fails. Nil if the process is successful.
func ProcessFileRequest(apiKey string, endpointPath string, uploadOptions map[string]string, fileFieldData map[string]FileField, fileName string, options options.CLI, logger log.Logger) error {

	// Check if the fileName itself should be excluded based on exclude patterns
	if len(options.Upload.Exclude) > 0 {
		if utils.IsFileExcluded(fileName, options.Upload.Exclude) {
			logger.Info(fmt.Sprintf("Skipping the upload of: %s (matches exclude pattern)", fileName))
			return nil
		}
	}

	if apiKey != "" {
		uploadOptions["apiKey"] = apiKey
	} else {
		return fmt.Errorf("missing api key, please specify using `--api-key`")
	}

	endpoint, err := endpoints.GetDefaultUploadEndpoint(apiKey, endpointPath, options)
	if err != nil {
		return fmt.Errorf("error getting upload endpoint: %w", err)
	}

	if !options.DryRun {
		logger.Info(fmt.Sprintf("Uploading %s to %s", filepath.Base(fileName), endpoint))

		// Create a builder function that constructs a fresh request for each attempt
		buildRequest := func() (*http.Request, error) {
			return buildFileRequest(endpoint, uploadOptions, fileFieldData)
		}

		err = processRequest(buildRequest, options.Upload.Timeout, options.Upload.Retries, logger)

		if err != nil {
			if strings.Contains(err.Error(), "409") {
				logger.Warn(fmt.Sprintf("Duplicate file detected, skipping upload of %s", filepath.Base(fileName)))
			} else {
				return err
			}
		} else {
			logger.Info("Uploaded " + filepath.Base(fileName))
		}
	} else {
		logger.Info(fmt.Sprintf("(dryrun) Skipping upload of %s to %s", filepath.Base(fileName), endpoint))
		logger.Debug("(dryrun) Upload payload:")
		prettyUploadOptions, _ := utils.PrettyPrintMap(uploadOptions)
		logger.Debug(prettyUploadOptions)
	}

	return nil
}

// ProcessBuildRequest processes a build request by creating an HTTP request with the provided payload,
// sending the request to the specified endpoint, and logging information based on the dryRun flag.
// It handles the API key, constructs the endpoint URL, and manages retries in case of failures.
//
// Parameters:
//   - apiKey: The project API key.
//   - payload: The payload to be sent in the request body.
//   - options: used to determine dry run, timeout, and retries.
//
// Returns:
//   - error: An error if any step of the build processing fails. Nil if the process is successful.
func ProcessBuildRequest(apiKey string, payload []byte, options options.CLI, logger log.Logger) error {
	if apiKey == "" {
		return fmt.Errorf("missing api key, please specify using `--api-key`")
	}

	endpoint, err := endpoints.GetDefaultBuildEndpoint(apiKey, options)
	if err != nil {
		return fmt.Errorf("error getting upload endpoint: %w", err)
	}

	if !options.DryRun {
		logger.Info(fmt.Sprintf("Sending build information to %s", endpoint))

		// Create a builder function that constructs a fresh request for each attempt
		buildRequest := func() (*http.Request, error) {
			req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
			if err != nil {
				return nil, err
			}
			req.Header.Add("Content-Type", "application/json")
			return req, nil
		}

		err := processRequest(buildRequest, options.Upload.Timeout, options.Upload.Retries, logger)
		if err != nil {
			return err
		}
	} else {
		logger.Info(fmt.Sprintf("(dryrun) Skipping sending build information to %s", endpoint))
		logger.Info("(dryrun) Build payload:")
		prettyUploadOptions, _ := utils.PrettyPrintJson(string(payload))
		fmt.Println(prettyUploadOptions)
	}

	return nil
}

// processRequest sends an HTTP request using sendRequest function with retry logic.
// It builds a fresh request for each attempt to avoid state pollution from previous attempts.
// It attempts to send the request multiple times, specified by retryCount parameter,
// and waits for a short duration between each attempt.
// If all attempts fail, it returns an error indicating the failure after the specified number of attempts.
// Parameters:
//   - buildRequest: A function that builds a fresh HTTP request for each attempt.
//   - timeout: Timeout duration for the HTTP request in seconds.
//   - retryCount: Number of times to retry the request in case of failure.
//
// Returns:
//   - error: An error indicating the reason for failure or nil if the request is successful.
func processRequest(buildRequest requestBuilder, timeout int, retryCount int, logger log.Logger) error {
	var err error
	i := 0
	for {
		// Build a fresh request for each attempt
		request, buildErr := buildRequest()
		if buildErr != nil {
			return errors.Wrap(buildErr, "failed to build request")
		}

		err = sendRequest(request, timeout, logger)
		if err == nil {
			return nil
		}

		i++

		if i > retryCount {
			break
		}

		logger.Warn(fmt.Sprintf("BugSnag API request attempt %d failed:", i))
		logger.Warn(err.Error())
		logger.Warn("Retrying...")

		time.Sleep(time.Second)
	}

	if err != nil {
		return errors.Errorf("failed after %d attempts. %s", i, err.Error())
	}

	return nil
}

// sendRequest sends an HTTP request using the provided request object and timeout.
//
// Parameters:
//   - request: The HTTP request to be sent.
//   - timeout: The timeout duration for the HTTP request in seconds.
//
// Returns:
//   - error: An error if any step of the request processing fails. Nil if the process is successful.
func sendRequest(request *http.Request, timeout int, logger log.Logger) error {
	// Configure transport to use HTTP/1.1 only
	var protocols http.Protocols
	protocols.SetHTTP1(true)
	protocols.SetHTTP2(false)

	transport := &http.Transport{
		Protocols: &protocols,
	}

	client := &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: transport,
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading body from response: %w", err)
	}

	contentType := response.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		warnings, err := utils.CheckResponseWarnings(responseBody)
		if err != nil {
			return err
		}

		for _, warning := range warnings {
			logger.Warn(warning.(string))
		}
	}

	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		return fmt.Errorf("%s: %s", response.Status, string(responseBody))
	}

	return nil
}
