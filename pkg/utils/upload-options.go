package utils

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
func BuildAndroidNDKUploadOptions(apiKey string, appId string, versionName string, versionCode string, projectRoot string, sharedObjectName string, overwrite bool) map[string]string {
	uploadOptions := make(map[string]string)

	uploadOptions["apiKey"] = apiKey
	uploadOptions["appId"] = appId
	uploadOptions["versionName"] = versionName
	uploadOptions["versionCode"] = versionCode
	uploadOptions["projectRoot"] = projectRoot
	uploadOptions["sharedObjectName"] = sharedObjectName

	if overwrite {
		uploadOptions["overwrite"] = "true"
	}

	return uploadOptions
}

// BuildAndroidProguardUploadOptions - Builds the upload options for processing dart files
func BuildAndroidProguardUploadOptions(apiKey string, appId string, versionName string, versionCode string, uuid string, overwrite bool) map[string]string {
	uploadOptions := make(map[string]string)

	uploadOptions["apiKey"] = apiKey
	uploadOptions["appId"] = appId
	uploadOptions["versionName"] = versionName
	uploadOptions["versionCode"] = versionCode
	uploadOptions["buildUUID"] = uuid

	if overwrite {
		uploadOptions["overwrite"] = "true"
	}

	return uploadOptions
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

	if overwrite {
		uploadOptions["overwrite"] = "true"
	}

	return uploadOptions
}
