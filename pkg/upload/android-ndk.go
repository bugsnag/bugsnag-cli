package upload

import (
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type AndroidNdkMapping struct {
	AndroidNdkRoot  utils.UploadPaths `help:"Path to Android NDK installation ($ANDROID_NDK_ROOT)"`
	AppManifestPath utils.UploadPaths `help:"(required) Path to directory or file to upload"`
	BuildUuid       string            `help:"Module Build UUID"`
	Configuration   string            `help:"Build type, like 'debug' or 'release'"`
	Path            utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	ProjectRoot     utils.UploadPaths `help:"path to remove from the beginning of the filenames in the mapping file"`
	VersionCode     string            `help:"Module version code"`
	VersionName     string            `help:"Module version name"`
}
