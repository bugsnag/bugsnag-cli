package upload

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

func ProcessAndroidAab(globalOptions options.CLI, endpoint string, logger log.Logger) error {

	var manifestData map[string]string
	var aabDir string
	var aabFile string
	var err error
	aabOptions := globalOptions.Upload.AndroidAab

	for _, path := range aabOptions.Path {
		// Look for AAB file if the upload command was run somewhere within the project root
		// based on an expected path of ${dir}/build/outputs/bundle/release/${dir}-release.aab
		// or ${dir}/build/outputs/bundle/release/${dir}-release-dexguard.aab
		if utils.IsDir(path) {
			if utils.FileExists(filepath.Join(path, "BUNDLE-METADATA")) {
				aabDir = path
			} else {
				arr := []string{"*", "build", "outputs", "bundle", "release", "*-release*.aab"}
				aabFile, err = android.FindAabPath(arr, path)

				if err != nil {
					return err
				}
			}
		} else if filepath.Ext(path) == ".aab" {
			aabFile = path
		}

		if aabFile != "" && aabDir == "" {
			logger.Debug(fmt.Sprintf("Extracting AAB file: %s", aabFile))
			aabDir, err = utils.ExtractFile(aabFile, "aab")

			defer os.RemoveAll(aabDir)

			if err != nil {
				return err
			}
		}
	}

	manifestData, err = android.MergeUploadOptionsFromAabManifest(aabDir, globalOptions.ApiKey, aabOptions.ApplicationId, aabOptions.BuildUuid, aabOptions.NoBuildUuid, aabOptions.VersionCode, aabOptions.VersionName, logger)

	if err != nil {
		return err
	}

	soFilePath := filepath.Join(aabDir, "BUNDLE-METADATA", "com.android.tools.build.debugsymbols")

	if utils.FileExists(soFilePath) {
		logger.Debug(fmt.Sprintf("Found NDK (.so) files at: %s", soFilePath))
		soFileList, err := utils.BuildFileList([]string{soFilePath})

		if err != nil {
			return err
		}

		if len(soFileList) > 0 {
			globalOptions.Upload.AndroidNdk = options.AndroidNdkMapping{
				ApplicationId: manifestData["applicationId"],
				Path:          soFileList,
				ProjectRoot:   aabOptions.ProjectRoot,
				VersionCode:   manifestData["versionCode"],
				VersionName:   manifestData["versionName"],
			}
			globalOptions.ApiKey = manifestData["apiKey"]
			err = ProcessAndroidNDK(globalOptions, endpoint, logger)

			if err != nil {
				return err
			}
		} else {
			logger.Info("No NDK (.so) files detected for upload.")
		}
	} else {
		logger.Info("No NDK (.so) files detected for upload.")
	}

	mappingFilePath := filepath.Join(aabDir, "BUNDLE-METADATA", "com.android.tools.build.obfuscation", "proguard.map")

	if utils.FileExists(mappingFilePath) {
		logger.Debug(fmt.Sprintf("Found Proguard (mapping.txt) file at: %s", mappingFilePath))
		globalOptions.Upload.AndroidProguard = options.AndroidProguardMapping{
			ApplicationId: manifestData["applicationId"],
			BuildUuid:     manifestData["buildUuid"],
			NoBuildUuid:   aabOptions.NoBuildUuid,
			DexFiles:      []string{filepath.Join(aabDir, "base", "dex")},
			Path:          []string{mappingFilePath},
			VersionCode:   manifestData["versionCode"],
			VersionName:   manifestData["versionName"],
		}
		globalOptions.ApiKey = manifestData["apiKey"]
		err = ProcessAndroidProguard(globalOptions, endpoint, logger)

		if err != nil {
			return err
		}
	} else {
		logger.Info("No Proguard (mapping.txt) file detected for upload.")
	}

	return nil
}
