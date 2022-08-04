package main

import (
	"github.com/alecthomas/kong"
	"os"
)

func main() {
	var commands struct {
		UploadServer string `help:"Bugsnag On-Premise upload server URL"`
		Upload       struct {
			// sub commands
			All            DiscoverAndUploadAny `cmd:"" help:"Find and upload any symbol files"`
			AndroidMapping AndroidMapping       `cmd:"" help:"Upload Android mapping file"`
			NdkLibrary     NdkLibrary           `cmd:"" help:"Upload Android NDK shared library symbols"`
			SourceMap      SourceMap            `cmd:"" help:"Upload JavaScript source map"`
			Dsym           AppleDebugSymbol     `cmd:"" help:"Upload dSYM files for Apple platforms"`
			DartSymbol     DartSymbol           `cmd:"" help:"Upload Dart symbol files"`
			BreakpadSymbol BreakpadSymbol       `cmd:"" help:"Upload Breakpad symbol files"`

			// shared options
			Overwrite bool   `help:"ignore existing upload with same version"`
			Timeout   int    `help:"seconds to wait before failing an upload request"`
			Retries   int    `help:"number of retry attempts before failing a request"`
			ApiKey    string `help:"Bugsnag project API key"`
		} `cmd:"" help:"Upload files"`
		CreateBuild CreateBuild `cmd:"" help:"Create or update build info"`
	}
	ctx := kong.Parse(&commands)

	// Check if we have an API key.
	// Look at moving this to its own utils file and make it more generic?
	if commands.Upload.ApiKey == "" {
		println("No API KEY...")
		println("Checking for ENV")
		if value, ok := os.LookupEnv("BUGSNAG_API_KEY"); ok{
			commands.Upload.ApiKey = value
			println("ENV found!")
		} else {
			println("No ENV for API key...")
			return
		}
	}

	switch ctx.Command() {
	case "uplod android-mapping <path>":
		println("mapping file!")
	case "upload ndk-library <path>":
		println("mapping file!")
	case "upload source-map <path>":
		println("source maps!")
	case "upload dsym <path>":
		println("Dsym!")
	case "upload dart-symbol <path>":
		println("Dart Symbol!")
		var uri = "https://upload.bugsnag.com/dart-symbol"
		DartUpload(uri, commands.Upload.ApiKey,commands.Upload.DartSymbol.BuildID  ,commands.Upload.DartSymbol.Path)
	case "upload breakpad-symbol <path>":
		println("BreakpadSymbol!")
	case "create-build":
		SendBuildInfo(ctx)
	default:
		println(ctx.Command())
	}
}
