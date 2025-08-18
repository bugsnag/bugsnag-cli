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

func ProcessUnityIos(globalOptions options.CLI, logger log.Logger) error {
	var (
		err                  error
		possibleDsymPath     string
		possibleXcodeProject string
		dsyms                []*ios.DwarfInfo
		tempDir              string
		tempDirs             []string
		lineMappingFile      string
		plistPath            string
		plistData            *ios.PlistData
	)

	unityOptions := globalOptions.Upload.UnityIos

	defer func() {
		for _, dir := range tempDirs {
			_ = os.RemoveAll(dir)
		}
	}()

	for _, path := range unityOptions.Path {
		if unityOptions.DsymShared.Scheme == "" {
			unityOptions.DsymShared.Scheme = "Unity-iPhone"
			logger.Debug(fmt.Sprintf("Using default Unity scheme: %s", unityOptions.DsymShared.Scheme))
		}

		if unityOptions.DsymShared.ProjectRoot == "" {
			if utils.IsDir(path) {
				unityOptions.DsymShared.ProjectRoot = path
			} else {
				unityOptions.DsymShared.ProjectRoot = filepath.Dir(path)
			}
		}

		logger.Debug(fmt.Sprintf("Using %s as the project root", unityOptions.DsymShared.ProjectRoot))

		if unityOptions.DsymPath == "" {
			if unityOptions.DsymShared.XcodeProject == "" {
				if ios.IsPathAnXcodeProjectOrWorkspace(path) {
					possibleXcodeProject = path
				} else {
					possibleXcodeProject = ios.FindXcodeProjOrWorkspace(path)
				}
			} else {
				possibleXcodeProject = ios.FindXcodeProjOrWorkspace(string(unityOptions.DsymShared.XcodeProject))
			}

			if possibleXcodeProject == "" {
				possibleDsymPath = path
			} else {
				globalOptions.Upload.XcodeArchive = options.XcodeArchive{
					Path:   utils.Paths{possibleXcodeProject},
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
				possibleDsymPath = xcarchivePath
			}

			logger.Debug(fmt.Sprintf("Searching for dSYMs in: %s", possibleDsymPath))
		} else {
			possibleDsymPath = string(unityOptions.DsymPath)
			logger.Debug(fmt.Sprintf("Using dSYM path: %s", possibleDsymPath))
		}

		dsyms, tempDir, err = ios.FindDsymsInPath(
			possibleDsymPath,
			unityOptions.DsymShared.IgnoreEmptyDsym,
			unityOptions.DsymShared.IgnoreMissingDwarf,
			logger,
		)

		tempDirs = append(tempDirs, tempDir)

		if err != nil {
			return fmt.Errorf("error locating dSYMs in %s: %w", possibleDsymPath, err)
		}

		if len(dsyms) == 0 {
			return fmt.Errorf("no dSYMs found in %s", possibleDsymPath)
		}

		logger.Info(fmt.Sprintf("Found %d dSYM files in: %s", len(dsyms), possibleDsymPath))

		if globalOptions.Upload.UnityIos.VersionName == "" || globalOptions.Upload.UnityIos.BundleVersion == "" || globalOptions.Upload.UnityIos.ApplicationId == "" {
			plistPath := string(unityOptions.DsymShared.Plist)
			if plistPath == "" {
				plistPath = filepath.Join(dsyms[0].Location, "..", "..", "Info.plist")
			}

			plistData, _ = ios.GetPlistData(plistPath)

			if plistData != nil {
				logger.Debug(fmt.Sprintf("Reading plist data from: %s", plistPath))

				if globalOptions.Upload.UnityIos.VersionName == "" {
					globalOptions.Upload.UnityIos.VersionName = plistData.VersionName
				}

				if globalOptions.Upload.UnityIos.BundleVersion == "" {
					globalOptions.Upload.UnityIos.BundleVersion = plistData.BundleVersion
				}

				if globalOptions.Upload.UnityIos.ApplicationId == "" {
					globalOptions.Upload.UnityIos.ApplicationId = plistData.BundleIdentifier
				}
			} else {
				logger.Debug("No plist file found")
			}
		}

		if unityOptions.UnityShared.NoUploadIl2cppMapping {
			logger.Debug("Skipping the upload of the LineNumberMappings.json file")
		} else if unityOptions.UnityShared.UploadIl2cppMapping != "" {
			lineMappingFile = string(unityOptions.UnityShared.UploadIl2cppMapping)
		} else {
			lineMappingFile, err = unity.GetIosLineMapping(path)
			if err != nil {
				return fmt.Errorf("failed to get line mapping file: %w", err)
			}
			logger.Info(fmt.Sprintf("Found line mapping file: %s", lineMappingFile))
		}

		for _, dsym := range dsyms {
			if dsym.Name == "UnityFramework" && lineMappingFile != "" {
				if dsym.UUID == "" {
					return fmt.Errorf("dSYM %s has no UUID, cannot upload line mappings", dsym.Name)
				}
			}

			err := ios.ProcessDsymUpload(plistPath, unityOptions.DsymShared.ProjectRoot, globalOptions, []*ios.DwarfInfo{dsym}, logger)

			if err != nil {
				return fmt.Errorf("Error uploading dSYM files: %w", err)
			}

			if dsym.Name == "UnityFramework" && lineMappingFile != "" {
				logger.Info(fmt.Sprintf("Uploading %s for dSYM %s, withID %s", lineMappingFile, dsym.Name, dsym.UUID))

				err = unity.UploadUnityLineMappings(
					globalOptions.ApiKey,
					"ios",
					dsym.UUID,
					globalOptions.Upload.UnityIos.ApplicationId,
					globalOptions.Upload.UnityIos.VersionName,
					globalOptions.Upload.UnityIos.BundleVersion,
					lineMappingFile,
					unityOptions.DsymShared.ProjectRoot,
					unityOptions.Overwrite,
					globalOptions,
					logger,
				)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
