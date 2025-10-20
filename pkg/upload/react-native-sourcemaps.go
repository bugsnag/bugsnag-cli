package upload

import (
	"fmt"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
)

func ProcessReactNativeSourcemaps(globalOptions options.CLI, logger log.Logger) error {
	var (
		uploadOptions = make(map[string]string)
		err           error
	)
	reactNativeOptions := globalOptions.Upload.ReactNativeSourcemap

	// Return early if all three of these are empty/undefined
	if reactNativeOptions.VersionName == "" && (reactNativeOptions.VersionCode == "" || reactNativeOptions.BundleVersion == "") && reactNativeOptions.CodeBundleId == "" {
		var platformSpecificTerminology string
		switch reactNativeOptions.Platform {
		case "android":
			platformSpecificTerminology = "version code"
		case "ios":
			platformSpecificTerminology = "bundle version"
		default:
			platformSpecificTerminology = "version code"
		}

		return fmt.Errorf("you must set at least the version name, %s and code bundle ID to uniquely identify the build", platformSpecificTerminology)
	}

	// If codeBundleId is set, use that instead of appVersion and versionCode
	if reactNativeOptions.CodeBundleId != "" {
		logger.Info("Using code bundle ID to identify build")
		uploadOptions["codeBundleId"] = reactNativeOptions.CodeBundleId
	} else {
		logger.Info("Using version name and version code to identify build")
		uploadOptions["appVersion"] = reactNativeOptions.VersionName

		if reactNativeOptions.VersionCode != "" {
			uploadOptions["appVersionCode"] = reactNativeOptions.VersionCode
		} else if reactNativeOptions.BundleVersion != "" {
			uploadOptions["appBundleVersion"] = reactNativeOptions.BundleVersion
		}
	}

	if reactNativeOptions.Dev {
		uploadOptions["dev"] = "true"
	}

	if reactNativeOptions.ProjectRoot == "" {
		uploadOptions["projectRoot"] = string(reactNativeOptions.Path)
	} else {
		uploadOptions["projectRoot"] = reactNativeOptions.ProjectRoot
	}

	uploadOptions["platform"] = reactNativeOptions.Platform

	if reactNativeOptions.Overwrite {
		uploadOptions["overwrite"] = "true"
	}

	fileFieldData := make(map[string]server.FileField)
	fileFieldData["sourceMap"] = server.LocalFile(reactNativeOptions.SourceMap)
	fileFieldData["bundle"] = server.LocalFile(reactNativeOptions.Bundle)

	err = server.ProcessFileRequest(globalOptions.ApiKey, "/react-native-source-map", uploadOptions, fileFieldData, reactNativeOptions.SourceMap, globalOptions, logger)

	if err != nil {
		return err
	}

	return nil
}
