package main

import (
	"github.com/alecthomas/kong"
	"github.com/bugsnag/bugsnag-cli/pkg/build"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
)

func main() {
	var commands struct {
		UploadAPIRootUrl  string `help:"Bugsnag On-Premise upload server URL. Can contain port number" default:"https://upload.bugsnag.com"`
		BuildApiRootUrl   string `help:"Bugsnag On-Premise build server URL. Can contain port number" default:"https://build.bugsnag.com"`
		Port              int    `help:"Port number for the upload server" default:"443"`
		ApiKey            string `help:"(required) Bugsnag integration API key for this application"`
		FailOnUploadError bool   `help:"Stops the upload when a mapping file fails to upload to Bugsnag successfully" default:"false"`
		AppVersion        string `help:"The version of the application."`
		AppVersionCode    string `help:"The version code for the application (Android only)."`
		AppBundleVersion  string `help:"The bundle version for the application (iOS only)."`
		Upload            struct {

			// shared options
			Overwrite bool `help:"Whether to overwrite any existing symbol file with a matching ID"`
			Timeout   int  `help:"Number of seconds to wait before failing an upload request" default:"300"`
			Retries   int  `help:"Number of retry attempts before failing an upload request" default:"0"`
			DryRun    bool `help:"Validate but do not upload"`

			// required options
			AndroidAab         upload.AndroidAabMapping      `cmd:"" help:"Process and upload application bundle files for Android"`
			All                upload.DiscoverAndUploadAny   `cmd:"" help:"Upload any symbol/mapping files"`
			AndroidNdk         upload.AndroidNdkMapping      `cmd:"" help:"Process and upload Proguard mapping files for Android"`
			AndroidProguard    upload.AndroidProguardMapping `cmd:"" help:"Process and upload NDK symbol files for Android"`
			DartSymbol         upload.DartSymbol             `cmd:"" help:"Process and upload symbol files for Flutter" name:"dart"`
			ReactNativeAndroid upload.ReactNativeAndroid     `cmd:"" help:"Upload source maps for React Native Android"`
		} `cmd:"" help:"Upload symbol/mapping files"`
		CreateBuild build.CreateBuild `cmd:"" help:"Provide extra information whenever you build, release, or deploy your application"`
	}

	// If running without any extra arguments, default to the --help flag
	// https://github.com/alecthomas/kong/issues/33#issuecomment-1207365879
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}

	ctx := kong.Parse(&commands)

	// Build connection URI
	endpoint, err := utils.BuildEndpointUrl(commands.UploadAPIRootUrl, commands.Port)

	if err != nil {
		log.Error("Failed to build upload url: "+err.Error(), 1)
	}

	switch ctx.Command() {

	case "upload all <path>":

		if commands.ApiKey == "" {
			log.Error("no API key provided", 1)
		}

		log.Info("Uploading files to: " + endpoint)

		err := upload.All(
			commands.Upload.All.Path,
			commands.Upload.All.UploadOptions,
			endpoint,
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.ApiKey,
			commands.FailOnUploadError)

		if err != nil {
			log.Error(err.Error(), 1)
		}

		log.Success("Upload(s) completed")

	case "upload android-aab <path>":

		log.Info("Uploading files to: " + endpoint)

		err := upload.ProcessAndroidAab(
			commands.ApiKey,
			commands.Upload.AndroidAab.AndroidNdkRoot,
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
			commands.Upload.DryRun,
		)

		if err != nil {
			log.Error(err.Error(), 1)
		}

		log.Success("Upload(s) completed")

	case "upload android-ndk <path>", "upload android-ndk":

		endpoint = endpoint + "/ndk-symbol"

		log.Info("Uploading files to: " + endpoint)

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
			commands.Upload.DryRun,
		)

		if err != nil {
			log.Error(err.Error(), 1)
		}

		log.Success("Upload(s) completed")

	case "upload android-proguard <path>", "upload android-proguard":

		if commands.ApiKey == "" {
			log.Error("no API key provided", 1)
		}

		log.Info("Uploading files to: " + endpoint)

		err := upload.ProcessAndroidProguard(
			commands.ApiKey,
			commands.Upload.AndroidProguard.ApplicationId,
			commands.Upload.AndroidProguard.AppManifest,
			commands.Upload.AndroidProguard.BuildUuid,
			commands.Upload.AndroidProguard.Path,
			commands.Upload.AndroidProguard.Variant,
			commands.Upload.AndroidProguard.VersionCode,
			commands.Upload.AndroidProguard.VersionName,
			endpoint,
			commands.Upload.Retries,
			commands.Upload.Timeout,
			commands.Upload.Overwrite,
			commands.Upload.DryRun,
		)

		if err != nil {
			log.Error(err.Error(), 1)
		}

		log.Success("Upload(s) completed")

	case "upload dart <path>":

		if commands.ApiKey == "" {
			log.Error("no API key provided", 1)
		}

		endpoint = endpoint + "/dart-symbol"

		log.Info("Uploading files to: " + endpoint)

		err := upload.Dart(commands.Upload.DartSymbol.Path,
			commands.AppVersion,
			commands.AppVersionCode,
			commands.AppBundleVersion,
			commands.Upload.DartSymbol.IosAppPath,
			endpoint,
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.ApiKey,
			commands.FailOnUploadError)

		if err != nil {
			log.Error(err.Error(), 1)
		}

		log.Success("Upload(s) completed")

	case "upload react-native-android", "upload react-native-android <path>":

		if commands.ApiKey == "" {
			log.Error("no API key provided", 1)
		}

		// Build endpoint URI
		endpoint = endpoint + "/react-native-source-map"

		err := upload.ProcessReactNativeAndroid(commands.Upload.ReactNativeAndroid.Path,
			commands.Upload.ReactNativeAndroid.AppManifestPath,
			commands.AppVersion,
			commands.AppVersionCode,
			commands.Upload.ReactNativeAndroid.CodeBundleId,
			commands.Upload.ReactNativeAndroid.Dev,
			commands.Upload.ReactNativeAndroid.SourceMapPath,
			commands.Upload.ReactNativeAndroid.BundlePath,
			commands.Upload.ReactNativeAndroid.ProjectRoot,
			endpoint,
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.ApiKey)

		if err != nil {
			log.Error(err.Error(), 1)
		}

		log.Success("Upload(s) completed")

	case "create-build":

		if commands.ApiKey == "" {
			log.Error("no API key provided", 1)
		}

		// Build connection URI
		endpoint, err := utils.BuildEndpointUrl(commands.BuildApiRootUrl, commands.Port)

		if err != nil {
			log.Error("Failed to build upload url: "+err.Error(), 1)
		}

		log.Info("Creating build on: " + endpoint)

		buildUploadError := build.ProcessBuildRequest(commands.ApiKey,
			commands.CreateBuild.BuilderName,
			commands.CreateBuild.ReleaseStage,
			commands.CreateBuild.Provider,
			commands.CreateBuild.Repository,
			commands.CreateBuild.Revision,
			commands.AppVersion,
			commands.AppVersionCode,
			commands.AppBundleVersion,
			commands.CreateBuild.Metadata,
			endpoint)
		if buildUploadError != nil {
			log.Error(buildUploadError.Error(), 1)
		}

		log.Success("Build created")
	default:
		println(ctx.Command())
	}
}
