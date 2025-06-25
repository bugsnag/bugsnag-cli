package options

import "github.com/bugsnag/bugsnag-cli/pkg/utils"

// UnityLineMapping is used to specify options for uploading a Unity IL2CPP mapping file.
type UnityLineMapping struct {
	NoUploadIl2cppMapping bool       `name:"no-upload-il2cpp-mapping" help:"Do not upload the IL2CPP mapping file (LineNumberMappings.json)"`
	UploadIl2cppMapping   utils.Path `name:"upload-il2cpp-mapping" help:"The path to the IL2CPP mapping file (LineNumberMappings.json) to upload"`
}

// UnityAndroid is used to specify options for uploading Unity symbols and AAB files.
type UnityAndroid struct {
	Path          utils.Paths `arg:"" name:"path" help:"The path to the Unity symbols (.zip) file to upload (or directory containing it)" type:"path" default:"."`
	AabPath       utils.Path  `help:"The path to an AAB file to upload alongside the Unity symbols"`
	ApplicationId string      `help:"A unique application ID, usually the package name, of the application"`
	BuildUuid     string      `help:"A unique identifier for this build of the application" xor:"no-build-uuid,build-uuid"`
	NoBuildUuid   bool        `help:"Prevents the automatically generated build UUID being uploaded with the build" xor:"build-uuid,no-build-uuid"`
	ProjectRoot   string      `help:"The path to strip from the beginning of source file names referenced in stacktraces on the BugSnag dashboard" type:"path"`
	VersionCode   string      `help:"The version code of this build of the application"`
	VersionName   string      `help:"The version of the application"`
	Overwrite     bool        `help:"Whether to ignore and overwrite existing uploads with same identifier, rather than failing if a matching file exists"`

	UnityShared UnityLineMapping `embed:""`
}

type UnityIos struct {
	Path          utils.Paths `arg:"" name:"path" help:"The path to the Unity symbols (.zip) file to upload (or directory containing it)" type:"path" default:"."`
	ApplicationId string      `help:"A unique application ID, usually the package name, of the application"`
	BundleVersion string      `help:"The bundle version of this build of the application (Apple platforms only)"`
	VersionName   string      `help:"The version of the application"`
	Overwrite     bool        `help:"Whether to ignore and overwrite existing uploads with same identifier, rather than failing if a matching file exists"`

	UnityShared UnityLineMapping `embed:""`
	DsymShared  DsymShared       `embed:""`
}
