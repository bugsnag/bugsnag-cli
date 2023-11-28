package main

import (
	"os"

	"github.com/alecthomas/kong"

	"github.com/bugsnag/bugsnag-cli/pkg/build"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

var package_version = "2.0.0"

// Global CLI options
type Globals struct {
	UploadAPIRootUrl  string            `help:"Bugsnag On-Premise upload server URL. Can contain port number" default:"https://upload.bugsnag.com"`
	BuildApiRootUrl   string            `help:"Bugsnag On-Premise build server URL. Can contain port number" default:"https://build.bugsnag.com"`
	Port              int               `help:"Port number for the upload server" default:"443"`
	ApiKey            string            `help:"(required) Bugsnag integration API key for this application"`
	FailOnUploadError bool              `help:"Stops the upload when a mapping file fails to upload to Bugsnag successfully" default:"false"`
	Version           utils.VersionFlag `name:"version" help:"Print version information and quit"`
	DryRun            bool              `help:"Validate but do not process"`
}

// Unique CLI options
type CLI struct {
	Globals

	Upload struct {
		// shared options
		Overwrite bool `help:"Whether to overwrite any existing symbol file with a matching ID"`
		Timeout   int  `help:"Number of seconds to wait before failing an upload request" default:"300"`
		Retries   int  `help:"Number of retry attempts before failing an upload request" default:"0"`

		// required options
		AndroidAab         upload.AndroidAabMapping      `cmd:"" help:"Process and upload application bundle files for Android"`
		All                upload.DiscoverAndUploadAny   `cmd:"" help:"Upload any symbol/mapping files"`
		AndroidNdk         upload.AndroidNdkMapping      `cmd:"" help:"Process and upload Proguard mapping files for Android"`
		AndroidProguard    upload.AndroidProguardMapping `cmd:"" help:"Process and upload NDK symbol files for Android"`
		DartSymbol         upload.DartSymbolOptions      `cmd:"" help:"Process and upload symbol files for Flutter" name:"dart"`
		ReactNativeAndroid upload.ReactNativeAndroid     `cmd:"" help:"Upload source maps for React Native Android"`
		ReactNativeIos     upload.ReactNativeIos         `cmd:"" help:"Upload source maps for React Native iOS"`
		UnityAndroid       upload.UnityAndroid           `cmd:"" help:"Upload Android mappings and NDK symbol files from Unity projects"`
	} `cmd:"" help:"Upload symbol/mapping files"`
	CreateBuild          build.CreateBuild          `cmd:"" help:"Provide extra information whenever you build, release, or deploy your application"`
	CreateAndroidBuildId build.CreateAndroidBuildId `cmd:"" help:"Generate a reproducible Build ID from .dex files"`
}

func main() {
	commands := CLI{}

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

	case "upload react-native-ios", "upload react-native-ios <path>":

		err := upload.ProcessReactNativeIos(
			commands.ApiKey,
			commands.Upload.ReactNativeIos.VersionName,
			commands.Upload.ReactNativeIos.BundleVersion,
			commands.Upload.ReactNativeIos.Scheme,
			commands.Upload.ReactNativeIos.SourceMap,
			commands.Upload.ReactNativeIos.Bundle,
			commands.Upload.ReactNativeIos.Plist,
			commands.Upload.ReactNativeIos.Xcworkspace,
			commands.Upload.ReactNativeIos.CodeBundleID,
			commands.Upload.ReactNativeIos.Dev,
			commands.Upload.ReactNativeIos.ProjectRoot,
			commands.Upload.ReactNativeIos.Path,
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
			commands.CreateBuild.VersionName,
			commands.CreateBuild.VersionCode,
			commands.CreateBuild.BundleVersion,
			commands.CreateBuild.Metadata,
			commands.CreateBuild.Path,
			endpoint,
			commands.DryRun,
		)

		if buildUploadError != nil {
			log.Error(buildUploadError.Error(), 1)
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
