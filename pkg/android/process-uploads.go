package android

import (
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// UploadAndroidNdk uploads a list of Android NDK symbol files to Bugsnag.
//
// This function processes each .so (shared object) file by preparing the necessary
// metadata and sending a request to the NDK symbol upload endpoint. If no files
// are provided, the function exits early.
//
// Parameters:
//   - fileList: List of paths to .so files to upload.
//   - apiKey: The Bugsnag project API key.
//   - applicationId: The application ID for the Android app.
//   - versionName: The version name of the build.
//   - versionCode: The version code of the build.
//   - projectRoot: The root path of the project.
//   - options: CLI options, including overwrite flag.
//   - logger: Logger for structured output and progress information.
//
// Returns:
//   - error: Non-nil if any step of the upload fails.
func UploadAndroidNdk(
	fileList []string,
	apiKey string,
	applicationId string,
	versionName string,
	versionCode string,
	projectRoot string,
	options options.CLI,
	logger log.Logger,
) error {
	fileFieldData := make(map[string]server.FileField)
	numberOfFiles := len(fileList)

	// No files to upload; exit early
	if numberOfFiles < 1 {
		logger.Info("No NDK files found to process")
		return nil
	}

	// Iterate over each provided .so file
	for _, file := range fileList {
		// Construct upload options using project metadata
		uploadOptions, err := utils.BuildAndroidNDKUploadOptions(
			applicationId, versionName, versionCode, projectRoot, filepath.Base(file), options.Upload.Overwrite,
		)
		if err != nil {
			return err
		}

		// Attach file to form field
		fileFieldData["soFile"] = server.LocalFile(file)

		// Upload the file to the /ndk-symbol endpoint
		err = server.ProcessFileRequest(apiKey, "/ndk-symbol", uploadOptions, fileFieldData, file, options, logger)
		if err != nil {
			return err
		}
	}

	return nil
}
