package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/breakpad"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"path/filepath"
)

func ProcessBreakpad(globalOptions options.CLI, endpoint string, logger log.Logger) error {
	breakpadOptions := globalOptions.Upload.Breakpad
	var symFileList []string

	for _, path := range breakpadOptions.SymFilePath {
		if utils.IsDir(path) {
			files, err := utils.BuildFileList([]string{path})
			if err != nil {
				return err
			}
			for _, file := range files {
				if filepath.Ext(file) == ".sym" {
					symFileList = append(symFileList, file)
				}
			}
		} else if filepath.Ext(path) == ".sym" {
			symFileList = append(symFileList, path)
		} else {
			logger.Warn(fmt.Sprintf("Skipping %s (not a .sym file or directory)", path))
		}
	}

	if len(symFileList) == 0 {
		logger.Info("No Breakpad .sym files found, skipping upload")
		return nil
	}

	logger.Info(fmt.Sprintf("Uploading %d Breakpad .sym files to Bugsnag", len(symFileList)))

	err := breakpad.UploadBreakpad(symFileList, globalOptions.ApiKey, breakpadOptions.ProjectRoot, endpoint, globalOptions, logger)
	if err != nil {
		return err
	}
	logger.Info("Breakpad symbol upload complete")
	return nil
}
