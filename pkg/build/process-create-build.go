package build

import (
	"encoding/json"
	"fmt"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type Payload struct {
	ApiKey           string            `json:"apiKey,omitempty"`
	BuilderName      string            `json:"builderName,omitempty"`
	ReleaseStage     string            `json:"releaseStage,omitempty"`
	SourceControl    SourceControl     `json:"sourceControl,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
	AppVersion       string            `json:"appVersion,omitempty"`
	AppVersionCode   string            `json:"appVersionCode,omitempty"`
	AppBundleVersion string            `json:"appBundleVersion,omitempty"`
}

// ProcessCreateBuild processes a build request by creating a payload from the provided
// build options, logging the build information, and sending an HTTP request to the specified endpoint.
//
// Parameters:
//   - buildOptions: An instance of CreateBuildInfo containing information for the build.
//   - endpoint: The target URL for the HTTP request.
//   - options: CLI options used for this command.
//
// Returns:
//   - error: An error if any step of the build processing fails. Nil if the process is successful.
func ProcessCreateBuild(
	buildOptions CreateBuildInfo,
	endpoint string,
	options options.CLI,
	logger log.Logger,
) error {
	buildPayload, err := json.Marshal(options.CreateBuild)
	if err != nil {
		return fmt.Errorf("Failed to create build information payload: %s", err.Error())
	}

	prettyBuildPayload, _ := utils.PrettyPrintJson(string(buildPayload))
	logger.Debug(fmt.Sprintf("Build information:\n%s", prettyBuildPayload))

	err = server.ProcessBuildRequest(endpoint, buildPayload, options, logger)
	if err != nil {
		return err
	}

	return nil
}
