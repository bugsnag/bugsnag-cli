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

var package_version = "2.9.1"

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

		err := upload.All(commands, endpoint, logger)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload android-aab <path>", "upload android-aab":

		err := upload.ProcessAndroidAab(commands, endpoint, logger)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload android-ndk <path>", "upload android-ndk":

		err := upload.ProcessAndroidNDK(commands, endpoint, logger)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload android-proguard <path>", "upload android-proguard":

		err := upload.ProcessAndroidProguard(commands, endpoint, logger)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload dart <path>":

		if commands.ApiKey == "" {
			logger.Fatal("missing api key, please specify using `--api-key`")
		}

		err := upload.Dart(commands, endpoint, logger)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload react-native", "upload react-native <path>":

		err := upload.ProcessReactNative(commands, endpoint, logger)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload react-native-android", "upload react-native-android <path>":

		err := upload.ProcessReactNativeAndroid(commands, endpoint, logger)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload react-native-ios", "upload react-native-ios <path>":

		err := upload.ProcessReactNativeIos(commands, endpoint, logger)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload js", "upload js <path>":

		err := upload.ProcessJs(commands, endpoint, logger)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload xcode-build", "upload xcode-build <path>":

		err := upload.ProcessXcodeBuild(commands, endpoint, logger)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload xcode-archive", "upload xcode-archive <path>":

		err := upload.ProcessXcodeArchive(commands, endpoint, logger)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload dsym", "upload dsym <path>":

		logger.Warn("The `upload dsym` command is deprecated and will be removed in a future release. Please use `upload xcode-build` instead.")

		err := upload.ProcessXcodeBuild(commands, endpoint, logger)

		if err != nil {
			logger.Fatal(err.Error())
		}

	case "upload unity-android <path>":

		if commands.ApiKey == "" {
			logger.Fatal("missing api key, please specify using `--api-key`")
		}

		err := upload.ProcessUnityAndroid(commands, endpoint, logger)

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

		err = build.ProcessCreateBuild(CreateBuildOptions, endpoint, commands, logger)

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
