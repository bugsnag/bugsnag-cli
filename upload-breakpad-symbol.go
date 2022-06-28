package main

type BreakpadSymbol struct {
	Path []string `arg:"" name:"path" help:"Path to directory or file to upload" type:"path"`

	AppVersion string `help:"Module version"`
	ProjectRoot string `help:"Directory where the app was built"`
}
