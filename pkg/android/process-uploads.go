package android

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
)

// buildUploadOptions constructs the form parameters required to upload an NDK symbol file.
//
// Parameters:
//   - appID: Android application ID (e.g., com.example.app)
//   - versionCode: Numeric build version (e.g., "42")
//   - versionName: Human-readable version name (e.g., "1.2.3")
//   - projectRoot: Root path of the project
//   - fileName: Original .so file name
//   - opts: CLI options containing upload flags
//
// Returns:
//   - A map[string]string of form values to be sent in the upload request.
func buildUploadOptions(appID, versionCode, versionName, projectRoot, fileName string, overwrite bool) map[string]string {
	uploadOpts := map[string]string{}

	if appID != "" {
		uploadOpts["appId"] = appID
	}
	if versionCode != "" {
		uploadOpts["versionCode"] = versionCode
	}
	if versionName != "" {
		uploadOpts["versionName"] = versionName
	}
	if projectRoot != "" {
		uploadOpts["projectRoot"] = projectRoot
	}
	if base := filepath.Base(fileName); base != "" {
		uploadOpts["sharedObjectName"] = base
	}
	if overwrite {
		uploadOpts["overwrite"] = "true"
	}

	return uploadOpts
}

// UploadAndroidNdk uploads extracted Android NDK symbol files (.so debug symbols)
// to the Bugsnag /ndk-symbol endpoint.
//
// Each file is uploaded with associated metadata including app ID, version code/name,
// and shared object name. If no symbol files are present, the upload is skipped.
//
// Parameters:
//   - symbolFiles: Map of original file path → generated symbol file path
//   - apiKey: Bugsnag project API key
//   - appID: Android app package name
//   - versionName: App version name (semantic)
//   - versionCode: App version code (build number)
//   - projectRoot: Absolute or relative root of the Android project
//   - opts: CLI options, including `--overwrite`
//   - logger: Logger instance for debug/info/error output
//
// A failure to upload one file does not stop the remaining files from being
// uploaded: the server may reject individual libraries (e.g. third-party
// libraries in an unsupported format), and that should not prevent the app's
// own symbols from being uploaded. Files are uploaded in a stable, sorted order.
//
// Returns:
//   - error: Non-nil if any file fails to upload; combines the errors of all failed files
func UploadAndroidNdk(
	symbolFiles map[string]string,
	apiKey string,
	appID string,
	versionName string,
	versionCode string,
	projectRoot string,
	opts options.CLI,
	overwrite bool,
	logger log.Logger,
) error {
	if len(symbolFiles) == 0 {
		logger.Info("No NDK files found to process")
		return nil
	}

	originalFiles := make([]string, 0, len(symbolFiles))
	for originalFile := range symbolFiles {
		originalFiles = append(originalFiles, originalFile)
	}
	sort.Strings(originalFiles)

	var uploadErrors []error

	for _, originalFile := range originalFiles {
		symbolPath := symbolFiles[originalFile]
		fileField := map[string]server.FileField{
			"soFile": server.LocalFile(symbolPath),
		}

		params := buildUploadOptions(appID, versionCode, versionName, projectRoot, originalFile, overwrite)

		err := server.ProcessFileRequest(
			apiKey,
			"/ndk-symbol",
			params,
			fileField,
			originalFile,
			opts,
			logger,
		)
		if err != nil {
			uploadErr := fmt.Errorf("failed to upload NDK symbol for %s: %w", originalFile, err)
			logger.Error(uploadErr.Error())
			uploadErrors = append(uploadErrors, uploadErr)
		}
	}

	return errors.Join(uploadErrors...)
}
