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
	ApplicationId string            `help:"Module application identifier"`
	BuildUuid     string            `help:"Module Build UUID ('none' to opt-out)"`
	Path          utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	ProjectRoot   string            `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	VersionCode   string            `help:"Module version code"`
	VersionName   string            `help:"Module version name"`
}

func ProcessAndroidAab(apiKey string, applicationId string, buildUuid string, paths []string, projectRoot string, versionCode string, versionName string, endpoint string, failOnUploadError bool, retries int, timeout int, overwrite bool, dryRun bool) error {

	var manifestData map[string]string
	var aabManifestPath string

	// Create temp working directory
	tempDir, err := os.MkdirTemp("", "bugsnag-cli-aab-unpacking-*")

	if err != nil {
		return fmt.Errorf("error creating temporary working directory " + err.Error())
	}

	defer os.RemoveAll(tempDir)

	for _, path := range paths {
		if filepath.Ext(path) == ".aab" {

			log.Info("Extracting " + filepath.Base(path) + " to " + tempDir)

			err = utils.Unzip(path, tempDir)

			if err != nil {
				return err
			}

			log.Success(filepath.Base(path) + " expanded")
		} else {
			return fmt.Errorf(path + " is not an AAB file")
		}
	}

	if aabManifestPath == "" {
		aabManifestPathExpected := filepath.Join(tempDir, "base", "manifest", "AndroidManifest.xml")
		if utils.FileExists(aabManifestPathExpected) {
			aabManifestPath = aabManifestPathExpected
		} else {
			log.Warn("AndroidManifest.xml not found in AAB file")
		}
	}

	// Check to see if we need to read the manifest file due to missing options
	if aabManifestPath != "" && (apiKey == "" || applicationId == "" || buildUuid == "" || versionCode == "" || versionName == "") {

		log.Info("Reading data from AndroidManifest.xml")

		manifestData, err = android.ReadAabManifest(filepath.Join(aabManifestPath))

		if err != nil {
			return fmt.Errorf("unable to read data from " + aabManifestPath + " " + err.Error())
		}

		if apiKey == "" {
			apiKey = manifestData["apiKey"]
			if apiKey != "" {
				log.Info("Using " + apiKey + " as API key from AndroidManifest.xml")
			}
		}

		if applicationId == "" {
			applicationId = manifestData["applicationId"]
			if applicationId != "" {
				log.Info("Using " + applicationId + " as application ID from AndroidManifest.xml")
			}
		}

		if buildUuid == "" {
			buildUuid = manifestData["buildUuid"]
			if buildUuid != "" {
				log.Info("Using " + buildUuid + " as build ID from AndroidManifest.xml")
			} else {
				buildUuid = android.GetDexBuildId(filepath.Join(tempDir, "base", "dex"))

				if buildUuid != "" {
					log.Info("Using " + buildUuid + " as build ID from dex signatures")
				}
			}
		} else if buildUuid == "none" {
			log.Info("No build ID will be used")
			buildUuid = ""
		}

		if versionCode == "" {
			versionCode = manifestData["versionCode"]
			if versionCode != "" {
				log.Info("Using " + versionCode + " as version code from AndroidManifest.xml")
			}
		}

		if versionName == "" {
			versionName = manifestData["versionName"]
			if versionName != "" {
				log.Info("Using " + versionName + " as version name from AndroidManifest.xml")
			}
		}
	}

	soFilePath := filepath.Join(tempDir, "BUNDLE-METADATA", "com.android.tools.build.debugsymbols")

	fileList, err := utils.BuildFileList([]string{soFilePath})

	if len(fileList) > 0 && err == nil {
		for _, file := range fileList {
			err = ProcessAndroidNDK(apiKey, applicationId, "", "", []string{file}, projectRoot, "", versionCode, versionName, endpoint, failOnUploadError, retries, timeout, overwrite, dryRun)

			if err != nil {
				return err
			}
		}
	} else {
		log.Info("No NDK (.so) files detected for upload. " + err.Error())
	}

	mappingFilePath := filepath.Join(tempDir, "BUNDLE-METADATA", "com.android.tools.build.obfuscation", "proguard.map")

	if utils.FileExists(mappingFilePath) {
		err = ProcessAndroidProguard(apiKey, applicationId, "", buildUuid, []string{mappingFilePath}, "", versionCode, versionName, endpoint, retries, timeout, overwrite, dryRun)

		if err != nil {
			return err
		}
	} else {
		log.Info("No Proguard (mapping.txt) file detected for upload.")
	}

	return nil
}
