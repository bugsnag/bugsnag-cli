package upload

import (
	"fmt"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type JsOptions struct {
	VersionName string      `help:"The version of the app that the source map applies to (as set in the JavaScript notifier). Defaults to the version in the package.json file."`
	BundleUrl   string      `help:"The url of the minified JavaScript file that the source map relates to. Asterisks can be used as a wildcard."`
	SourceMap   string      `help:"Path to the source map file. This usually has the .min.js extension." type:"path"`
	Bundle      string      `help:"Path to the minified JavaScript file that the source map relates to. If this is not provided then the file will be obtained when an error event is received." type:"path"`
	ProjectRoot string      `help:"The root path of the project on the filesystem where the source map was built. This will be stripped off from the file name in the displayed stacktraces." type:"path"`
	Path        utils.Paths `arg:"" name:"path" help:"Path to directory or file to upload" type:"path" default:"."`
}

func ProcessJs(
	apiKey string,
	jsOptions JsOptions,
	endpoint string,
	timeout int,
	retries int,
	overwrite bool,
	dryRun bool,
	logger log.Logger,
) error {

	var rootDirPath string

	for _, path := range jsOptions.Path {
		rootDirPath = path

		logger.Info(fmt.Sprintf("Uploading js %v %v %v", apiKey, jsOptions, rootDirPath))

	}

	return nil
}
