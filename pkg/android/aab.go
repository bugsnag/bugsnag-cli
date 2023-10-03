package android

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"path/filepath"
)

func MergeUploadOptionsFromAabManifest(path string, apiKey string, applicationId string, buildUuid string, versionCode string, versionName string) (map[string]string, error) {

	var manifestData map[string]string
	var err error
	var aabManifestPath string
	aabUploadOptions := make(map[string]string)

	aabUploadOptions["apiKey"] = apiKey
	aabUploadOptions["applicationId"] = applicationId
	aabUploadOptions["buildUuid"] = buildUuid
	aabUploadOptions["versionCode"] = versionCode
	aabUploadOptions["versionName"] = versionName

	if apiKey == "" || applicationId == "" || buildUuid == "" || versionCode == "" || versionName == "" {

		aabManifestPathExpected := filepath.Join(path, "base", "manifest", "AndroidManifest.xml")

		if utils.FileExists(aabManifestPathExpected) {
			aabManifestPath = aabManifestPathExpected
		} else {
			return nil, fmt.Errorf("AndroidManifest.xml not found in AAB file")
		}

		log.Info("Reading data from AndroidManifest.xml")

		manifestData, err = ReadAabManifest(filepath.Join(aabManifestPath))

		if err != nil {
			return nil, fmt.Errorf("unable to read data from " + path + " " + err.Error())
		}

		if aabUploadOptions["apiKey"] == "" {
			aabUploadOptions["apiKey"] = manifestData["apiKey"]
			if aabUploadOptions["apiKey"] != "" {
				log.Info("Using " + aabUploadOptions["apiKey"] + " as API key from AndroidManifest.xml")
			}
		}

		if aabUploadOptions["applicationId"] == "" {
			aabUploadOptions["applicationId"] = manifestData["applicationId"]
			if aabUploadOptions["applicationId"] != "" {
				log.Info("Using " + aabUploadOptions["applicationId"] + " as application ID from AndroidManifest.xml")
			}
		}

		if aabUploadOptions["buildUuid"] == "" {
			aabUploadOptions["buildUuid"] = manifestData["buildUuid"]
			if aabUploadOptions["buildUuid"] != "" {
				log.Info("Using " + aabUploadOptions["buildUuid"] + " as build ID from AndroidManifest.xml")
			} else {
				aabUploadOptions["buildUuid"] = GetDexBuildId(filepath.Join(path, "base", "dex"))
				if aabUploadOptions["buildUuid"] != "" {
					log.Info("Using " + aabUploadOptions["buildUuid"] + " as build ID from dex signatures")
				}
			}
		} else if aabUploadOptions["buildUuid"] == "none" {
			log.Info("No build ID will be used")
			aabUploadOptions["buildUuid"] = ""
		}

		if aabUploadOptions["versionCode"] == "" {
			aabUploadOptions["versionCode"] = manifestData["versionCode"]
			if aabUploadOptions["versionCode"] != "" {
				log.Info("Using " + aabUploadOptions["versionCode"] + " as version code from AndroidManifest.xml")
			}
		}

		if aabUploadOptions["versionName"] == "" {
			aabUploadOptions["versionName"] = manifestData["versionName"]
			if aabUploadOptions["versionName"] != "" {
				log.Info("Using " + aabUploadOptions["versionName"] + " as version name from AndroidManifest.xml")
			}
		}
		return aabUploadOptions, nil
	}
	return aabUploadOptions, nil
}
