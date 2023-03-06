package upload

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type ReactNativeAndroid struct {
	CodeBundleId  string `help:"A unique identifier to identify a code bundle release when using tools like CodePush"`
	Dev           bool   `help:"Indicates whether the application is a debug or release build"`
	SourceMapPath string `help:"Path to the source map file" type:"path"`
	BundlePath    string `help:"Path to the bundle file" type:"path"`
	ProjectRoot   string `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
}

func ProcessReactNativeAndroid(appVersion string, appVersionCode string, codeBundleId string, dev bool, sourceMapPath string, bundlePath string, projectRoot string, endpoint string, timeout int, retries int, overwrite bool, apiKey string) error {

	if appVersion == "" {
		return fmt.Errorf("`--app-version` missing from options")
	}

	if sourceMapPath == "" {
		return fmt.Errorf("`--source-map-path` missing from options")
	}

	if bundlePath == "" {
		return fmt.Errorf("`--bundle-path` missing from options")
	}

	// Check if we have project root
	if projectRoot == "" {
		return fmt.Errorf("`--project-root` missing from options")
	}

	// Check if SourceMapPath exists
	if !utils.FileExists(sourceMapPath) {
		return fmt.Errorf(sourceMapPath + " does not exist on the system.")
	}

	// Check if bundlePath exists
	if !utils.FileExists(bundlePath) {
		return fmt.Errorf(bundlePath + " does not exist on the system.")
	}

	log.Info("Uploading debug information for React Native Android")

	uploadOptions := utils.BuildReactNativeAndroidUploadOptions(apiKey, appVersion, appVersionCode, codeBundleId, dev, projectRoot, overwrite)

	req, err := BuildFileRequest(endpoint, uploadOptions, sourceMapPath, bundlePath)

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

	statusOK := res.StatusCode >= 200 && res.StatusCode < 300

	if !statusOK {
		return fmt.Errorf("%s : %s", res.Status, string(b))
	}

	return nil
}

func BuildFileRequest(url string, fieldData map[string]string, sourceMapPath string, bundlePath string) (*http.Request, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for key, value := range fieldData {
		writer.WriteField(key, value)
	}

	sourceMapFile, err := os.Open(sourceMapPath)

	if err != nil {
		return nil, fmt.Errorf("Unable to open " + sourceMapPath + "\n " + err.Error())
	}

	defer sourceMapFile.Close()

	sourceMapFilePart, err := writer.CreateFormFile("sourceMap", filepath.Base(sourceMapFile.Name()))

	if err != nil {
		return nil, err
	}

	io.Copy(sourceMapFilePart, sourceMapFile)

	bundleFile, err := os.Open(bundlePath)

	if err != nil {
		return nil, fmt.Errorf("Unable to open " + bundlePath + "\n " + err.Error())
	}

	defer bundleFile.Close()

	bundleFilePart, err := writer.CreateFormFile("bundle", filepath.Base(bundleFile.Name()))

	if err != nil {
		return nil, err
	}

	io.Copy(bundleFilePart, sourceMapFile)

	writer.Close()

	request, err := http.NewRequest("POST", url, body)

	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

	return request, nil
}
