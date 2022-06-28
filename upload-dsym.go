package main

type AppleDebugSymbol struct {
	Path []string `arg:"" name:"path" help:"Path to directory or file to upload" type:"path"`

	InfoPlistPath string `help:"Path to the app Info.plist file"`
	ProjectRoot string `help:"Directory where the app was built"`
}
