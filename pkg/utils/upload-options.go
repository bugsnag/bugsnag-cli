package utils

import (
	"fmt"
)

// BuildDartUploadOptions - Builds the upload options for processing dart files
func BuildDartUploadOptions(apiKey string, uuid string, platform string, overwrite bool, appVersion string, appExtraVersion string) map[string]string {
	uploadOptions := make(map[string]string)

	uploadOptions["apiKey"] = apiKey

	uploadOptions["buildId"] = uuid

	uploadOptions["platform"] = platform

	if overwrite {
		uploadOptions["overwrite"] = "true"
	}

	if platform == "ios" {
		if appVersion != "" {
			uploadOptions["appVersion"] = appVersion
		}

		if appExtraVersion != "" {
			uploadOptions["AppBundleVersion"] = appExtraVersion
		}
	}

	if platform == "android" {
		if appVersion != "" {
			uploadOptions["appVersion"] = appVersion
		}

		if appExtraVersion != "" {
			uploadOptions["appVersionCode"] = appExtraVersion
		}
	}

	return uploadOptions
}

// BuildAndroidNDKUploadOptions - Builds the upload options for processing dart files
func BuildAndroidNDKUploadOptions(apiKey string, applicationId string, versionName string, versionCode string, projectRoot string, sharedObjectName string, overwrite bool) (map[string]string, error) {
	uploadOptions := make(map[string]string)

	if apiKey != "" {
		uploadOptions["apiKey"] = apiKey
	} else {
		return nil, fmt.Errorf("missing api key, please specify using `--api-key`")
	}

	if applicationId != "" {
		uploadOptions["appId"] = applicationId
	}

	if versionCode != "" {
		uploadOptions["versionCode"] = versionCode
	}

	if versionName != "" {
		uploadOptions["versionName"] = versionName
	}

	if projectRoot != "" {
		uploadOptions["projectRoot"] = projectRoot
	}

	if sharedObjectName != "" {
		uploadOptions["sharedObjectName"] = sharedObjectName
	}

	if overwrite {
		uploadOptions["overwrite"] = "true"
	}

	if uploadOptions["appId"] == "" && uploadOptions["versionName"] == "" && uploadOptions["versionCode"] == "" {
		return nil, fmt.Errorf("missing application ID, version code or version name, please specify using `--application-id`, `--version-code` or `--version-name`")
	}

	return uploadOptions, nil
}

// BuildAndroidProguardUploadOptions - Builds the upload options for processing dart files
func BuildAndroidProguardUploadOptions(apiKey string, applicationId string, versionName string, versionCode string, buildUuid string, overwrite bool) (map[string]string, error) {
	uploadOptions := make(map[string]string)

	if apiKey != "" {
		uploadOptions["apiKey"] = apiKey
	} else {
		return nil, fmt.Errorf("missing api key, please specify using `--api-key`")
	}

	if applicationId != "" {
		uploadOptions["appId"] = applicationId
	}

	if versionCode != "" {
		uploadOptions["versionCode"] = versionCode
	}

	if versionName != "" {
		uploadOptions["versionName"] = versionName
	}

	if buildUuid != "" {
		uploadOptions["buildUuid"] = buildUuid
	}

	if overwrite {
		uploadOptions["overwrite"] = "true"
	}

	if uploadOptions["buildUuid"] == "" && uploadOptions["versionName"] == "" && uploadOptions["versionCode"] == "" {
		return nil, fmt.Errorf("you must set one at least the application ID, version name, version code or build UUID to uniquely identify the build")
	}

	return uploadOptions, nil
}

func BuildReactNativeAndroidUploadOptions(apiKey string, appVersion string, appVersionCode string, codeBundleId string, dev bool, projectRoot string, overwrite bool) map[string]string {
	uploadOptions := make(map[string]string)

	uploadOptions["apiKey"] = apiKey
	uploadOptions["appVersion"] = appVersion

	if appVersionCode != "" {
		uploadOptions["appVersionCode"] = appVersionCode
	}

	if codeBundleId != "" {
		uploadOptions["codeBundleId"] = codeBundleId
	}

	if dev {
		uploadOptions["dev"] = "true"
	}

	uploadOptions["projectRoot"] = projectRoot

	uploadOptions["platform"] = "android"

	if overwrite {
		uploadOptions["overwrite"] = "true"
	}

	return uploadOptions
}
