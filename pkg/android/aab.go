package android

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"path/filepath"
)

func MergeUploadOptionsFromAabManifest(
	path string,
	apiKey string,
	applicationId string,
	buildUuid string,
	noBuildUuid bool,
	versionCode string,
	versionName string,
	logger log.Logger,
) (map[string]string, error) {

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

			logger.Info("Reading data from AndroidManifest.xml")

			manifestData, err = ReadAabManifest(filepath.Join(aabManifestPath))

			if err != nil {
				return aabUploadOptions, fmt.Errorf("unable to read data from " + path + " " + err.Error())
			}
		} else {
			return aabUploadOptions, fmt.Errorf("AndroidManifest.xml not found in AAB file")
		}

		if aabUploadOptions["apiKey"] == "" && manifestData["apiKey"] != "" {
			aabUploadOptions["apiKey"] = manifestData["apiKey"]
			logger.Info("Using " + manifestData["apiKey"] + " as API key from AndroidManifest.xml")
		}

		if aabUploadOptions["applicationId"] == "" && manifestData["applicationId"] != "" {
			aabUploadOptions["applicationId"] = manifestData["applicationId"]
			logger.Info("Using " + aabUploadOptions["applicationId"] + " as application ID from AndroidManifest.xml")
		}

		if aabUploadOptions["buildUuid"] == "" && !noBuildUuid {
			aabUploadOptions["buildUuid"] = manifestData["buildUuid"]
			if aabUploadOptions["buildUuid"] != "" {
				logger.Info("Using " + aabUploadOptions["buildUuid"] + " as build ID from AndroidManifest.xml")
			} else {
				aabUploadOptions["buildUuid"] = GetDexBuildId(filepath.Join(path, "base", "dex"))
				if aabUploadOptions["buildUuid"] != "" {
					logger.Info("Using " + aabUploadOptions["buildUuid"] + " as build ID from dex signatures")
				}
			}
		} else if aabUploadOptions["buildUuid"] == "none" || noBuildUuid {
			logger.Info("No build ID will be used")
			aabUploadOptions["buildUuid"] = ""
		}

		if aabUploadOptions["versionCode"] == "" && manifestData["versionCode"] != "" {
			aabUploadOptions["versionCode"] = manifestData["versionCode"]
			logger.Info("Using " + aabUploadOptions["versionCode"] + " as version code from AndroidManifest.xml")
		}

		if aabUploadOptions["versionName"] == "" && manifestData["versionName"] != "" {
			aabUploadOptions["versionName"] = manifestData["versionName"]
			logger.Info("Using " + aabUploadOptions["versionName"] + " as version name from AndroidManifest.xml")
		}
	}
	return aabUploadOptions, nil
}
