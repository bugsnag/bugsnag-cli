package upload

import (
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

func ProcessReactNative(globalOptions options.CLI, endpoint string, logger log.Logger) error {
	reactNativeOptions := globalOptions.Upload.ReactNative

	// The commands must be run from either the android/ or ios/ subdirectory.
	androidPath := []string{}
	iosPath := []string{}
	for _, basePath := range reactNativeOptions.Path {
		androidPath = append(androidPath, filepath.Join(basePath, "android"))
		iosPath = append(iosPath, filepath.Join(basePath, "ios"))
	}

	logger.Info("Uploading JavaScript source maps for Android")
	globalOptions.Upload.ReactNativeAndroid = options.ReactNativeAndroid{
		Path:        androidPath,
		ProjectRoot: reactNativeOptions.ProjectRoot,
		ReactNative: reactNativeOptions.Shared,
		Android:     reactNativeOptions.AndroidSpecific,
	}
	if err := ProcessReactNativeAndroid(globalOptions, endpoint, logger); err != nil {
		return err
	}

	logger.Info("Uploading JavaScript source maps for iOS")
	globalOptions.Upload.ReactNativeIos = options.ReactNativeIos{
		Path:        iosPath,
		ProjectRoot: reactNativeOptions.ProjectRoot,
		ReactNative: reactNativeOptions.Shared,
		Ios:         reactNativeOptions.IosSpecific,
	}
	if err := ProcessReactNativeIos(globalOptions, endpoint, logger); err != nil {
		return err
	}

	logger.Info("Uploading Android Proguard mappings")
	// Missing: ApplicationId BuildUuid NoBuildUuid DexFiles
	globalOptions.Upload.AndroidProguard = options.AndroidProguardMapping{
		Path:        androidPath,
		VersionName: reactNativeOptions.Shared.VersionName,
		AppManifest: reactNativeOptions.AndroidSpecific.AppManifest,
		Variant:     reactNativeOptions.AndroidSpecific.Variant,
		VersionCode: reactNativeOptions.AndroidSpecific.VersionCode,
	}
	if err := ProcessAndroidProguard(globalOptions, endpoint, logger); err != nil {
		return err
	}

	logger.Info("Uploading iOS dSYMs")
	// Missing: IgnoreEmptyDsym IgnoreMissingDwarf
	globalOptions.Upload.Dsym = options.Dsym{
		Path:         iosPath,
		VersionName:  reactNativeOptions.Shared.VersionName,
		ProjectRoot:  reactNativeOptions.ProjectRoot,
		Plist:        utils.Path(reactNativeOptions.IosSpecific.Plist),
		Scheme:       reactNativeOptions.IosSpecific.Scheme,
		XcodeProject: utils.Path(reactNativeOptions.IosSpecific.XcodeProject),
	}
	if err := ProcessDsym(globalOptions, endpoint, logger); err != nil {
		return err
	}

	logger.Info("Uploading Android NDK symbols")
	// Missing: ApplicationId AndroidNdkRoot
	globalOptions.Upload.AndroidNdk = options.AndroidNdkMapping{
		Path:        androidPath,
		VersionName: reactNativeOptions.Shared.VersionName,
		ProjectRoot: reactNativeOptions.ProjectRoot,
		AppManifest: reactNativeOptions.AndroidSpecific.AppManifest,
		Variant:     reactNativeOptions.AndroidSpecific.Variant,
		VersionCode: reactNativeOptions.AndroidSpecific.VersionCode,
	}
	if err := ProcessAndroidNDK(globalOptions, endpoint, logger); err != nil {
		return err
	}

	return nil
}
