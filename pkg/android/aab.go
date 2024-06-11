package android

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"path/filepath"
)

func FindAabPath(arr []string, path string) (string, error) {

	// Look for AAB file based on an expected path
	iterations := len(arr)
	for i := 1; i <= iterations; i++ {
		path_ending := filepath.Join(arr...)
		combinedPath := filepath.Join(path, path_ending)
		matchingPaths, err := filepath.Glob(combinedPath)

		if err != nil {
			return "", err
		}
		if matchingPaths != nil {
			if len(matchingPaths) > 1 {
				// Return an error if more than one AAB file was found
				return "", fmt.Errorf("Path ambiguous: more than one AAB file was found within %s", filepath.Dir(combinedPath))
			}
			if len(matchingPaths) == 1 {
				aabPath := matchingPaths[0]
				return aabPath, nil
			}
		}
		arr = arr[1:]
	}
	return "", fmt.Errorf("No AAB file was found")
}

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

			logger.Debug("Reading data from AndroidManifest.xml")

			manifestData, err = ReadAabManifest(filepath.Join(aabManifestPath))

			if err != nil {
				return aabUploadOptions, fmt.Errorf("unable to read data from %s %s", path, err.Error())
			}
		} else {
			return aabUploadOptions, fmt.Errorf("AndroidManifest.xml not found in AAB file")
		}

		if aabUploadOptions["apiKey"] == "" && manifestData["apiKey"] != "" {
			aabUploadOptions["apiKey"] = manifestData["apiKey"]
			logger.Debug(fmt.Sprintf("Using %s as API key from AndroidManifest.xml", manifestData["apiKey"]))
		}

		if aabUploadOptions["applicationId"] == "" && manifestData["applicationId"] != "" {
			aabUploadOptions["applicationId"] = manifestData["applicationId"]
			logger.Debug(fmt.Sprintf("Using %s as application ID from AndroidManifest.xml", aabUploadOptions["applicationId"]))
		}

		if aabUploadOptions["buildUuid"] == "" && !noBuildUuid {
			aabUploadOptions["buildUuid"] = manifestData["buildUuid"]
			if aabUploadOptions["buildUuid"] != "" {
				logger.Debug(fmt.Sprintf("Using %s as build ID from AndroidManifest.xml", aabUploadOptions["buildUuid"]))
			} else {
				aabUploadOptions["buildUuid"] = GetDexBuildId(filepath.Join(path, "base", "dex"))
				if aabUploadOptions["buildUuid"] != "" {
					logger.Debug(fmt.Sprintf("Using %s as build ID from dex signatures", aabUploadOptions["buildUuid"]))
				}
			}
		} else if aabUploadOptions["buildUuid"] == "none" || noBuildUuid {
			logger.Debug("No build ID will be used")
			aabUploadOptions["buildUuid"] = ""
		}

		if aabUploadOptions["versionCode"] == "" && manifestData["versionCode"] != "" {
			aabUploadOptions["versionCode"] = manifestData["versionCode"]
			logger.Debug(fmt.Sprintf("Using %s as version code from AndroidManifest.xml", aabUploadOptions["versionCode"]))
		}

		if aabUploadOptions["versionName"] == "" && manifestData["versionName"] != "" {
			aabUploadOptions["versionName"] = manifestData["versionName"]
			logger.Debug(fmt.Sprintf("Using %s as version name from AndroidManifest.xml", aabUploadOptions["versionName"]))
		}
	}
	return aabUploadOptions, nil
}
