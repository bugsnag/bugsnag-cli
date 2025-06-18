package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/unity"
	"os"
)

func ProcessUnityIos(globalOptions options.CLI, endpoint string, logger log.Logger) error {
	unityOptions := globalOptions.Upload.UnityIos
	var lineMappingFile string

	for _, path := range unityOptions.Path {
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

		if unityOptions.UnityLineMapping.NoUploadIl2cppMappingFile {
			logger.Debug("Skipping the upload of the LineNumberMappings.json file")
		} else {
			lineMappingFile, err = unity.GetiOSLineMapping(string(unityOptions.UnityLineMapping.UploadIl2cppMappingFile), path)
			if err != nil {
				return err
			}
			logger.Debug(fmt.Sprintf("Found line mapping file: %s", lineMappingFile))
		}

		for _, dsym := range dsyms {
			logger.Info(fmt.Sprintf("Processing dSYM: %s (%s)", dsym.Name, dsym.UUID))
			if dsym.Name == "UnityFramework" && !unityOptions.UnityLineMapping.NoUploadIl2cppMappingFile {
				logger.Info(fmt.Sprintf("Uploading %s for build ID %s", lineMappingFile, dsym.UUID))
				//err = unity.UploadIosLineMappings(
				//	lineMappingFile,
				//	dsym.UUID,
				//	endpoint,
				//	globalOptions,
				//	logger,
				//)
				//
				//if err != nil {
				//	return err
				//}
			}
		}

		return UploadDsyms(dsyms, plistPath, endpoint, globalOptions, logger)
	}
	return nil
}
