package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/elf"
	"github.com/bugsnag/bugsnag-cli/pkg/unity"
	"os"
	"path/filepath"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

func ProcessUnityAndroid(globalOptions options.CLI, endpoint string, logger log.Logger) error {
	var (
		zipPath         string
		archList        []string
		symbolFileList  []string
		manifestData    map[string]string
		lineMappingFile string
		buildDirectory  string
		aabPath         string
		fileList        []string
	)

	unityOptions := globalOptions.Upload.UnityAndroid

	for _, path := range unityOptions.Path {
		aabPath = string(unityOptions.AabPath)

		if utils.IsDir(path) {
			buildDirectory = path
			zipPath, _ = utils.FindLatestFileWithSuffix(path, ".symbols.zip")

			if aabPath == "" {
				aabPath, _ = utils.FindLatestFileWithSuffix(path, ".aab")
			}
		} else if strings.HasSuffix(path, ".symbols.zip") {
			zipPath = path
			if aabPath == "" {
				buildDirectory = filepath.Dir(path)
				aabPath, _ = utils.FindLatestFileWithSuffix(buildDirectory, ".aab")
			}
		} else {
			return fmt.Errorf("%s is not a .symbols.zip file or containing directory", path)
		}
	}

	if aabPath != "" {
		logger.Debug(fmt.Sprintf("Extracting %s into a temporary directory", filepath.Base(aabPath)))

		aabDir, err := utils.ExtractFile(aabPath, "aab")

		if err != nil {
			return err
		}

		defer os.RemoveAll(aabDir)

		manifestData, err = android.MergeUploadOptionsFromAabManifest(aabDir, globalOptions.ApiKey, unityOptions.ApplicationId, unityOptions.BuildUuid, unityOptions.NoBuildUuid, unityOptions.VersionCode, unityOptions.VersionName, logger)

		if err != nil {
			return err
		}

		globalOptions.ApiKey = manifestData["apiKey"]
		globalOptions.Upload.AndroidAab = options.AndroidAabMapping{
			ApplicationId: manifestData["applicationId"],
			BuildUuid:     manifestData["buildUuid"],
			NoBuildUuid:   unityOptions.NoBuildUuid,
			Path:          []string{aabDir},
			ProjectRoot:   unityOptions.ProjectRoot,
			VersionCode:   manifestData["versionCode"],
			VersionName:   manifestData["versionName"],
		}
		err = ProcessAndroidAab(globalOptions, endpoint, logger)

		if err != nil {
			return err
		}
	}

	if zipPath == "" {
		logger.Info("No Unity Android symbols.zip file found, skipping")
	} else {
		logger.Debug(fmt.Sprintf("Extracting %s into a temporary directory", filepath.Base(zipPath)))

		if manifestData == nil {
			manifestData, _ = android.MergeUploadOptionsFromAabManifest("", globalOptions.ApiKey, unityOptions.ApplicationId, unityOptions.BuildUuid, unityOptions.NoBuildUuid, unityOptions.VersionCode, unityOptions.VersionName, logger)
		}

		unityDir, err := utils.ExtractFile(zipPath, "unity-android")

		if err != nil {
			return err
		}

		defer os.RemoveAll(unityDir)

		archList, err = utils.BuildDirectoryList([]string{unityDir})

		if err != nil {
			return err
		}

		if unityOptions.UnityShared.NoUploadIl2cppMapping {
			logger.Debug("Skipping the upload of the LineNumberMappings.json file")
		} else if unityOptions.UnityShared.UploadIl2cppMapping != "" {
			lineMappingFile = string(unityOptions.UnityShared.UploadIl2cppMapping)
		} else {
			lineMappingFile, err = unity.GetAndroidLineMapping(buildDirectory)
			if err != nil {
				return err
			}
			logger.Debug(fmt.Sprintf("Found line mapping file: %s", lineMappingFile))
		}

		for _, arch := range archList {
			soPath := filepath.Join(unityDir, arch)
			fileList, err = utils.BuildFileList([]string{soPath})
			if err != nil {
				return err
			}
			for _, file := range fileList {
				if filepath.Base(file) == "libil2cpp.sym.so" && utils.ContainsString(fileList, "libil2cpp.dbg.so") {
					continue
				}
				if filepath.Base(file) == "libil2cpp.so" && !unityOptions.UnityShared.NoUploadIl2cppMapping {
					_, err := elf.GetBuildId(file)
					if err != nil {
						return fmt.Errorf("failed to get build ID from %s: %w", file, err)
					}
				}
				symbolFileList = append(symbolFileList, file)
			}
		}

		for _, file := range symbolFileList {
			err = android.UploadAndroidNdk(
				file,
				manifestData["apiKey"],
				manifestData["applicationId"],
				manifestData["versionName"],
				manifestData["versionCode"],
				unityOptions.ProjectRoot,
				endpoint,
				globalOptions,
				unityOptions.Overwrite,
				logger,
			)

			if err != nil {
				return err
			}

			if filepath.Base(file) == "libil2cpp.so" && !unityOptions.UnityShared.NoUploadIl2cppMapping {
				buildId, _ := elf.GetBuildId(file)
				logger.Info(fmt.Sprintf("Uploading %s for build ID %s", lineMappingFile, buildId))
				err = unity.UploadUnityLineMappings(
					manifestData["apiKey"],
					"android",
					buildId,
					manifestData["applicationId"],
					manifestData["versionName"],
					manifestData["versionCode"],
					lineMappingFile,
					unityOptions.ProjectRoot,
					unityOptions.Overwrite,
					endpoint,
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
