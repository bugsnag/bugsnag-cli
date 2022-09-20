package upload

type DartSymbol struct {
	Path []string `arg:"" name:"path" help:"Path to directory or file to upload" type:"path"`

	AppVersion string `help:"(optional) the version of the application."`
	AppVersionCode string `help:"(optional) the version code for the application (Android only)."`
	AppBundleVersion string `help:"(optional) the bundle version for the application (iOS only)."`
	IosAppPath string `help:"(optional) the path to the built IOS app."`
	Platform string `help:"(optional) the application platform, either android or ios"`
}