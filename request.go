package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Create a multi-part form request adding a file as a parameter
func BuildFileRequest(url string, fieldName string, filename string) (*http.Request, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(fieldName, filepath.Base(file.Name()))

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

func SendRequest(request *http.Request) (*http.Response, error) {
	client := &http.Client{}

    response, err := client.Do(request)
    if err != nil {
        return nil, err
    }

	return response, nil
}
