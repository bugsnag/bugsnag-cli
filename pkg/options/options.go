package options

import (
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// Global CLI options
type Globals struct {
	ApiKey   string            `help:"The BugSnag API key for the application"`
	DryRun   bool              `help:"Performs a dry-run of the command without sending any information to BugSnag"`
	LogLevel string            `help:"Sets the level of logging to debug, info, warn or fatal" default:"info"`
	Port     int               `help:"The port number for the BugSnag upload server" default:"443"`
	Verbose  bool              `name:"verbose" help:"Sets the level of the logging to its highest."`
	Version  utils.VersionFlag `name:"version" help:"Prints the version information for this CLI"`
}

type DiscoverAndUploadAny struct {
	Path          utils.Paths       `arg:"" name:"path" help:"(required) Path to directory or file to upload" type:"path"`
	UploadOptions map[string]string `help:"Additional arguments to pass to the upload request" mapsep:","`
}

type AndroidAabMapping struct {
	Path          utils.Paths `arg:"" name:"path" help:"The path to the AAB file to upload (or directory containing it)" type:"path" default:"."`
	ApplicationId string      `help:"A unique application ID, usually the package name, of the application"`
	BuildUuid     string      `help:"A unique identifier for this build of the application" xor:"no-build-uuid,build-uuid"`
	NoBuildUuid   bool        `help:"Prevents the automatically generated build UUID being uploaded with the build" xor:"build-uuid,no-build-uuid"`
	ProjectRoot   string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`
	VersionCode   string      `help:"The version code of this build of the application"`
	VersionName   string      `help:"The version of the application"`
}

type AndroidNdkMapping struct {
	Path           utils.Paths `arg:"" name:"path" help:"The path to the directory or file to upload" type:"path" default:"."`
	ApplicationId  string      `help:"A unique application ID, usually the package name, of the application"`
	AndroidNdkRoot string      `help:"The path to your NDK installation, used to access the objcopy tool for extracting symbol information"`
	AppManifest    string      `help:"The path to a manifest file (AndroidManifest.xml) from which to obtain build information" type:"path"`
	ProjectRoot    string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`
	Variant        string      `help:"The build type/flavor (e.g. debug, release) used to disambiguate the between built files when searching the project directory"`
	VersionCode    string      `help:"The version code of this build of the application"`
	VersionName    string      `help:"The version of the application"`
}

type AndroidProguardMapping struct {
	Path          utils.Paths `arg:"" name:"path" help:"The path to the directory or file to upload" type:"path" default:"."`
	ApplicationId string      `help:"A unique application ID, usually the package name, of the application"`
	AppManifest   string      `help:"The path to a manifest file (AndroidManifest.xml) from which to obtain build information" type:"path"`
	BuildUuid     string      `help:"A unique identifier for this build of the application" xor:"no-build-uuid,build-uuid"`
	NoBuildUuid   bool        `help:"Prevents the automatically generated build UUID being uploaded with the build" xor:"build-uuid,no-build-uuid"`
	DexFiles      []string    `help:"The path to classes.dex files or directory used to calculate a build UUID" type:"path" default:""`
	Variant       string      `help:"The build type/flavor (e.g. debug, release) used to disambiguate the between built files when searching the project directory"`
	VersionCode   string      `help:"The version code of this build of the application"`
	VersionName   string      `help:"The version of the application"`
}

type DartSymbol struct {
	Path          utils.Paths `arg:"" name:"path" help:"The path to the directory or file to upload" type:"path"`
	BundleVersion string      `help:"The bundle version of this build of the application (Apple platforms only)" xor:"app-bundle-version,bundle-version"`
	IosAppPath    utils.Path  `help:"The path to the iOS application binary, used to determine a unique build ID." type:"path"`
	VersionName   string      `help:"The version of the application." xor:"app-version,version-name"`
	VersionCode   string      `help:"The version code of this build of the application (Android only)" xor:"app-version-code,version-code"`
}

type Dsym struct {
	Path   utils.Paths `arg:"" name:"path" help:"The path to the directory or file to upload" type:"path" default:"."`
	Shared DsymShared  `embed:""`
}

type XcodeBuild struct {
	Path   utils.Paths `arg:"" name:"path" help:"The path to the directory or file to upload" type:"path" default:"."`
	Shared DsymShared  `embed:""`
}

type XcodeArchive struct {
	Path   utils.Paths `arg:"" name:"path" help:"The path to the directory or file to upload" type:"path" default:"."`
	Shared DsymShared  `embed:""`
}

type DsymShared struct {
	IgnoreEmptyDsym    bool       `help:"Throw warnings instead of errors when a dSYM file is found, rather than the expected dSYM directory"`
	IgnoreMissingDwarf bool       `help:"Throw warnings instead of errors when a dSYM with missing DWARF data is found"`
	Configuration      string     `help:"The configuration used to build the application"`
	Scheme             string     `help:"The name of the Xcode options.Scheme used to build the application"`
	ProjectRoot        string     `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`
	Plist              utils.Path `help:"The path to a .plist file from which to obtain build information" type:"path"`
	XcodeProject       utils.Path `help:"The path to an Xcode project, workspace or containing directory from which to obtain build information" type:"path"`
}

type Js struct {
	Path         utils.Paths `arg:"" name:"path" help:"The path to the directory or file to upload" type:"path" default:"."`
	BaseUrl      string      `help:"For directory-based uploads, the URL of the base directory for the minified JavaScript files that the source maps relate to. The relative path is appended onto this for each file. Asterisks can be used as a wildcard."`
	Bundle       string      `help:"Path to the minified JavaScript file that the source map relates to. If this is not provided then the file will be obtained when an error event is received." type:"path"`
	BundleUrl    string      `help:"For single file uploads, the URL of the minified JavaScript file that the source map relates to. Asterisks can be used as a wildcard."`
	ProjectRoot  string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`
	SourceMap    string      `help:"Path to the source map file. This usually has the .min.js extension." type:"path"`
	VersionName  string      `help:"The version of the app that the source map applies to. Defaults to the version in the package.json file (if found)."`
	CodeBundleId string      `help:"A unique identifier for the JavaScript bundle"`
}

type ReactNative struct {
	Path            utils.Paths                `arg:"" name:"path" help:"The path to the root of the React Native project to upload files from" type:"path" default:"."`
	ProjectRoot     string                     `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`
	Shared          ReactNativeShared          `embed:""`
	AndroidSpecific ReactNativeAndroidSpecific `embed:"" prefix:"android-"`
	IosSpecific     ReactNativeIosSpecific     `embed:"" prefix:"ios-"`
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
}

type UnityAndroid struct {
	Path          utils.Paths `arg:"" name:"path" help:"The path to the Unity symbols (.zip) file to upload (or directory containing it)" type:"path"`
	AabPath       utils.Path  `help:"The path to an AAB file to upload alongside the Unity symbols"`
	ApplicationId string      `help:"A unique application ID, usually the package name, of the application"`
	BuildUuid     string      `help:"A unique identifier for this build of the application" xor:"no-build-uuid,build-uuid"`
	NoBuildUuid   bool        `help:"Prevents the automatically generated build UUID being uploaded with the build" xor:"build-uuid,no-build-uuid"`
	ProjectRoot   string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`
	VersionCode   string      `help:"The version code of this build of the application"`
	VersionName   string      `help:"The version of the application"`
}
type Breakpad struct {
	Path            utils.Paths `arg:"" name:"path" help:"The path to the symbol files (.sym) to upload (or directory containing them)" type:"path"`
	CpuArch         string      `help:"The CPU architecture that the module was built for"`
	CodeFile        string      `help:"The basename of the module"`
	DebugFile       string      `help:"The basename of the debug file"`
	DebugIdentifier string      `help:"The debug file's identifier"`
	ProductName     string      `help:"The product name"`
	ProjectRoot     string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`
	OsName          string      `help:"The name of the operating system that the module was built for"`
	VersionName     string      `help:"The version of the application"`
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
		All                DiscoverAndUploadAny   `cmd:"" help:"Upload any symbol/mapping files"`
		AndroidAab         AndroidAabMapping      `cmd:"" help:"Process and upload application bundle files for Android"`
		AndroidNdk         AndroidNdkMapping      `cmd:"" help:"Process and upload NDK symbol files for Android"`
		AndroidProguard    AndroidProguardMapping `cmd:"" help:"Process and upload Proguard/R8 mapping files for Android"`
		DartSymbol         DartSymbol             `cmd:"" help:"Process and upload symbol files for Flutter" name:"dart"`
		XcodeBuild         XcodeBuild             `cmd:"" help:"Upload dSYMs for iOS from a build"`
		Dsym               Dsym                   `cmd:"" help:"(deprecated) Upload dSYMs for iOS"`
		XcodeArchive       XcodeArchive           `cmd:"" help:"Upload dSYMs for iOS from a Xcode archive"`
		Js                 Js                     `cmd:"" help:"Upload source maps for JavaScript"`
		ReactNative        ReactNative            `cmd:"" help:"Upload source maps for React Native"`
		ReactNativeAndroid ReactNativeAndroid     `cmd:"" help:"Upload source maps for React Native Android"`
		ReactNativeIos     ReactNativeIos         `cmd:"" help:"Upload source maps for React Native iOS"`
		UnityAndroid       UnityAndroid           `cmd:"" help:"Upload Android mappings and NDK symbol files from Unity projects"`
		Breakpad           Breakpad               `cmd:"" help:"Upload breakpad .sym files"`
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
