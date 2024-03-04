package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
)

func ProcessDysmNew(
	apiKey string,
	scheme string,
	xcodeProjectPath string,
	plistPath string,
	paths []string,
	endpoint string,
	timeout int,
	retries int,
	overwrite bool,
	dryRun bool,
) error {
	var dwarfInfo []*ios.DwarfInfo
	var tempDirs []string

	// Performs an automatic cleanup of temporary directories at the end
	defer func() {
		for _, tempDir := range tempDirs {
			_ = os.RemoveAll(tempDir)
		}
	}()

	for _, path := range paths {
		if utils.IsDir(path) {
			log.Info("directory")
			var tempDir string
			dwarfInfo, tempDir, _ = ios.FindDsymsInPath(path, false, false)
			tempDirs = append(tempDirs, tempDir)

			fmt.Println(dwarfInfo)

		} else if ios.IsDsymFile(path) {
			log.Info("file")
			// If the path is a file, we need to check if it's a dSYM
		} else {
			return fmt.Errorf("Invalid path: %s", path)
		}
	}

	return nil
}
