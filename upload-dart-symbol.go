package main

type DartSymbol struct {
	Path []string `arg:"" name:"path" help:"Path to directory or file to upload" type:"path"`

	AppVersion string `help:"Application version"`
	AppVersionCode string `help:"Module version code (Android only)"`
	AppBundleVersion string `help:"App bundle version (Apple platforms only)"`
	Platform string `help:"platform the program was built for"`
}
