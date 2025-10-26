package options

import "github.com/bugsnag/bugsnag-cli/pkg/utils"

type ReactNative struct {
	Path            utils.Paths                `arg:"" name:"path" help:"The path to the root of the React Native project to upload files from" type:"path" default:"."`
	ProjectRoot     string                     `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`
	Shared          ReactNativeShared          `embed:""`
	AndroidSpecific ReactNativeAndroidSpecific `embed:"" prefix:"android-"`
	IosSpecific     ReactNativeIosSpecific     `embed:"" prefix:"ios-"`
	Overwrite       bool                       `help:"Whether to ignore and overwrite existing uploads with same identifier, rather than failing if a matching file exists"`
}

type ReactNativeAndroidSpecific struct {
	AppManifest string `help:"The path to a manifest file (AndroidManifest.xml) from which to obtain build information" type:"path"`
	Variant     string `help:"The build type/flavor (e.g. debug, release) used to disambiguate the between built files when searching the project directory"`
	VersionCode string `help:"The version code of this build of the application"`
}

type ReactNativeAndroid struct {
	Path        utils.Paths `arg:"" name:"path" help:"The path to the root of the React Native project to upload files from" type:"path" default:"."`
	ProjectRoot string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`

	ReactNative ReactNativeShared          `embed:""`
	Android     ReactNativeAndroidSpecific `embed:""`
	Overwrite   bool                       `help:"Whether to ignore and overwrite existing uploads with same identifier, rather than failing if a matching file exists"`
}

type ReactNativeShared struct {
	Bundle       string `help:"The path to the bundled JavaScript file to upload" type:"path"`
	CodeBundleId string `help:"A unique identifier for the JavaScript bundle"`
	Dev          bool   `help:"Indicates whether this is a debug or release build"`
	SourceMap    string `help:"The path to the source map file to upload" type:"path"`
	VersionName  string `help:"The version of the application"`
}

type ReactNativeIosSpecific struct {
	BundleVersion string     `help:"The bundle version of this build of the application (Apple platforms only)"`
	Plist         string     `help:"The path to a .plist file from which to obtain build information" type:"path"`
	Scheme        string     `help:"The name of the Xcode options.Ios.Scheme used to build the application"`
	XcodeProject  string     `help:"The path to an Xcode project, workspace or containing directory from which to obtain build information" type:"path"`
	XcarchivePath utils.Path `help:"The path to the Xcode archive to process if it has been exported" type:"path"`
}

type ReactNativeIos struct {
	Path        utils.Paths `arg:"" name:"path" help:"The path to the root of the React Native project to upload files from" type:"path" default:"."`
	ProjectRoot string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`

	ReactNative ReactNativeShared      `embed:""`
	Ios         ReactNativeIosSpecific `embed:""`
	Overwrite   bool                   `help:"Whether to ignore and overwrite existing uploads with same identifier, rather than failing if a matching file exists"`
}

type ReactNativeSourcemaps struct {
	Path          utils.Path     `arg:"" name:"path" help:"The path to the root of the React Native project to upload files from" type:"path" default:"."`
	ProjectRoot   string         `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`
	VersionName   string         `help:"The version of the application"`
	VersionCode   string         `help:"The version code of this build of the application" xor:"version-code,bundle-version"`
	BundleVersion string         `help:"The bundle version of this build of the application (Apple platforms only)" xor:"version-code,bundle-version"`
	CodeBundleId  string         `help:"A unique identifier for the JavaScript bundle"`
	Dev           bool           `help:"Indicates whether this is a debug or release build"`
	Overwrite     bool           `help:"Whether to ignore and overwrite existing uploads with same identifier, rather than failing if a matching file exists"`
	SourceMap     string         `help:"The path to the source map file to upload" type:"path"`
	Bundle        string         `help:"The path to the bundled JavaScript file to upload" type:"path"`
	Platform      utils.Platform `help:"The platform of the React Native build"`
}
