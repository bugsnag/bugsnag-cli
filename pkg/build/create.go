package build

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/log"
	"github.com/bugsnag/bugsnag-cli/pkg/server"
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
)

type CreateBuild struct {
	BuilderName  string            `help:"The name of the entity that triggered the build. Could be a user, system etc."`
	Metadata     map[string]string `help:"Additional build information"`
	ReleaseStage string            `help:"The release stage (eg, production, staging) that is being released (if applicable)."`
	Provider     string            `help:"The name of the source control provider that contains the source code for the build."`
	Repository   string            `help:"The URL of the repository containing the source code being deployed."`
	Revision     string            `help:"The source control SHA-1 hash for the code that has been built (short or long hash)"`
}

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

type SourceControl struct {
	Provider   string `json:"provider,omitempty"`
	Repository string `json:"repository,omitempty"`
	Revision   string `json:"revision,omitempty"`
}

func ProcessBuildRequest(apiKey string, builderName string, releaseStage string, provider string, repository string, revision string, appVersion string, appVersionCode string, appBundleVersion string, metadata map[string]string, endpoint string) error {
	if appVersion == "" {
		log.Error("Missing app version, please provide this via the command line options", 1)
	}

	builderName, err := SetBuilderName(builderName)

	if err != nil {
		log.Error("Failed to set builder name from system. Please provide this via the command line options. "+err.Error(), 1)
	}

	repoInfo := GetRepoInfo(provider, repository, revision)

	payload := Payload{
		ApiKey:       apiKey,
		BuilderName:  builderName,
		ReleaseStage: releaseStage,
		SourceControl: SourceControl{
			Provider:   repoInfo["provider"],
			Repository: repoInfo["repository"],
			Revision:   repoInfo["revision"],
		},
		Metadata:         metadata,
		AppVersion:       appVersion,
		AppVersionCode:   appVersionCode,
		AppBundleVersion: appBundleVersion,
	}

	buildPayload, err := json.Marshal(payload)

	if err != nil {
		log.Error("Failed to create build information payload: "+err.Error(), 1)
	}

	prettyBuildPayload, _ := utils.PrettyPrintJson(string(buildPayload))
	log.Info("Build information: \n" + prettyBuildPayload)

	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(buildPayload))

	req.Header.Add("Content-Type", "application/json")

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
	return nil
}

func GetRepoInfo(repoProvider string, repoUrl string, repoHash string) map[string]string {
	repoInfo := make(map[string]string)

	if repoUrl == "" {
		repoUrl, _ = utils.GetRepoUrl()
	}

	repoInfo["repository"] = repoUrl

	if repoProvider != "" {
		repoInfo["provider"] = repoProvider
	}

	if repoHash == "" {
		repoHash, _ = utils.GetCommitHash()
	}

	repoInfo["revision"] = repoHash

	return repoInfo
}

func SetBuilderName(name string) (string, error) {
	if name == "" {
		builder, err := utils.GetSystemUser()
		if err != nil {
			return name, err
		}
		return builder, nil
	}
	return name, nil
}
