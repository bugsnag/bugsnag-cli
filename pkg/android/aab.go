package android

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"path/filepath"
)

func GetUploadOptionsFromAabManifest(path string, apiKey string, applicationId string, buildUuid string, versionCode string, versionName string) (map[string]string, error) {

	var manifestData map[string]string
	var err error
	aabUploadOptions := make(map[string]string)

	manifestData, err = ReadAabManifest(filepath.Join(path))

	if err != nil {
		return nil, fmt.Errorf("unable to read data from " + path + " " + err.Error())
	}

	if apiKey == "" {
		apiKey = manifestData["apiKey"]
		if apiKey != "" {
			log.Info("Using " + apiKey + " as API key from AndroidManifest.xml")
		}
	}

	aabUploadOptions["apiKey"] = apiKey

	if applicationId == "" {
		applicationId = manifestData["applicationId"]
		if applicationId != "" {
			log.Info("Using " + applicationId + " as application ID from AndroidManifest.xml")
		}
	}

	aabUploadOptions["applicationId"] = applicationId

	if buildUuid == "" {
		buildUuid = manifestData["buildUuid"]
		if buildUuid != "" {
			log.Info("Using " + buildUuid + " as build ID from AndroidManifest.xml")
		} else {
			buildUuid = GetDexBuildId(filepath.Join(path, "..", "..", "dex"))

			if buildUuid != "" {
				log.Info("Using " + buildUuid + " as build ID from dex signatures")
			}
		}
	} else if buildUuid == "none" {
		log.Info("No build ID will be used")
		buildUuid = ""
	}

	aabUploadOptions["buildUuid"] = buildUuid

	if versionCode == "" {
		versionCode = manifestData["versionCode"]
		if versionCode != "" {
			log.Info("Using " + versionCode + " as version code from AndroidManifest.xml")
		}
	}

	aabUploadOptions["versionCode"] = versionCode

	if versionName == "" {
		versionName = manifestData["versionName"]
		if versionName != "" {
			log.Info("Using " + versionName + " as version name from AndroidManifest.xml")
		}
	}

	aabUploadOptions["versionName"] = versionName

	return aabUploadOptions, nil
}
