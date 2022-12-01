package build

import (
	"github.com/alecthomas/kong"
)

type CreateBuild struct {
	AppVersion       string            `help:"The version number of the application."`
	AppVersionCode   string            `help:"The version code for the application (Android only)."`
	AppBundleVersion string            `help:"The bundle version for the application (iOS only)."`
	BuilderName      string            `help:"The name of the entity that triggered the build. Could be a user, system etc."`
	Metadata         map[string]string `help:"Additional build information"`
	ReleaseStage     string            `help:"The release stage (eg, production, staging) that is being released (if applicable)."`
	Provider         string            `help:"he name of the source control provider that contains the source code for the build."`
	Repository       string            `help:"The URL of the repository containing the source code being deployed."`
	Revision         string            `help:"The source control SHA-1 hash for the code that has been built (short or long hash)"`
}
