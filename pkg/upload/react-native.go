package upload

import (
	"fmt"
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// ProcessReactNative handles the upload process for React Native projects.
//
// It processes both Android and iOS assets, including JavaScript source maps,
// Proguard mappings, dSYMs, and NDK symbols.
// It constructs the necessary paths based on the provided options and uploads
// the assets to the Bugsnag server.
//
// Parameters:
//   - globalOptions: CLI options containing React Native upload settings.
//   - logger: Logger instance for logging progress and errors.
//
// Returns:
//   - error: non-nil if any step fails during processing or uploading.
func ProcessReactNative(globalOptions options.CLI, logger log.Logger) error {
	reactNativeOptions := globalOptions.Upload.ReactNative

	// Construct Android and iOS paths
	androidPath, iosPath := generatePaths(reactNativeOptions.Path, "android", "ios")

	logger.Info("Starting upload of React Native assets")

	// Process React Native Android
	logger.Info("Uploading JavaScript source maps for Android")
	globalOptions.Upload.ReactNativeAndroid = options.ReactNativeAndroid{
		Path:        androidPath,
		ProjectRoot: reactNativeOptions.ProjectRoot,
		ReactNative: reactNativeOptions.Shared,
		Android:     reactNativeOptions.AndroidSpecific,
	}
	if err := ProcessReactNativeAndroid(globalOptions, logger); err != nil {
		return fmt.Errorf("failed to upload JavaScript source maps for Android: %w", err)
	}

	// Process React Native iOS
	logger.Info("Uploading JavaScript source maps for iOS")
	globalOptions.Upload.ReactNativeIos = options.ReactNativeIos{
		Path:        iosPath,
		ProjectRoot: reactNativeOptions.ProjectRoot,
		ReactNative: reactNativeOptions.Shared,
		Ios:         reactNativeOptions.IosSpecific,
	}
	if err := ProcessReactNativeIos(globalOptions, logger); err != nil {
		return fmt.Errorf("failed to upload JavaScript source maps for iOS: %w", err)
	}

	// Process Android Proguard mappings
	logger.Info("Uploading Android Proguard mappings")
	globalOptions.Upload.AndroidProguard = options.AndroidProguardMapping{
		Path:        androidPath,
		VersionName: reactNativeOptions.Shared.VersionName,
		AppManifest: reactNativeOptions.AndroidSpecific.AppManifest,
		Variant:     reactNativeOptions.AndroidSpecific.Variant,
		VersionCode: reactNativeOptions.AndroidSpecific.VersionCode,
	}
	if err := ProcessAndroidProguard(globalOptions, logger); err != nil {
		return fmt.Errorf("failed to upload Android Proguard mappings: %w", err)
	}

	// Process iOS dSYMs
	logger.Info("Uploading iOS dSYMs")
	globalOptions.Upload.Dsym = options.Dsym{
		Path: iosPath,
		Shared: options.DsymShared{
			ProjectRoot:  reactNativeOptions.ProjectRoot,
			Scheme:       reactNativeOptions.IosSpecific.Scheme,
			Plist:        utils.Path(reactNativeOptions.IosSpecific.Plist),
			XcodeProject: utils.Path(reactNativeOptions.IosSpecific.XcodeProject),
		},
	}
	if err := ProcessDsym(globalOptions, logger); err != nil {
		return fmt.Errorf("failed to upload iOS dSYMs: %w", err)
	}

	// Process Android NDK symbols
	logger.Info("Uploading Android NDK symbols")
	globalOptions.Upload.AndroidNdk = options.AndroidNdkMapping{
		Path:        androidPath,
		VersionName: reactNativeOptions.Shared.VersionName,
		ProjectRoot: reactNativeOptions.ProjectRoot,
		AppManifest: reactNativeOptions.AndroidSpecific.AppManifest,
		Variant:     reactNativeOptions.AndroidSpecific.Variant,
		VersionCode: reactNativeOptions.AndroidSpecific.VersionCode,
	}
	if err := ProcessAndroidNDK(globalOptions, logger); err != nil {
		return fmt.Errorf("failed to upload Android NDK symbols: %w", err)
	}

	logger.Info("Successfully uploaded all React Native assets")
	return nil
}

// generatePaths constructs platform-specific paths based on the base paths provided.
func generatePaths(basePaths []string, androidSubPath, iosSubPath string) ([]string, []string) {
	var androidPaths, iosPaths []string
	for _, basePath := range basePaths {
		androidPaths = append(androidPaths, filepath.Join(basePath, androidSubPath))
		iosPaths = append(iosPaths, filepath.Join(basePath, iosSubPath))
	}
	return androidPaths, iosPaths
}
