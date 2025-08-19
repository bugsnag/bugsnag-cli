package options

import "github.com/bugsnag/bugsnag-cli/pkg/utils"

type LinuxOptions struct {
	Path          utils.Paths `arg:"" name:"path" help:"The path to the directory or file to upload" type:"path" default:"."`
	ApplicationId string      `help:"A unique application ID, usually the package name, of the application"`
	ProjectRoot   string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path" default:"."`
	Variant       string      `help:"The build type/flavor (e.g. Debug, Release) used to disambiguate the between built files when searching the project directory" default:"release"`
	VersionCode   string      `help:"The version code of this build of the application"`
	VersionName   string      `help:"The version of the application"`
	BuildFolder   utils.Path  `help:"The path to the build folder" type:"path"`
}
