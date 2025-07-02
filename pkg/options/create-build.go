package options

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"strings"
)

// UnmarshalText parses a comma-separated list of key=value pairs into Metadata.
//
// Parameters:
//   - text: A byte slice representing metadata in the format "key1=value1,key2=value2".
//
// Returns:
//   - error: An error if any entry is not formatted correctly.
func (t *Metadata) UnmarshalText(text []byte) error {
	result := make(map[string]string)
	parts := strings.Split(string(text), ",")
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("invalid format: %s", part)
		}
		result[kv[0]] = kv[1]
	}
	*t = result
	return nil
}

// AndroidBuildOptions holds build-specific options for Android builds.
type AndroidBuildOptions struct {
	AndroidAab  utils.Path `help:"The path to an Android AAB file from which to obtain build information"`
	AppManifest utils.Path `help:"The path to an Android manifest file (AndroidManifest.xml) from which to obtain build information"`
	VersionCode string     `help:"The version code of this build of the application (Android only)." aliases:"app-version-code,version-code" xor:"version-code,bundle-version"`
}

// IosBuildOptions holds build-specific options for iOS builds.
type IosBuildOptions struct {
	BundleVersion string `help:"The bundle version of this build of the application (Apple platforms only)" aliases:"app-bundle-version,bundle-version"`
}

// Metadata represents custom key-value data to be associated with a build.
type Metadata map[string]string

// CreateBuild defines the configuration for creating a new build through the CLI.
type CreateBuild struct {
	Path              utils.Paths    `arg:"" name:"path" help:"Path to the project directory" type:"path" default:"."`
	AutoAssignRelease bool           `help:"Whether to automatically associate this build with any new error events and sessions that are received for the release stage"`
	BuildApiRootUrl   string         `help:"The build server hostname, optionally containing port number"`
	BuilderName       string         `help:"The name of the person or entity who built the app"`
	Metadata          Metadata       `help:"Custom build information to be associated with the release on the BugSnag dashboard"`
	Provider          utils.Provider `help:"The name of the source control provider that contains the source code for the build"`
	ReleaseStage      string         `help:"The release stage (eg, production, staging) of the application build"`
	Repository        string         `help:"The URL of the repository containing the source code that has been built."`
	Retries           int            `help:"The number of retry attempts before failing the command" default:"0"`
	Revision          string         `help:"The source control SHA-1 hash for the code that has been built (short or long hash)"`
	Timeout           int            `help:"The number of seconds to wait before failing the command" default:"300"`
	VersionName       string         `help:"The version of the application" aliases:"app-version,version-name"`

	AndroidBuildOptions
	IosBuildOptions
}
