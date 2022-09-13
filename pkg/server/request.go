package server

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// BuildFileRequest - Create a multi-part form request adding a file as a parameter
func BuildFileRequest(url string, fieldData map[string]string, fileFieldName string, fileName string) (*http.Request, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, value := range fieldData {
		writer.WriteField(key, value)
	}

	part, err := writer.CreateFormFile(fileFieldName, filepath.Base(file.Name()))

	if err != nil {
		return nil, err
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("POST", url, body)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, nil
}

// SendRequest Sends request
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
