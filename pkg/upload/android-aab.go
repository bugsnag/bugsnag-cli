package upload

import (
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type AndroidAabMapping struct {
	ApplicationId string            `help:"Module application identifier"`
	BuildUuid     string            `help:"Module Build UUID"`
	Configuration string            `help:"Build type, like 'debug' or 'release'"`
	Path          utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	ProjectRoot   string            `help:"path to remove from the beginning of the filenames in the mapping file"`
	VersionCode   string            `help:"Module version code"`
	VersionName   string            `help:"Module version name"`
}
