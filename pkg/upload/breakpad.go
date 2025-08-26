package upload

import (
	"fmt"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// ProcessBreakpad uploads Breakpad symbol files (.sym) to Bugsnag.
//
// It validates required options, builds a list of symbol files,
// prepares upload parameters, and sends files to the server.
//
// Parameters:
// - globalOptions: CLI options including upload configuration and API key.
// - logger: logger instance for outputting progress and errors.
//
// Returns:
// - error: if any step fails during processing or uploading.
func ProcessBreakpad(globalOptions options.CLI, logger log.Logger) error {
	breakpadOptions := globalOptions.Upload.Breakpad
	apiKey := globalOptions.ApiKey
	projectRoot := globalOptions.Upload.Breakpad.ProjectRoot

	if apiKey == "" {
		return fmt.Errorf("missing api key, please specify using --api-key")
	}

	// Collect all .sym files from given paths
	symFileList, err := utils.BuildFileList(breakpadOptions.Path)
	if err != nil {
		return err
	}

	if len(symFileList) == 0 {
		logger.Error("No .sym files found")
		return nil
	}

	logger.Debug(fmt.Sprintf("Uploading %d .sym files", len(symFileList)))

	for _, file := range symFileList {
		// Build form fields for the upload
		formFields, err := utils.BuildBreakpadUploadOptions(
			breakpadOptions.CpuArch,
			breakpadOptions.CodeFile,
			breakpadOptions.DebugFile,
			breakpadOptions.DebugIdentifier,
			breakpadOptions.ProductName,
			breakpadOptions.OsName,
			breakpadOptions.VersionName,
		)
		if err != nil {
			return err
		}

		// Prepare the file to be uploaded
		fileFieldData := map[string]server.FileField{
			"symbol_file": server.LocalFile(file),
		}

		// Build query parameters for the request
		queryParams := fmt.Sprintf("?api_key=%s&overwrite=%t&project_root=%s",
			strings.ReplaceAll(apiKey, " ", "%20"),
			breakpadOptions.Overwrite,
			strings.ReplaceAll(projectRoot, " ", "%20"),
		)

		// Send the file upload request to the Breakpad symbol endpoint
		err = server.ProcessFileRequest(apiKey, "/breakpad-symbol"+queryParams, formFields, fileFieldData, file, globalOptions, logger)
		if err != nil {
			return err
		}
	}

	return nil
}
