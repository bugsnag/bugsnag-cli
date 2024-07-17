package upload

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
	"path/filepath"
)

type AndroidAabMapping struct {
	ApplicationId string      `help:"Module application identifier"`
	BuildUuid     string      `help:"Module Build UUID" xor:"no-build-uuid,build-uuid"`
	NoBuildUuid   bool        `help:"Upload with no Build UUID" xor:"build-uuid,no-build-uuid"`
	Path          utils.Paths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path" default:"."`
	ProjectRoot   string      `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	VersionCode   string      `help:"Module version code"`
	VersionName   string      `help:"Module version name"`
}

func ProcessAndroidAab(
	apiKey string,
	applicationId string,
	buildUuid string,
	noBuildUuid bool,
	paths []string,
	projectRoot string,
	versionCode string,
	versionName string,
	endpoint string,
	retries int,
	timeout int,
	overwrite bool,
	dryRun bool,
	logger log.Logger,
) error {

	var manifestData map[string]string
	var aabDir string
	var err error

	for _, path := range paths {
		// Look for AAB file if the upload command was run somewhere within the project root
		// based on an expected path of ${dir}/build/outputs/bundle/release/${dir}-release.aab
		// or ${dir}/build/outputs/bundle/release/${dir}-release-dexguard.aab
		if utils.IsDir(path) {
			if utils.FileExists(filepath.Join(path, "BUNDLE-METADATA")) {
				aabDir = path
			} else {
				arr := []string{"*", "build", "outputs", "bundle", "release", "*-release*.aab"}
				path, err = android.FindAabPath(arr, path)

				if err != nil {
					return err
				}
			}
		} else if filepath.Ext(path) == ".aab" {
			aabDir, err = utils.ExtractFile(path, "aab")

			defer os.RemoveAll(aabDir)

			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("%s is not an AAB file/directory", path)
		}
	}

	manifestData, err = android.MergeUploadOptionsFromAabManifest(aabDir, apiKey, applicationId, buildUuid, noBuildUuid, versionCode, versionName, logger)

	if err != nil {
		return err
	}

	soFilePath := filepath.Join(aabDir, "BUNDLE-METADATA", "com.android.tools.build.debugsymbols")

	if utils.FileExists(soFilePath) {
		soFileList, err := utils.BuildFileList([]string{soFilePath})

		if err != nil {
			return err
		}

		if len(soFileList) > 0 {
			err = ProcessAndroidNDK(
				manifestData["apiKey"],
				manifestData["applicationId"],
				"",
				"",
				soFileList,
				projectRoot,
				"",
				manifestData["versionCode"],
				manifestData["versionName"],
				endpoint,
				retries,
				timeout,
				overwrite,
				dryRun,
				logger,
			)

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
		err = ProcessAndroidProguard(
			manifestData["apiKey"],
			manifestData["applicationId"],
			"",
			manifestData["buildUuid"],
			noBuildUuid,
			[]string{filepath.Join(aabDir, "base", "dex")},
			[]string{mappingFilePath},
			"",
			manifestData["versionCode"],
			manifestData["versionName"],
			endpoint,
			retries,
			timeout,
			overwrite,
			dryRun,
			logger,
		)

		if err != nil {
			return err
		}
	} else {
		logger.Info("No Proguard (mapping.txt) file detected for upload.")
	}

	return nil
}
