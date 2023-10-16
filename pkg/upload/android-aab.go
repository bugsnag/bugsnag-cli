package upload

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type AndroidAabMapping struct {
	ApplicationId string      `help:"Module application identifier"`
	BuildUuid     string      `help:"Module Build UUID ('none' to opt-out)"`
	Path          utils.Paths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	ProjectRoot   string      `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	VersionCode   string      `help:"Module version code"`
	VersionName   string      `help:"Module version name"`
}

func ProcessAndroidAab(apiKey string, applicationId string, buildUuid string, paths []string, projectRoot string, versionCode string, versionName string, endpoint string, failOnUploadError bool, retries int, timeout int, overwrite bool, dryRun bool) error {

	var manifestData map[string]string
	var aabDir string
	var err error

	for _, path := range paths {
		// Check to see if we are dealing with a .aab file and extract it into a temp directory
		if filepath.Ext(path) == ".aab" {

			log.Info("Extracting " + filepath.Base(path) + " into a temporary directory")
			aabDir, err = utils.ExtractFile(path, "aab")

			defer os.RemoveAll(aabDir)

			if err != nil {
				return err
			}
		} else if utils.IsDir(path) {
			aabDir = path
		} else {
			return fmt.Errorf(path + " is not an AAB file/directory")
		}
	}

	manifestData, err = android.MergeUploadOptionsFromAabManifest(aabDir, apiKey, applicationId, buildUuid, versionCode, versionName)

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
			for _, file := range soFileList {
				err = ProcessAndroidNDK(manifestData["apiKey"], manifestData["applicationId"], "", "", []string{file}, projectRoot, "", manifestData["versionCode"], manifestData["versionName"], endpoint, failOnUploadError, retries, timeout, overwrite, dryRun)

				if err != nil {
					return err
				}
			}
		} else {
			log.Info("No NDK (.so) files detected for upload.")
		}
	} else {
		log.Info("No NDK (.so) files detected for upload.")
	}

	mappingFilePath := filepath.Join(aabDir, "BUNDLE-METADATA", "com.android.tools.build.obfuscation", "proguard.map")

	if utils.FileExists(mappingFilePath) {
		err = ProcessAndroidProguard(
			manifestData["apiKey"],
			manifestData["applicationId"],
			"",
			manifestData["buildUuid"],
			nil,
			[]string{mappingFilePath},
			"",
			manifestData["versionCode"],
			manifestData["versionName"],
			endpoint,
			retries,
			timeout,
			overwrite,
			dryRun,
		)

		if err != nil {
			return err
		}
	} else {
		log.Info("No Proguard (mapping.txt) file detected for upload.")
	}

	return nil
}
