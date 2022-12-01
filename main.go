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
		Port              int    `help:"Port number for the upload server" default:"443"`
		ApiKey            string `help:"Bugsnag integration API key for this application"`
		FailOnUploadError bool   `help:"Stops the upload when a mapping file fails to upload to Bugsnag successfully" default:false`
		Upload            struct {

			// shared options
			Overwrite bool `help:"Whether to overwrite any existing symbol file with a matching ID"`
			Timeout   int  `help:"Number of seconds to wait before failing an upload request" default:"300"`
			Retries   int  `help:"Number of retry attempts before failing an upload request" default:"0"`

			// required options
			All        upload.DiscoverAndUploadAny `cmd:"" help:"Upload any symbol/mapping files"`
			DartSymbol upload.DartSymbol           `cmd:"" help:"Process and upload symbol files for Flutter" name:"dart"`
		} `cmd:"" help:"Upload symbol/mapping files"`
		CreateBuild build.CreateBuild `cmd:"" help:"Create or update build info"`
	}

	// If running without any extra arguments, default to the --help flag
	// https://github.com/alecthomas/kong/issues/33#issuecomment-1207365879
	if len(os.Args) < 2 {
		os.Args = append(os.Args, "--help")
	}

	ctx := kong.Parse(&commands)

	// Check if we have an apiKey in the request
	if commands.ApiKey == "" {
		log.Error("no API key provided", 1)
	}

	// Build connection URI
	endpoint, err := utils.BuildEndpointUrl(commands.UploadAPIRootUrl, commands.Port)

	if err != nil {
		log.Error("Failed to build upload url: "+err.Error(), 1)
	}

	switch ctx.Command() {

	// Upload command
	case "upload all <path>":
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

	case "upload dart <path>":
		log.Info("Uploading files to: " + endpoint)

		err := upload.Dart(commands.Upload.DartSymbol.Path,
			commands.Upload.DartSymbol.AppVersion,
			commands.Upload.DartSymbol.AppVersionCode,
			commands.Upload.DartSymbol.AppBundleVersion,
			commands.Upload.DartSymbol.IosAppPath,
			endpoint+"/dart-symbol",
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.ApiKey,
			commands.FailOnUploadError)

		if err != nil {
			log.Error(err.Error(), 1)
		}

		log.Success("Upload(s) completed")

	case "create-build":
		log.Info("Creating build on: " + endpoint)
	default:
		println(ctx.Command())
	}
}
