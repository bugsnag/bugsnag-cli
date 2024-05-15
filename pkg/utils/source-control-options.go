package utils

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/log"
)

func SourceControl(provider string) string {
	if provider == "" {
		return provider
	}
	switch provider {
	case "github", "github-enterprise", "bitbucket", "bitbucket-server", "gitlab", "gitlab-onpremise":
		return provider
	default:
		log.Warn(fmt.Sprintf("%s is not an accepted value for the source control provider, proceeding without a provider", provider))
		return ""
	}
}
