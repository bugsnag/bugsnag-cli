package upload

import (
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type ReactNativeCocoa struct {
	AppVersion       string      `help:"The version of the application."`
	AppBundleVersion string      `help:"Bundle version for the application. (iOS only)"`
	Scheme           string      `help:"The name of the scheme to use when building the application."`
	SourceMap        string      `help:"Path to the source map file" type:"path"`
	Bundle           string      `help:"Path to the bundle file" type:"path"`
	Plist            string      `help:"Path to the Info.plist file" type:"path"`
	Xcworkspace      string      `help:"Path to the .xcworkspace file" type:"path"`
	CodeBundleID     string      `help:"A unique identifier to identify a code bundle release when using tools like CodePush"`
	Dev              bool        `help:"Indicates whether the application is a debug or release build"`
	ProjectRoot      string      `help:"path to remove from the beginning of the filenames in the mapping file" type:"path"`
	Path             utils.Paths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
}

func ProcessReactNativeCocoa(
	apiKey string,
	appVersion string,
	appBundleVersion string,
	scheme string,
	sourceMapPath string,
	bundlePath string,
	plistPath string,
	xcworkspacePath string,
	codeBundleId string,
	dev bool,
	projectRoot string,
	paths []string,
	endpoint string,
	timeout int,
	retries int,
	overwrite bool,
	dryRun bool,
) error {
	return nil
}
