package upload

import (
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type AndroidProguardMapping struct {
	ApplicationId   string            `help:"Module application identifier"`
	AppManifestPath utils.UploadPaths `help:"(required) Path to directory or file to upload"`
	BuildUuid       string            `help:"Module Build UUID"`
	Configuration   string            `help:"Build type, like 'debug' or 'release'"`
	Path            utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	VersionCode     string            `help:"Module version code"`
	VersionName     string            `help:"Module version name"`
}
