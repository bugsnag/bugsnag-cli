package options

import (
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// Global CLI options
type Globals struct {
	ApiKey            string            `help:"The BugSnag API key for the application"`
	DryRun            bool              `help:"Performs a dry-run of the command without sending any information to BugSnag"`
	FailOnUploadError bool              `help:"Stops the upload when a file fails to upload successfully" default:"false"`
	LogLevel          string            `help:"Sets the level of logging to debug, info, warn or fatal" default:"info"`
	Port              int               `help:"The port number for the BugSnag upload server" default:"443"`
	Verbose           bool              `name:"verbose" help:"Sets the level of the logging to its highest."`
	Version           utils.VersionFlag `name:"version" help:"Prints the version information for this CLI"`
}

// Unique CLI options
type CLI struct {
	Globals

	CreateAndroidBuildId CreateAndroidBuildId `cmd:"" help:"Generate a reproducible Build ID from .dex files"`
	CreateBuild          CreateBuild          `cmd:"" help:"Provide extra information whenever you build, release, or deploy your application"`
	Upload               struct {
		// shared options
		Overwrite        bool   `help:"Whether to ignore and overwrite existing uploads with same identifier, rather than failing if a matching file exists"`
		Retries          int    `help:"The number of retry attempts before failing an upload request" default:"0"`
		Timeout          int    `help:"The number of seconds to wait before failing an upload request" default:"300"`
		UploadAPIRootUrl string `help:"The upload server hostname, optionally containing port number" default:"https://upload.bugsnag.com"`

		// required options
		All                upload.DiscoverAndUploadAny   `cmd:"" help:"Upload any symbol/mapping files"`
		AndroidAab         upload.AndroidAabMapping      `cmd:"" help:"Process and upload application bundle files for Android"`
		AndroidNdk         upload.AndroidNdkMapping      `cmd:"" help:"Process and upload NDK symbol files for Android"`
		AndroidProguard    upload.AndroidProguardMapping `cmd:"" help:"Process and upload Proguard/R8 mapping files for Android"`
		DartSymbol         upload.DartSymbolOptions      `cmd:"" help:"Process and upload symbol files for Flutter" name:"dart"`
		Dsym               upload.Dsym                   `cmd:"" help:"Upload dSYMs for iOS"`
		Js                 upload.JsOptions              `cmd:"" help:"Upload source maps for JavaScript"`
		ReactNativeAndroid upload.ReactNativeAndroid     `cmd:"" help:"Upload source maps for React Native Android"`
		ReactNativeIos     upload.ReactNativeIos         `cmd:"" help:"Upload source maps for React Native iOS"`
		UnityAndroid       upload.UnityAndroid           `cmd:"" help:"Upload Android mappings and NDK symbol files from Unity projects"`
	} `cmd:"" help:"Upload symbol/mapping files"`
}

type CreateAndroidBuildId struct {
	Path utils.Paths `arg:"" name:"path" help:"Path to the project directory" type:"path"`
}

type AndroidBuildOptions struct {
	AndroidAab  utils.Path `help:"The path to an Android AAB file from which to obtain build information"`
	AppManifest utils.Path `help:"The path to an Android manifest file (AndroidManifest.xml) from which to obtain build information"`
	VersionCode string     `help:"The version code of this build of the application (Android only)." aliases:"app-version-code,version-code" xor:"version-code,bundle-version"`
}

type IosBuildOptions struct {
	BundleVersion string `help:"The bundle version of this build of the application (Apple platforms only)" aliases:"app-bundle-version,bundle-version"`
}

type CreateBuild struct {
	Path              utils.Paths       `arg:"" name:"path" help:"Path to the project directory" type:"path" default:"."`
	AutoAssignRelease bool              `help:"Whether to automatically associate this build with any new error events and sessions that are received for the release stage"`
	BuildApiRootUrl   string            `help:"The build server hostname, optionally containing port number" default:"https://build.bugsnag.com"`
	BuilderName       string            `help:"The name of the person or entity who built the app"`
	Metadata          map[string]string `help:"Custom build information to be associated with the release on the BugSnag dashboard"`
	Provider          utils.Provider    `help:"The name of the source control provider that contains the source code for the build"`
	ReleaseStage      string            `help:"The release stage (eg, production, staging) of the application build"`
	Repository        string            `help:"The URL of the repository containing the source code that has been built."`
	Retries           int               `help:"The number of retry attempts before failing the command" default:"0"`
	Revision          string            `help:"The source control SHA-1 hash for the code that has been built (short or long hash)"`
	Timeout           int               `help:"The number of seconds to wait before failing the command" default:"300"`
	VersionName       string            `help:"The version of the application" aliases:"app-version,version-name"`
	AndroidBuildOptions
	IosBuildOptions
}
