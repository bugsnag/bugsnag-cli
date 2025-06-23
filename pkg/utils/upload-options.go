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

func BuildDsymUploadOptions(apiKey string, projectRoot string) (map[string]string, error) {
	uploadOptions := make(map[string]string)

	if apiKey != "" {
		uploadOptions["apiKey"] = apiKey
	} else {
		return nil, fmt.Errorf("missing api key, please specify using `--api-key`")
	}

	uploadOptions["projectRoot"] = projectRoot

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

		return nil, fmt.Errorf("you must set at least the version name, %s and code bundle ID to uniquely identify the build", platformSpecificTerminology)
	}

	if apiKey != "" {
		uploadOptions["apiKey"] = apiKey
	} else {
		return nil, fmt.Errorf("missing api key, please specify using `--api-key`")
	}

	// If codeBundleId is set, use that instead of appVersion and versionCode
	if codeBundleId != "" {
		uploadOptions["codeBundleId"] = codeBundleId
	} else {
		uploadOptions["appVersion"] = appVersion

		if platform == "android" {
			uploadOptions["appVersionCode"] = versionCode
		} else if platform == "ios" {
			uploadOptions["appBundleVersion"] = versionCode
		}
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

func BuildJsUploadOptions(apiKey string, versionName string, codeBundleId string, bundleUrl string, projectRoot string, overwrite bool) (map[string]string, error) {
	uploadOptions := make(map[string]string)

	if apiKey != "" {
		uploadOptions["apiKey"] = apiKey
	} else {
		return nil, fmt.Errorf("missing api key, please specify using `--api-key`")
	}

	if versionName != "" {
		uploadOptions["appVersion"] = versionName
	}

	if bundleUrl != "" {
		uploadOptions["minifiedUrl"] = bundleUrl
	} else {
		return nil, fmt.Errorf("missing minified URL, please specify using `--bundle-url`")
	}

	if projectRoot != "" {
		uploadOptions["projectRoot"] = projectRoot
	}

	if codeBundleId != "" {
		uploadOptions["codeBundleId"] = codeBundleId
	}

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

// BuildBreakpadUploadOptions - Builds the upload options for processing breakpad symbol files
func BuildBreakpadUploadOptions(CpuArch string, CodeFile string, DebugFile string, DebugIdentifier string, ProductName string, OsName string, VersionName string) (map[string]string, error) {
	uploadOptions := make(map[string]string)

	if CpuArch != "" {
		uploadOptions["cpu"] = CpuArch
	}

	if CodeFile != "" {
		uploadOptions["code_file"] = CodeFile
	}

	if DebugFile != "" {
		uploadOptions["debug_file"] = DebugFile
	}

	if DebugIdentifier != "" {
		uploadOptions["debug_identifier"] = DebugIdentifier
	}

	if ProductName != "" {
		uploadOptions["product"] = ProductName
	}

	if OsName != "" {
		uploadOptions["os"] = OsName
	}

	if VersionName != "" {
		uploadOptions["version"] = VersionName
	}

	return uploadOptions, nil
}

type UnityLineMappingOptions struct {
	APIKey           string
	AppID            string
	AppVersion       string
	AppVersionCode   string // Android only
	AppBundleVersion string // iOS only
	SOBuildID        string // Android only
	DSYMUUUID        string // iOS only
	ProjectRoot      string
	Overwrite        bool
}

func BuildUnityLineMappingUploadOptions(opts UnityLineMappingOptions) (map[string]string, error) {
	uploadOptions := make(map[string]string)

	if opts.APIKey != "" {
		uploadOptions["apiKey"] = opts.APIKey
	} else {
		return nil, fmt.Errorf("missing api key, please specify using `--api-key`")
	}

	if opts.SOBuildID != "" {
		uploadOptions["soBuildId"] = opts.SOBuildID
	}

	if opts.DSYMUUUID != "" {
		uploadOptions["dsymUUID"] = opts.DSYMUUUID
	}

	if opts.AppID != "" {
		uploadOptions["appId"] = opts.AppID
	}

	if opts.AppVersionCode != "" {
		uploadOptions["appVersionCode"] = opts.AppVersionCode
	}

	if opts.AppBundleVersion != "" {
		uploadOptions["appBundleVersion"] = opts.AppBundleVersion
	}

	if opts.AppVersion != "" {
		uploadOptions["appVersion"] = opts.AppVersion
	}

	if opts.ProjectRoot != "" {
		uploadOptions["projectRoot"] = opts.ProjectRoot
	}

	if opts.Overwrite {
		uploadOptions["overwrite"] = "true"
	}

	return uploadOptions, nil
}
