package options

import (
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// Global CLI options
type Globals struct {
	UploadAPIRootUrl  string            `help:"Bugsnag On-Premise upload server URL. Can contain port number" default:"https://upload.bugsnag.com"`
	BuildApiRootUrl   string            `help:"Bugsnag On-Premise build server URL. Can contain port number" default:"https://build.bugsnag.com"`
	Port              int               `help:"Port number for the upload server" default:"443"`
	ApiKey            string            `help:"(required) Bugsnag integration API key for this application"`
	FailOnUploadError bool              `help:"Stops the upload when a mapping file fails to upload to Bugsnag successfully" default:"false"`
	Version           utils.VersionFlag `name:"version" help:"Print version information and quit"`
	DryRun            bool              `help:"Validate but do not process"`
}

// Unique CLI options
type CLI struct {
	Globals

	Upload struct {
		// shared options
		Overwrite bool `help:"Whether to overwrite any existing symbol file with a matching ID"`
		Timeout   int  `help:"Number of seconds to wait before failing an upload request" default:"300"`
		Retries   int  `help:"Number of retry attempts before failing an upload request" default:"0"`

		// required options
		AndroidAab         upload.AndroidAabMapping      `cmd:"" help:"Process and upload application bundle files for Android"`
		All                upload.DiscoverAndUploadAny   `cmd:"" help:"Upload any symbol/mapping files"`
		AndroidNdk         upload.AndroidNdkMapping      `cmd:"" help:"Process and upload Proguard mapping files for Android"`
		AndroidProguard    upload.AndroidProguardMapping `cmd:"" help:"Process and upload NDK symbol files for Android"`
		DartSymbol         upload.DartSymbolOptions      `cmd:"" help:"Process and upload symbol files for Flutter" name:"dart"`
		ReactNativeAndroid upload.ReactNativeAndroid     `cmd:"" help:"Upload source maps for React Native Android"`
		ReactNativeIos     upload.ReactNativeIos         `cmd:"" help:"Upload source maps for React Native iOS"`
		UnityAndroid       upload.UnityAndroid           `cmd:"" help:"Upload Android mappings and NDK symbol files from Unity projects"`
	} `cmd:"" help:"Upload symbol/mapping files"`
	CreateBuild          CreateBuild          `cmd:"" help:"Provide extra information whenever you build, release, or deploy your application"`
	CreateAndroidBuildId CreateAndroidBuildId `cmd:"" help:"Generate a reproducible Build ID from .dex files"`
}

type CreateAndroidBuildId struct {
	Path utils.Paths `arg:"" name:"path" help:"Path to the project directory" type:"path"`
}

type AndroidBuildOptions struct {
	VersionCode string     `help:"The version code for the application (Android only)." xor:"app-version-code,version-code"`
	AppManifest utils.Path `help:"TBC"`
	AndroidAab  utils.Path `help:"TBC"`
}

type IosBuildOptions struct {
	BundleVersion string `help:"The bundle version for the application (iOS only)." xor:"app-bundle-version,bundle-version"`
}

type CreateBuild struct {
	BuilderName       string            `help:"The name of the entity that triggered the build. Could be a user, system etc."`
	Metadata          map[string]string `help:"Additional build information"`
	ReleaseStage      string            `help:"The release stage (eg, production, staging) that is being released (if applicable)."`
	Provider          string            `help:"The name of the source control provider that contains the source code for the build."`
	Repository        string            `help:"The URL of the repository containing the source code being deployed."`
	Revision          string            `help:"The source control SHA-1 hash for the code that has been built (short or long hash)"`
	Path              utils.Paths       `arg:"" name:"path" help:"Path to the project directory" type:"path" default:"."`
	VersionName       string            `help:"The version of the application." xor:"app-version,version-name"`
	AutoAssignRelease bool              `help:"TBC"`
	AndroidBuildOptions
	IosBuildOptions
}
