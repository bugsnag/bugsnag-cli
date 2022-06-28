package main

import "github.com/alecthomas/kong"

type CreateBuild struct {
	ApiKey           string `name:"apikey" help:"Project API key"`
	AppVersion       string `name:"app-version" help:"App version"`
	AppVersionCode   int    `name:"version-code" help:"App version code (Android)"`
	AppBundleVersion string `name:"bundle-version" help:"App bundle version (Apple platforms)"`
	BuilderName      string `name:"builder-name" help:"Name of the build creator"`
	Metadata         map[string]string `help:"Additional build information"`
	ReleaseStage     string `name:"release-stage" help:"Project API key"`
}

func SendBuildInfo(ctx *kong.Context) error {
	return nil
}
