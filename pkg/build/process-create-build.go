package build

import (
	"encoding/json"
	"fmt"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

// ProcessCreateBuild marshals build metadata into JSON and sends it to the build endpoint.
//
// Parameters:
//   - buildOptions: A structure containing all metadata for the build (implements CreateBuildInfo).
//   - options: CLI options including endpoint and retry configuration.
//   - logger: Logger used for debug and error output.
//
// Returns:
//   - error: Non-nil if JSON marshalling fails or the request to the server fails.
func ProcessCreateBuild(
	buildOptions CreateBuildInfo,
	options options.CLI,
	logger log.Logger,
) error {
	// Marshal the build options into a JSON payload
	buildPayload, err := json.Marshal(buildOptions)
	if err != nil {
		return fmt.Errorf("Failed to create build information payload: %s", err.Error())
	}

	// Output the build payload in a human-readable format for debugging
	prettyBuildPayload, _ := utils.PrettyPrintJson(string(buildPayload))
	logger.Debug(fmt.Sprintf("Build information:\n%s", prettyBuildPayload))

	// Send the build payload to the configured Bugsnag build endpoint
	err = server.ProcessBuildRequest(buildOptions.ApiKey, buildPayload, options, logger)
	if err != nil {
		return err
	}

	return nil
}
