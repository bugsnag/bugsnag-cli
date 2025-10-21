package upload

import (
	"fmt"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
)

// ProcessReactNativeSourcemaps handles uploading React Native source maps and bundles to Bugsnag.
//
// This function prepares and sends a multipart request containing the JavaScript bundle
// and its corresponding source map, along with metadata identifying the build.
// It validates key identifiers (version name, code bundle ID, etc.) and constructs
// upload parameters according to React Native platform conventions.
//
// Parameters:
//   - globalOptions: Global CLI options including the API key and upload configuration.
//   - logger: Logger instance for structured debug, info, and error output.
//
// Behavior:
//   - Validates presence of required identifiers (versionName, versionCode/bundleVersion, codeBundleId).
//   - Builds metadata for either iOS or Android React Native builds.
//   - Uploads both the source map and the JS bundle to the Bugsnag /react-native-source-map endpoint.
//
// Returns:
//   - error: Non-nil if validation or upload fails.
func ProcessReactNativeSourcemaps(globalOptions options.CLI, logger log.Logger) error {
	reactNativeOpts := globalOptions.Upload.ReactNativeSourcemaps
	uploadOpts := make(map[string]string)
	// Validate versioning identifiers
	if reactNativeOpts.VersionName == "" && reactNativeOpts.CodeBundleId == "" {
		var versionLabel string
		switch reactNativeOpts.Platform {
		case "android":
			versionLabel = "version name, version code, or code bundle ID"
		case "ios":
			versionLabel = "version name, bundle version, or code bundle ID"
		default:
			versionLabel = "version name"
		}

		return fmt.Errorf(
			"missing required identifiers: you must set at least %s to uniquely identify the build",
			versionLabel,
		)
	}

	// Select identification strategy
	if reactNativeOpts.CodeBundleId != "" {
		logger.Info(fmt.Sprintf("Using code bundle ID to identify build: %s", reactNativeOpts.CodeBundleId))
		uploadOpts["codeBundleId"] = reactNativeOpts.CodeBundleId
	} else {
		logger.Info("Using version name and version code/bundle version to identify build")
		uploadOpts["appVersion"] = reactNativeOpts.VersionName

		if reactNativeOpts.VersionCode != "" {
			logger.Debug(fmt.Sprintf("Using Android version code: %s", reactNativeOpts.VersionCode))
			uploadOpts["appVersionCode"] = reactNativeOpts.VersionCode
		} else if reactNativeOpts.BundleVersion != "" {
			logger.Debug(fmt.Sprintf("Using iOS bundle version: %s", reactNativeOpts.BundleVersion))
			uploadOpts["appBundleVersion"] = reactNativeOpts.BundleVersion
		}
	}

	// Development mode flag
	if reactNativeOpts.Dev {
		logger.Debug("Development build detected — marking upload as dev=true")
		uploadOpts["dev"] = "true"
	}

	// Determine project root
	if reactNativeOpts.ProjectRoot == "" {
		uploadOpts["projectRoot"] = string(reactNativeOpts.Path)
		logger.Debug(fmt.Sprintf("Project root not provided — using inferred path: %s", uploadOpts["projectRoot"]))
	} else {
		uploadOpts["projectRoot"] = reactNativeOpts.ProjectRoot
		logger.Debug(fmt.Sprintf("Using specified project root: %s", reactNativeOpts.ProjectRoot))
	}

	// Platform and overwrite flag
	uploadOpts["platform"] = reactNativeOpts.Platform
	if reactNativeOpts.Overwrite {
		uploadOpts["overwrite"] = "true"
	}

	// Prepare upload files
	fileFields := map[string]server.FileField{
		"sourceMap": server.LocalFile(reactNativeOpts.SourceMap),
		"bundle":    server.LocalFile(reactNativeOpts.Bundle),
	}

	logger.Info(fmt.Sprintf("Uploading React Native source map for platform: %s", reactNativeOpts.Platform))

	// Perform upload request
	err := server.ProcessFileRequest(
		globalOptions.ApiKey,
		"/react-native-source-map",
		uploadOpts,
		fileFields,
		reactNativeOpts.SourceMap,
		globalOptions,
		logger,
	)
	if err != nil {
		return fmt.Errorf("failed to upload React Native source map for platform %q: %w", reactNativeOpts.Platform, err)
	}

	return nil
}
