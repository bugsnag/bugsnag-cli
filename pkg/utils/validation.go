package utils

import (
	"fmt"
	"os"
	"strings"
)

type Paths []string
type Path string
type Provider string

// Validate that the path(s) exist
func (p Paths) Validate() error {
	for _, path := range p {
		if _, err := os.Stat(path); err != nil {
			return err
		}
	}
	return nil
}

// Validate that the path exist
func (p Path) Validate() error {
	if _, err := os.Stat(string(p)); err != nil {
		return err
	}
	return nil
}

func ContainsString(slice []string, target string) bool {
	for _, element := range slice {
		if strings.Contains(element, target) {
			return true
		}
	}
	return false
}

func (p *Provider) Validate() error {
	switch strings.ToLower(string(*p)) {
	case "github", "github-enterprise", "bitbucket", "bitbucket-server", "gitlab", "gitlab-onpremise":
		return nil
	case "":
		return fmt.Errorf("missing source control provider, please specify using `--provider`. Accepted values are: github, github-enterprise, bitbucket, bitbucket-server, gitlab, gitlab-onpremise")
	default:
		return fmt.Errorf("%s is not an accepted value for the source control provider. Accepted values are: github, github-enterprise, bitbucket, bitbucket-server, gitlab, gitlab-onpremise", *p)
	}
}
