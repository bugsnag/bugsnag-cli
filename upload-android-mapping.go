package main

type AndroidMapping struct {
	Path []string `arg:"" name:"path" help:"Path to module or file upload" type:"path"`

	AppManifestPath string `help:"Path to the merged app manifest"`
	ApplicationID string `help:"Module application identifier"`
	BuildUUID string `help:"Module Build UUID"`
	Configuration string `help:"Build type, like 'debug' or 'release'"`
	VersionCode string `help:"Module version code"`
	VersionName string `help:"Module version name"`
}
