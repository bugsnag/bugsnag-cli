package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"strings"
)

func ProcessBreakpad(globalOptions options.CLI, endpoint string, logger log.Logger) error {
	breakpadOptions := globalOptions.Upload.Breakpad
	apiKey := globalOptions.ApiKey
	projectRoot := globalOptions.Upload.Breakpad.ProjectRoot

	symFileList, err := utils.BuildFileList(breakpadOptions.Path)
	if err != nil {
		return err
	}

	if len(symFileList) == 0 {
		logger.Error("No .sym files found")
		return nil
	}

	logger.Debug(fmt.Sprintf("Uploading %d .sym files", len(symFileList)))

	return UploadBreakpadSymbols(symFileList, apiKey, projectRoot, endpoint, globalOptions, logger)
}

func UploadBreakpadSymbols(
	fileList []string,
	apiKey string,
	projectRoot string,
	endpoint string,
	options options.CLI,
	logger log.Logger,
) error {
	if apiKey == "" {
		return fmt.Errorf("missing api key, please specify using --api-key")
	}

	breakpadOptions := options.Upload.Breakpad

	for _, file := range fileList {
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

		fileFieldData := map[string]server.FileField{
			"symbol_file": server.LocalFile(file),
		}

		queryParams := fmt.Sprintf("?api_key=%s&overwrite=%t&project_root=%s",
			strings.ReplaceAll(apiKey, " ", "%20"),
			options.Upload.Overwrite,
			strings.ReplaceAll(projectRoot, " ", "%20"),
		)

		err = server.ProcessFileRequest(endpoint+"/breakpad-symbol"+queryParams, formFields, fileFieldData, file, options, logger)
		if err != nil {
			return err
		}
	}
	return nil
}
