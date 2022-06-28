package main

import "github.com/alecthomas/kong"

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

	switch ctx.Command() {
	case "android-mapping <path>":
		println("mapping file!")
	case "ndk-library <path>":
		println("ndk library!")
	case "source-map <path>":
		println("source maps!")
	case "create-build":
		SendBuildInfo(ctx)
	default:
		println(ctx.Command())
	}
}
