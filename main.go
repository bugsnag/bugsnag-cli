package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"

	"github.com/bugsnag/bugsnag-cli/pkg/build"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

var package_version = "2.4.1"

func main() {
	commands := options.CLI{}

	// If running without any extra arguments, default to the --help flag
	// https://github.com/alecthomas/kong/issues/33#issuecomment-1207365879
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}

	kongCtx := kong.Parse(&commands,
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": package_version,
		})

	if commands.Verbose {
		commands.LogLevel = "debug"
	}

	logger := log.NewLoggerWrapper(commands.LogLevel)

	// Build connection URI
	endpoint, err := utils.BuildEndpointUrl(commands.Upload.UploadAPIRootUrl, commands.Port)

	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to build upload url: %s", err.Error()))
	}

	if commands.DryRun {
		logger.Info("Performing dry run - no data will be sent to BugSnag")
	}

	if commands.FailOnUploadError {
		logger.Warn("The `--fail-on-upload-error` flag is deprecated and will be removed in a future release. All commands now fail if the upload is unsuccessful.")
	}

	switch kongCtx.Command() {

	case "upload all <path>":

		if commands.ApiKey == "" {
			logger.Fatal("missing api key, please specify using `--api-key`")
		}

		err := upload.All(
			commands.Upload.All.Path,
			commands.Upload.All.UploadOptions,
			endpoint,
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.ApiKey,
			commands.DryRun,
			logger,
		)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload android-aab <path>", "upload android-aab":

		err := upload.ProcessAndroidAab(
			commands.ApiKey,
			commands.Upload.AndroidAab.ApplicationId,
			commands.Upload.AndroidAab.BuildUuid,
			commands.Upload.AndroidAab.NoBuildUuid,
			commands.Upload.AndroidAab.Path,
			commands.Upload.AndroidAab.ProjectRoot,
			commands.Upload.AndroidAab.VersionCode,
			commands.Upload.AndroidAab.VersionName,
			endpoint,
			commands.Upload.Retries,
			commands.Upload.Timeout,
			commands.Upload.Overwrite,
			commands.DryRun,
			logger,
		)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload android-ndk <path>", "upload android-ndk":

		err := upload.ProcessAndroidNDK(
			commands.ApiKey,
			commands.Upload.AndroidNdk.ApplicationId,
			commands.Upload.AndroidNdk.AndroidNdkRoot,
			commands.Upload.AndroidNdk.AppManifest,
			commands.Upload.AndroidNdk.Path,
			commands.Upload.AndroidNdk.ProjectRoot,
			commands.Upload.AndroidNdk.Variant,
			commands.Upload.AndroidNdk.VersionCode,
			commands.Upload.AndroidNdk.VersionName,
			endpoint,
			commands.Upload.Retries,
			commands.Upload.Timeout,
			commands.Upload.Overwrite,
			commands.DryRun,
			logger,
		)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload android-proguard <path>", "upload android-proguard":

		err := upload.ProcessAndroidProguard(
			commands.ApiKey,
			commands.Upload.AndroidProguard.ApplicationId,
			commands.Upload.AndroidProguard.AppManifest,
			commands.Upload.AndroidProguard.BuildUuid,
			commands.Upload.AndroidProguard.NoBuildUuid,
			commands.Upload.AndroidProguard.DexFiles,
			commands.Upload.AndroidProguard.Path,
			commands.Upload.AndroidProguard.Variant,
			commands.Upload.AndroidProguard.VersionCode,
			commands.Upload.AndroidProguard.VersionName,
			endpoint,
			commands.Upload.Retries,
			commands.Upload.Timeout,
			commands.Upload.Overwrite,
			commands.DryRun,
			logger,
		)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload dart <path>":

		if commands.ApiKey == "" {
			logger.Fatal("missing api key, please specify using `--api-key`")
		}

		err := upload.Dart(
			commands.Upload.DartSymbol.Path,
			commands.Upload.DartSymbol.VersionName,
			commands.Upload.DartSymbol.VersionCode,
			commands.Upload.DartSymbol.BundleVersion,
			string(commands.Upload.DartSymbol.IosAppPath),
			endpoint,
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.ApiKey,
			commands.DryRun,
			logger,
		)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload react-native-android", "upload react-native-android <path>":

		err := upload.ProcessReactNativeAndroid(
			commands.ApiKey,
			commands.Upload.ReactNativeAndroid.AppManifest,
			commands.Upload.ReactNativeAndroid.Bundle,
			commands.Upload.ReactNativeAndroid.CodeBundleId,
			commands.Upload.ReactNativeAndroid.Dev,
			commands.Upload.ReactNativeAndroid.Path,
			commands.Upload.ReactNativeAndroid.ProjectRoot,
			commands.Upload.ReactNativeAndroid.Variant,
			commands.Upload.ReactNativeAndroid.VersionName,
			commands.Upload.ReactNativeAndroid.VersionCode,
			commands.Upload.ReactNativeAndroid.SourceMap,
			endpoint,
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.DryRun,
			logger,
		)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload react-native-ios", "upload react-native-ios <path>":

		err := upload.ProcessReactNativeIos(
			commands.ApiKey,
			commands.Upload.ReactNativeIos.VersionName,
			commands.Upload.ReactNativeIos.BundleVersion,
			commands.Upload.ReactNativeIos.Scheme,
			commands.Upload.ReactNativeIos.SourceMap,
			commands.Upload.ReactNativeIos.Bundle,
			commands.Upload.ReactNativeIos.Plist,
			commands.Upload.ReactNativeIos.XcodeProject,
			commands.Upload.ReactNativeIos.CodeBundleID,
			commands.Upload.ReactNativeIos.Dev,
			commands.Upload.ReactNativeIos.ProjectRoot,
			commands.Upload.ReactNativeIos.Path,
			endpoint,
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.DryRun,
			logger,
		)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload dsym", "upload dsym <path>":

		err := upload.ProcessDsym(
			commands.ApiKey,
			commands.Upload.Dsym.Scheme,
			string(commands.Upload.Dsym.XcodeProject),
			string(commands.Upload.Dsym.Plist),
			commands.Upload.Dsym.ProjectRoot,
			commands.Upload.Dsym.IgnoreMissingDwarf,
			commands.Upload.Dsym.IgnoreEmptyDsym,
			commands.Upload.Dsym.Path,
			endpoint,
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.DryRun,
			logger,
		)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload unity-android <path>":

		if commands.ApiKey == "" {
			logger.Fatal("missing api key, please specify using `--api-key`")
		}

		err := upload.ProcessUnityAndroid(
			commands.ApiKey,
			string(commands.Upload.UnityAndroid.AabPath),
			commands.Upload.UnityAndroid.ApplicationId,
			commands.Upload.UnityAndroid.VersionCode,
			commands.Upload.UnityAndroid.BuildUuid,
			commands.Upload.UnityAndroid.NoBuildUuid,
			commands.Upload.UnityAndroid.VersionName,
			commands.Upload.UnityAndroid.ProjectRoot,
			commands.Upload.UnityAndroid.Path,
			endpoint,
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.DryRun,
			logger,
		)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "create-build", "create-build <path>":
		// Create Build Info
		CreateBuildOptions, err := build.GatherBuildInfo(commands)

		if err != nil {
			logger.Fatal(err.Error())
		}

		// Validate Build Info
		err = CreateBuildOptions.Validate()

		if err != nil {
			logger.Fatal(err.Error())
		}

		// Get Endpoint URL
		endpoint, err = utils.BuildEndpointUrl(commands.CreateBuild.BuildApiRootUrl, commands.Port)

		if err != nil {
			logger.Fatal(fmt.Sprintf("Failed to build upload url: %s", err.Error()))
		}

		err = build.ProcessCreateBuild(CreateBuildOptions, endpoint, commands.DryRun, commands.CreateBuild.Timeout, commands.CreateBuild.Retries, logger)

		if err != nil {
			logger.Fatal(err.Error())
		}

		logger.Info("Build created")

	case "create-android-build-id", "create-android-build-id <path>":
		err := build.PrintAndroidBuildId(commands.CreateAndroidBuildId.Path)

		if err != nil {
			logger.Fatal(err.Error())
		}

	default:
		println(kongCtx.Command())
	}
}
