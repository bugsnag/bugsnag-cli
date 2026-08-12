package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

// BuildDartUploadOptions - Builds the upload options for processing dart files
func BuildDartUploadOptions(uuid string, platform string, overwrite bool, appVersion string, appExtraVersion string) map[string]string {
	uploadOptions := make(map[string]string)

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
func BuildAndroidProguardUploadOptions(applicationId string, versionName string, versionCode string, buildUuid string, overwrite bool) (map[string]string, error) {
	uploadOptions := make(map[string]string)

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

func BuildDsymUploadOptions(projectRoot string) (map[string]string, error) {
	uploadOptions := make(map[string]string)

	// Normalize projectRoot to a clean absolute path.
	normalizedRoot, err := normalizeProjectRoot(projectRoot)
	if err != nil {
		return nil, err
	}
	if normalizedRoot != "" {
		uploadOptions["projectRoot"] = normalizedRoot
	}

	return uploadOptions, nil
}

// normalizeProjectRoot converts a path to a cleaned absolute path, handling various input formats:
//   - Empty string: returned as-is (treated as "no project root" — no stripping performed)
//   - Absolute path (starts with "/"): cleaned with filepath.Clean to collapse duplicate/trailing slashes
//   - Dot-relative path (starts with "." e.g. ".", "../../../"): resolved to absolute via the CWD
//   - Other (no leading slash): treated as an absolute path with a missing "/", prepended and cleaned
//
// Returns an error only if the path exceeds 1024 characters or CWD resolution fails.
func normalizeProjectRoot(projectRoot string) (string, error) {
	if projectRoot == "" {
		return "", nil
	}

	const maxPathLength = 1024
	if len(projectRoot) > maxPathLength {
		return "", fmt.Errorf("project root path exceeds the maximum allowed length of %d characters", maxPathLength)
	}

	var resolved string

	switch {
	case strings.HasPrefix(projectRoot, "/"):
		// Already absolute: clean to collapse double/trailing slashes
		resolved = filepath.Clean(projectRoot)
	case strings.HasPrefix(projectRoot, "."):
		// Dot-relative (e.g. ".", "../../../"): resolve against CWD
		abs, err := filepath.Abs(projectRoot)
		if err != nil {
			return "", fmt.Errorf("failed to resolve relative project root %q: %w", projectRoot, err)
		}
		resolved = abs
	default:
		// No leading slash: treat as absolute path with missing "/", clean after prepending
		resolved = filepath.Clean("/" + projectRoot)
	}

	return resolved, nil
}

func BuildJsUploadOptions(versionName string, codeBundleId string, bundleUrl string, projectRoot string, overwrite bool) (map[string]string, error) {
	uploadOptions := make(map[string]string)

	// If codeBundleId is set, use that instead of appVersion
	if codeBundleId != "" {
		uploadOptions["codeBundleId"] = codeBundleId
	} else if versionName != "" {
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
