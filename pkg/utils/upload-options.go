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
			uploadOptions["appBundleVersion"] = appExtraVersion
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

// BuildAndroidProguardUploadOptions - Builds the upload options for processing Proguard files
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
		uploadOptions["buildUUID"] = buildUuid
	}

	if overwrite {
		uploadOptions["overwrite"] = "true"
	}

	if uploadOptions["appId"] == "" && uploadOptions["buildUUID"] == "" && uploadOptions["versionName"] == "" && uploadOptions["versionCode"] == "" {
		return nil, fmt.Errorf("you must set at least the application ID, version name, version code or build UUID to uniquely identify the build")
	}

	return uploadOptions, nil
}

func BuildReactNativeUploadOptions(apiKey string, appVersion string, versionCode string, codeBundleId string, dev bool, projectRoot string, overwrite bool, platform string) (map[string]string, error) {
	uploadOptions := make(map[string]string)

	// Return early if all three of these are empty/undefined
	if appVersion == "" && versionCode == "" && codeBundleId == "" {
		var platformSpecificTerminology string
		if platform == "android" {
			platformSpecificTerminology = "version code"
		} else if platform == "ios" {
			platformSpecificTerminology = "bundle version"
		}

		return nil, fmt.Errorf("you must set at least the version name, %s or code bundle ID to uniquely identify the build", platformSpecificTerminology)
	}

	if apiKey != "" {
		uploadOptions["apiKey"] = apiKey
	} else {
		return nil, fmt.Errorf("missing api key, please specify using `--api-key`")
	}

	uploadOptions["appVersion"] = appVersion

	if platform == "android" {
		uploadOptions["appVersionCode"] = versionCode

	} else if platform == "ios" {
		uploadOptions["appBundleVersion"] = versionCode
	}

	if codeBundleId != "" {
		uploadOptions["codeBundleId"] = codeBundleId
	}

	if dev {
		uploadOptions["dev"] = "true"
	}

	uploadOptions["projectRoot"] = projectRoot

	uploadOptions["platform"] = platform

	if overwrite {
		uploadOptions["overwrite"] = "true"
	}

	return uploadOptions, nil
}

// BuildAndroidNDKUploadOptions - Builds the upload options for processing NDK files
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

	return uploadOptions, nil
}
