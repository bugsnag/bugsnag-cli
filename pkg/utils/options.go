package utils

type AndroidOptions struct {
	BuildUuid     string `help:"Module Build UUID"`
	Configuration string `help:"Build type, like 'debug' or 'release'"`
	VersionCode   string `help:"Module version code"`
	VersionName   string `help:"Module version name"`
}
