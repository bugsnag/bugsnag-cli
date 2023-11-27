package build

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"io"
	"net/http"
	"strings"
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

func ProcessBuildRequest(buildOptions CreateBuildInfo, endpoint string, dryRun bool) error {

	// Validate the required options for the API
	err := buildOptions.Validate()

	if err != nil {
		return err
	}

	buildPayload, err := json.Marshal(buildOptions)

	if err != nil {
		log.Error("Failed to create build information payload: "+err.Error(), 1)
	}

	prettyBuildPayload, _ := utils.PrettyPrintJson(string(buildPayload))
	log.Info("Build information: \n" + prettyBuildPayload)

	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(buildPayload))

	req.Header.Add("Content-Type", "application/json")

	if !dryRun {
		res, err := server.SendRequest(req, 300)

		if err != nil {
			return fmt.Errorf("error sending file request: %w", err)
		}

		b, err := io.ReadAll(res.Body)

		if strings.Contains(string(b), "Source control provider is missing") {
			log.Info("Source control provider is missing and could not be inferred. Please resend using one of: [github-enterprise, github, gitlab-onpremise, gitlab, bitbucket-server, bitbucket]. Request was still processed but source control information was ignored.")
		}

		if err != nil {
			return fmt.Errorf("error reading body from response: %w", err)
		}

		if res.StatusCode != 200 {
			return fmt.Errorf("%s : %s", res.Status, string(b))
		}
	} else {
		log.Info("(dryrun) Skipping sending build information to " + endpoint)
	}

	return nil
}
