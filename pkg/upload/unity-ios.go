package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"os"
)

func ProcessUnityIos(globalOptions options.CLI, endpoint string, logger log.Logger) error {
	unityOptions := globalOptions.Upload.UnityIos

	globalOptions.Upload.XcodeBuild = options.XcodeBuild{
		Path: unityOptions.Path,
		Shared: options.DsymShared{
			IgnoreEmptyDsym:    unityOptions.Shared.IgnoreEmptyDsym,
			IgnoreMissingDwarf: unityOptions.Shared.IgnoreMissingDwarf,
			Configuration:      unityOptions.Shared.Configuration,
			Scheme:             unityOptions.Shared.Scheme,
			ProjectRoot:        unityOptions.Shared.ProjectRoot,
			Plist:              unityOptions.Shared.Plist,
			XcodeProject:       unityOptions.Shared.XcodeProject,
		},
	}

	dsyms, plistPath, tempDirs, err := FindDsymsAndSettings(globalOptions, logger)
	defer func() {
		for _, tempDir := range tempDirs {
			_ = os.RemoveAll(tempDir)
		}
	}()
	if err != nil {
		return err
	}

	for _, dsym := range dsyms {
		dsymInfo := fmt.Sprintf("(UUID: %s, Name: %s, Arch: %s)", dsym.UUID, dsym.Name, dsym.Arch)
		logger.Debug(fmt.Sprintf("Processing dSYM %s", dsymInfo))
	}

	if plistPath == "" {
		logger.Warn("No Info.plist found, using default settings")
	}

	return nil
}
