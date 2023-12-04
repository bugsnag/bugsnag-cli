package main

import (
	"github.com/alecthomas/kong"
	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/bugsnag/bugsnag-cli/pkg/build"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
)

var package_version = "2.0.0"

func main() {
	commands := options.CLI{}

	// If running without any extra arguments, default to the --help flag
	// https://github.com/alecthomas/kong/issues/33#issuecomment-1207365879
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}

	ctx := kong.Parse(&commands,
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": package_version,
		})

	// Build connection URI
	endpoint, err := utils.BuildEndpointUrl(commands.UploadAPIRootUrl, commands.Port)

	if err != nil {
		log.Error("Failed to build upload url: "+err.Error(), 1)
	}

	if commands.DryRun {
		log.Info("Performing dry run - no data will be sent to BugSnag")
	}

	switch ctx.Command() {

	case "upload all <path>":

		if commands.ApiKey == "" {
			log.Error("no API key provided", 1)
		}

		err := upload.All(
			commands.Upload.All.Path,
			commands.Upload.All.UploadOptions,
			endpoint,
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.ApiKey,
			commands.FailOnUploadError,
			commands.DryRun,
		)

		if err != nil {
			log.Error(err.Error(), 1)
		}

	case "upload android-aab <path>":

		err := upload.ProcessAndroidAab(
			commands.ApiKey,
			commands.Upload.AndroidAab.ApplicationId,
			commands.Upload.AndroidAab.BuildUuid,
			commands.Upload.AndroidAab.Path,
			commands.Upload.AndroidAab.ProjectRoot,
			commands.Upload.AndroidAab.VersionCode,
			commands.Upload.AndroidAab.VersionName,
			endpoint,
			commands.FailOnUploadError,
			commands.Upload.Retries,
			commands.Upload.Timeout,
			commands.Upload.Overwrite,
			commands.DryRun,
		)

		if err != nil {
			log.Error(err.Error(), 1)
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
			commands.FailOnUploadError,
			commands.Upload.Retries,
			commands.Upload.Timeout,
			commands.Upload.Overwrite,
			commands.DryRun,
		)

		if err != nil {
			log.Error(err.Error(), 1)
		}

	case "upload android-proguard <path>", "upload android-proguard":

		err := upload.ProcessAndroidProguard(
			commands.ApiKey,
			commands.Upload.AndroidProguard.ApplicationId,
			commands.Upload.AndroidProguard.AppManifest,
			commands.Upload.AndroidProguard.BuildUuid,
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
		)

		if err != nil {
			log.Error(err.Error(), 1)
		}

	case "upload dart <path>":

		if commands.ApiKey == "" {
			log.Error("no API key provided", 1)
		}

		err := upload.Dart(commands.Upload.DartSymbol.Path,
			commands.Upload.DartSymbol.VersionName,
			commands.Upload.DartSymbol.VersionCode,
			commands.Upload.DartSymbol.BundleVersion,
			commands.Upload.DartSymbol.IosAppPath,
			endpoint,
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.ApiKey,
			commands.FailOnUploadError,
			commands.DryRun,
		)

		if err != nil {
			log.Error(err.Error(), 1)
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
		)

		if err != nil {
			log.Error(err.Error(), 1)
		}

	case "upload unity-android <path>":

		if commands.ApiKey == "" {
			log.Error("no API key provided", 1)
		}

		err := upload.ProcessUnityAndroid(
			commands.ApiKey,
			string(commands.Upload.UnityAndroid.AabPath),
			commands.Upload.UnityAndroid.ApplicationId,
			commands.Upload.UnityAndroid.VersionCode,
			commands.Upload.UnityAndroid.BuildUuid,
			commands.Upload.UnityAndroid.VersionName,
			commands.Upload.UnityAndroid.ProjectRoot,
			commands.Upload.UnityAndroid.Path,
			endpoint,
			commands.FailOnUploadError,
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.DryRun,
		)

		if err != nil {
			log.Error(err.Error(), 1)
		}

	case "create-build", "create-build <path>":
		var androidManifestPath string
		var BaseOptions build.CreateBuildInfo

		UserBuildOptions := build.PopulateFromCliOpts(commands)

		BaseOptions = build.PopulateFromPath(commands.CreateBuild.Path[0])

		if commands.CreateBuild.AndroidAab != "" {
			androidManifestPath, err = android.GetAndroidManifestFileFromAAB(string(commands.CreateBuild.AndroidAab))

			if err != nil {
				log.Error(err.Error(), 1)
			}
		}

		if androidManifestPath == "" {
			androidManifestPath = string(commands.CreateBuild.AppManifest)
		}

		if androidManifestPath != "" {
			ManifestBuildOptions := build.PopulateFromAndroidManifest(androidManifestPath)
			BaseOptions = build.TheGreatMerge(BaseOptions, ManifestBuildOptions)
		}

		FinalMerge := build.TheGreatMerge(UserBuildOptions, BaseOptions)

		// Build connection URI
		endpoint, err := utils.BuildEndpointUrl(commands.BuildApiRootUrl, commands.Port)

		if err != nil {
			log.Error("Failed to build upload url: "+err.Error(), 1)
		}

		// Validate the required options for the API
		err = FinalMerge.Validate()

		if err != nil {
			log.Error(err.Error(), 1)
		}

		log.Info("Creating build on: " + endpoint)

		err = build.ProcessBuildRequest(FinalMerge, endpoint, commands.DryRun)

		if err != nil {
			log.Error(err.Error(), 1)
		}

		log.Success("Build created")

	case "create-android-build-id", "create-android-build-id <path>":
		err := build.PrintAndroidBuildId(commands.CreateAndroidBuildId.Path)

		if err != nil {
			log.Error(err.Error(), 1)
		}

	default:
		println(ctx.Command())
	}
}
