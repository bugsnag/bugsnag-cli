package android

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"path/filepath"
)

func MergeUploadOptionsFromAabManifest(path string, apiKey string, applicationId string, buildUuid string, noBuildUuid bool, versionCode string, versionName string) (map[string]string, error) {

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

			log.Info("Reading data from AndroidManifest.xml")

			manifestData, err = ReadAabManifest(filepath.Join(aabManifestPath))

			if err != nil {
				return aabUploadOptions, fmt.Errorf("unable to read data from %s %s", path, err.Error())
			}
		} else {
			return aabUploadOptions, fmt.Errorf("AndroidManifest.xml not found in AAB file")
		}

		if aabUploadOptions["apiKey"] == "" && manifestData["apiKey"] != "" {
			aabUploadOptions["apiKey"] = manifestData["apiKey"]
			log.Info(fmt.Sprintf("Using %s as API key from AndroidManifest.xml", manifestData["apiKey"]))
		}

		if aabUploadOptions["applicationId"] == "" && manifestData["applicationId"] != "" {
			aabUploadOptions["applicationId"] = manifestData["applicationId"]
			log.Info(fmt.Sprintf("Using %s as application ID from AndroidManifest.xml", aabUploadOptions["applicationId"]))
		}

		if aabUploadOptions["buildUuid"] == "" && !noBuildUuid {
			aabUploadOptions["buildUuid"] = manifestData["buildUuid"]
			if aabUploadOptions["buildUuid"] != "" {
				log.Info(fmt.Sprintf("Using %s as build ID from AndroidManifest.xml", aabUploadOptions["buildUuid"]))
			} else {
				aabUploadOptions["buildUuid"] = GetDexBuildId(filepath.Join(path, "base", "dex"))
				if aabUploadOptions["buildUuid"] != "" {
					log.Info(fmt.Sprintf("Using %s as build ID from dex signatures", aabUploadOptions["buildUuid"]))
				}
			}
		} else if aabUploadOptions["buildUuid"] == "none" || noBuildUuid {
			log.Info("No build ID will be used")
			aabUploadOptions["buildUuid"] = ""
		}

		if aabUploadOptions["versionCode"] == "" && manifestData["versionCode"] != "" {
			aabUploadOptions["versionCode"] = manifestData["versionCode"]
			log.Info(fmt.Sprintf("Using %s as version code from AndroidManifest.xml", aabUploadOptions["versionCode"]))
		}

		if aabUploadOptions["versionName"] == "" && manifestData["versionName"] != "" {
			aabUploadOptions["versionName"] = manifestData["versionName"]
			log.Info(fmt.Sprintf("Using %s as version name from AndroidManifest.xml", aabUploadOptions["versionName"]))
		}
	}
	return aabUploadOptions, nil
}
