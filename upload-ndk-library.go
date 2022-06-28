package main

type NdkLibrary struct {
	Path []string `arg:"" name:"path" help:"Path to upload" type:"path"`

	AndroidNdkRoot string `help:"Path to Android NDK installation" env:"ANDROID_NDK_ROOT"`
	AppManifestPath string `help:"Path to the merged app manifest"`
	Configuration string `help:"Build type, like 'debug' or 'release'"`
	ProjectRoot string `help:"path to remove from the beginning of the filenames in the mapping file"`
	VersionCode string `help:"Module version code"`
	VersionName string `help:"Module version name"`
}
