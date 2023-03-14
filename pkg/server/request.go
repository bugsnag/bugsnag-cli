package server

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
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

		//defer file.Close()

		part, err := writer.CreateFormFile(key, filepath.Base(file.Name()))

		if err != nil {
			return nil, err
		}

		io.Copy(part, file)
	}

	for key, value := range fieldData {
		writer.WriteField(key, value)
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

// ProcessRequest - Builds and sends file requests to the API
func ProcessRequest(endpoint string, uploadOptions map[string]string, fileFieldData map[string]string, timeout int) error {
	req, err := BuildFileRequest(endpoint, uploadOptions, fileFieldData)

	if err != nil {
		return fmt.Errorf("error building file request: %w", err)
	}

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

	return nil
}
