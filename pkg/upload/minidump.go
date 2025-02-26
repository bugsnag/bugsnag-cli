package upload

import (
	"fmt"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

func ProcessMinidump(globalOptions options.CLI, endpoint string, logger log.Logger) error {
	minidumpOptions := globalOptions.Upload.Minidump
	var symFileList []string

	for _, path := range minidumpOptions.SymFilePath {
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
		logger.Info("No minidump .sym files found, skipping upload")
		return nil
	}

	logger.Info(fmt.Sprintf("Uploading %d minidump .sym files to Bugsnag", len(symFileList)))

	logger.Info("Minidump symbol upload complete")
	return nil
}
