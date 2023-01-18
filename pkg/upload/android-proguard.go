package upload

import (
	"path/filepath"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type AndroidProguardMapping struct {
	ApplicationId   string            `help:"Module application identifier"`
	AppManifestPath string            `help:"(required) Path to directory or file to upload" type:"path"`
	BuildUuid       string            `help:"Module Build UUID"`
	Configuration   string            `help:"Build type, like 'debug' or 'release'"`
	Path            utils.UploadPaths `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	VersionCode     string            `help:"Module version code"`
	VersionName     string            `help:"Module version name"`
	DryRun          bool              `help:"Validate but do not upload"`
}

func ProcessAndroidProguard(paths []string, applicationId string, appManifestPath string, buildUuid string, configuration string, versionCode string, versionName string, endpoint string, timeout int, retries int, overwrite bool, apiKey string, failOnUploadError bool, dryRun bool) error {

	return nil
}
