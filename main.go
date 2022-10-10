package main

import (
	"github.com/alecthomas/kong"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
)

func main() {
	var commands struct {
		UploadAPIRootUrl string `help:"Bugsnag On-Premise upload server URL. Can contain port number" default:"https://upload.bugsnag.com"`
		Port		 int	`help:"Port number for the upload server" default:"443"`
		ApiKey       string `help:"Bugsnag project API key"`
		Upload       struct {

			// shared options
			Overwrite     bool              `help:"ignore existing upload with same version"`
			Timeout       int               `help:"seconds to wait before failing an upload request" default:"300"`
			Retries       int               `help:"number of retry attempts before failing a request" default:"0"`

			// required options
			All            upload.DiscoverAndUploadAny `cmd:"" help:"Find and upload any symbol files"`
			DartSymbol     upload.DartSymbol           `cmd:"" help:"Upload Dart symbol files" name:"dart"`
		} `cmd:"" help:"Upload files"`
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
	endpoint := utils.BuildEndpointUrl(commands.UploadAPIRootUrl, commands.Port)

	log.Info("Uploading files to: " + endpoint)

	switch ctx.Command() {

	// Upload command
	case "upload all <path>":
		err := upload.All(commands.Upload.All.Path, commands.Upload.All.UploadOptions, endpoint, commands.Upload.Timeout,
			commands.Upload.Retries, commands.Upload.Overwrite, commands.ApiKey)

		if err != nil {
			log.Error(err.Error(), 1)
		}
		
	case "upload dart <path>":
		err := upload.Dart(commands.Upload.DartSymbol.Path,
			commands.Upload.DartSymbol.AppVersion,
			commands.Upload.DartSymbol.AppVersionCode,
			commands.Upload.DartSymbol.AppBundleVersion,
			commands.Upload.DartSymbol.IosAppPath,
			endpoint + "/dart-symbol",
			commands.Upload.Timeout,
			commands.Upload.Retries,
			commands.Upload.Overwrite,
			commands.ApiKey)

		if err != nil {
			log.Error(err.Error(), 1)
		}
	default:
		println(ctx.Command())
	}

	log.Success("Upload(s) completed")

}
