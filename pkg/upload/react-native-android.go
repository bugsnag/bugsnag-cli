package upload

import (
	"fmt"
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

	fileFieldData := make(map[string]string)
	fileFieldData["sourceMap"] = sourceMapPath
	fileFieldData["bundle"] = bundlePath

	requestStatus := server.ProcessRequest(endpoint, uploadOptions, fileFieldData, timeout)

	if requestStatus != nil {
		return requestStatus

	}

	return nil
}
