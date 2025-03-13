package upload

import (
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
)

func ProcessDsym(globalOptions options.CLI, endpoint string, logger log.Logger) error {
	dsymOptions := globalOptions.Upload.Dsym

	logger.Info("Searching for dSYM files to upload")

	globalOptions.Upload.XcodeArchive = options.XcodeArchive{
		Path:   dsymOptions.Path,
		Shared: dsymOptions.Shared,
	}

	// Attempt to process dSYMs from the Xcode archive
	if err := ProcessXcodeArchive(globalOptions, endpoint, logger); err != nil {
		if err.Error() == "No xcarchive found in specified paths" {
			logger.Info("No dSYM files found in the xcarchive, searching in the build directory")
		} else {
			return err
		}
	} else {
		logger.Info("dSYM files successfully uploaded from the xcarchive")
		return nil
	}

	// Attempt to process dSYMs from the Xcode build directory
	globalOptions.Upload.XcodeBuild = options.XcodeBuild{
		Path:   dsymOptions.Path,
		Shared: dsymOptions.Shared,
	}

	if err := ProcessXcodeBuild(globalOptions, endpoint, logger); err != nil {
		return err
	}

	logger.Info("dSYM files successfully uploaded from the build directory")
	return nil
}
