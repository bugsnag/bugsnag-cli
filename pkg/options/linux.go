package options

import "github.com/bugsnag/bugsnag-cli/pkg/utils"

type LinuxOptions struct {
	Path          utils.Paths `arg:"" name:"path" help:"The path to the directory or file to upload" type:"path" default:"."`
	ApplicationId string      `help:"A unique application ID, usually the package name, of the application"`
	ProjectRoot   string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path" default:"."`
	VersionName   string      `help:"The version of the application"`
	Overwrite     bool        `help:"Whether to ignore and overwrite existing uploads with same identifier, rather than failing if a matching file exists"`
}
