package upload

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/unity"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

func ProcessUnityIos(globalOptions options.CLI, endpoint string, logger log.Logger) error {
	unityOptions := globalOptions.Upload.UnityIos
	var tempDirs []string

	defer func() {
		for _, dir := range tempDirs {
			_ = os.RemoveAll(dir)
		}
	}()

	for _, path := range unityOptions.Path {
		var (
			err             error
			pathToCheck     string
			dsyms           []*ios.DwarfInfo
			tempDir         string
			lineMappingFile string
			plistData       *ios.PlistData
		)

		if unityOptions.DsymShared.Scheme == "" {
			unityOptions.DsymShared.Scheme = "Unity-iPhone"
			logger.Info(fmt.Sprintf("Using scheme: %s", unityOptions.DsymShared.Scheme))
		}

		// Handle Xcode project detection
		if ios.IsPathAnXcodeProjectOrWorkspace(path) {
			globalOptions.Upload.XcodeArchive = options.XcodeArchive{
				Path:   utils.Paths{path},
				Shared: unityOptions.DsymShared,
			}

			xcarchivePath, err := ios.FindXcarchivePath(globalOptions, logger)
			if err != nil {
				return fmt.Errorf("failed to find Xcode archive path: %w", err)
			}
			if xcarchivePath == "" {
				return fmt.Errorf("no Xcode archive found in specified paths")
			}

			logger.Info(fmt.Sprintf("Found Xcode archive at %s", xcarchivePath))
			pathToCheck = xcarchivePath
		} else {
			pathToCheck = path
		}

		dsyms, tempDir, err = ios.FindDsymsInPath(
			pathToCheck,
			unityOptions.DsymShared.IgnoreEmptyDsym,
			unityOptions.DsymShared.IgnoreMissingDwarf,
			logger,
		)
		tempDirs = append(tempDirs, tempDir)

		if err != nil {
			return fmt.Errorf("error locating dSYMs in %s: %w", pathToCheck, err)
		}
		if len(dsyms) == 0 {
			return fmt.Errorf("no dSYMs found in %s", pathToCheck)
		}

		logger.Info(fmt.Sprintf("Found %d dSYM files in: %s", len(dsyms), pathToCheck))

		// Resolve plist path if not explicitly provided
		plistPath := string(unityOptions.DsymShared.Plist)
		if plistPath == "" {
			plistPath = filepath.Join(dsyms[0].Location, "..", "..", "Info.plist")
		}

		if !utils.FileExists(plistPath) {
			return fmt.Errorf("plist file not found at: %s", plistPath)
		}
		logger.Info(fmt.Sprintf("Using plist path: %s", plistPath))

		plistData, err = ios.GetPlistData(plistPath)
		if err != nil {
			return fmt.Errorf("failed to read plist: %w", err)
		}

		if plistData != nil {
			if globalOptions.Upload.UnityIos.VersionName == "" {
				globalOptions.Upload.UnityIos.VersionName = plistData.VersionName
			}

			if globalOptions.Upload.UnityIos.BundleVersion == "" {
				globalOptions.Upload.UnityIos.BundleVersion = plistData.BundleVersion
			}

			if globalOptions.Upload.UnityIos.ApplicationId == "" {
				globalOptions.Upload.UnityIos.ApplicationId = plistData.BundleIdentifier
			}
		}

		if unityOptions.UnityShared.NoUploadIl2cppMappingFile {
			logger.Debug("Skipping the upload of the LineNumberMappings.json file")
		} else if unityOptions.UnityShared.UploadIl2cppMappingFile != "" {
			lineMappingFile = string(unityOptions.UnityShared.UploadIl2cppMappingFile)
		} else {
			lineMappingFile, err = unity.GetIosLineMapping(path)
			if err != nil {
				return fmt.Errorf("failed to get line mapping file: %w", err)
			}
			logger.Info(fmt.Sprintf("Found line mapping file: %s", lineMappingFile))
		}

		// Log dSYM details
		for _, dsym := range dsyms {
			if dsym.Name == "UnityFramework" && lineMappingFile != "" {
				logger.Info(fmt.Sprintf("Uploading %s for dSYM %s, withID %s", lineMappingFile, dsym.Name, dsym.UUID))
				err = unity.UploadIosLineMappings(
					lineMappingFile,
					dsym.UUID,
					endpoint,
					globalOptions,
					logger,
				)

				if err != nil {
					return err
				}
			}
		}

		// Upload dSYMs and plist
		err = UploadDsyms(dsyms, plistPath, endpoint, globalOptions, logger)
		if err != nil {
			return fmt.Errorf("failed to upload dSYMs: %w", err)
		}
	}

	return nil
}
