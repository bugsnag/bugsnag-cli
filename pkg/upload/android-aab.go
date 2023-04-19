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
	AndroidNdkRoot string            `help:"Path to Android NDK installation ($ANDROID_NDK_ROOT)"`
	Path           utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	ProjectRoot    string            `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
}

func ProcessAndroidAab(apiKey string, androidNdkRoot string, paths []string, projectRoot string, endpoint string, failOnUploadError bool, retries int, timeout int, overwrite bool, dryRun bool) error {

	var manifestData map[string]string

	// Create temp directory
	tempDir, err := os.MkdirTemp("", "bugsnag-cli-aab-unpacking-*")

	if err != nil {
		return err
	}

	if dryRun {
		log.Info("Performing dry run - no files will be uploaded")
	}

	for _, path := range paths {
		if filepath.Ext(path) == ".aab" {

			log.Info("Extracting " + filepath.Base(path) + " to " + tempDir)

			err = utils.Unzip(path, tempDir)

			if err != nil {
				return err
			}

			log.Success(filepath.Base(path) + " expanded")

			aabManifestPath := filepath.Join(tempDir, "base", "manifest", "AndroidManifest.xml")

			if utils.FileExists(aabManifestPath) {
				manifestData, err = android.ReadAabManifest(filepath.Join(aabManifestPath))

				if err != nil {
					return fmt.Errorf("error reading raw AAB manifest data. " + err.Error())
				}
			} else {
				return fmt.Errorf("unable to read data from " + aabManifestPath + " " + err.Error())
			}
		} else {
			// Not an AAB file
			return fmt.Errorf("")
		}
	}

	if apiKey == "" {
		apiKey = manifestData["apiKey"]
		log.Info("Using " + apiKey + " as API key from AndroidManifest.xml")
	}

	applicationId := manifestData["applicationId"]
	log.Info("Using " + applicationId + " as application ID from AndroidManifest.xml")

	versionCode := manifestData["versionCode"]
	log.Info("Using " + versionCode + " as version code from AndroidManifest.xml")

	versionName := manifestData["versionName"]
	log.Info("Using " + versionName + " as version name from AndroidManifest.xml")

	buildUuid := manifestData["buildUuid"]
	log.Info("Using " + buildUuid + " as build UUID from AndroidManifest.xml")

	soFilePath := filepath.Join(tempDir, "BUNDLE-METADATA", "com.android.tools.build.debugsymbols")

	fileList, err := utils.BuildFileList([]string{soFilePath})

	if err != nil {
		return fmt.Errorf("error building `.so` file list. " + err.Error())
	}

	mappingFilePath := filepath.Join(tempDir, "BUNDLE-METADATA", "com.android.tools.build.obfuscation", "proguard.map")

	for _, file := range fileList {
		err = ProcessAndroidNDK(apiKey, applicationId, androidNdkRoot, "", []string{file}, projectRoot, "", versionCode, versionName, endpoint+"/ndk-symbol", failOnUploadError, retries, timeout, overwrite, dryRun)

		if err != nil {
			return err
		}
	}

	err = ProcessAndroidProguard(apiKey, applicationId, "", buildUuid, []string{mappingFilePath}, "", versionCode, versionName, endpoint, retries, timeout, overwrite, dryRun)

	if err != nil {
		return err
	}

	return nil
}
