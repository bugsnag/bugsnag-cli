package utils

import (
	"fmt"
	"os"
	"strings"
)

type Paths []string
type Path string
type Provider string
type LogLevels string
type Platform string

// Validate that the path(s) exist
func (p Paths) Validate() error {
	for _, path := range p {
		if _, err := os.Stat(path); err != nil {
			return err
		}
	}
	return nil
}

// Validate that the path exists
func (p Path) Validate() error {
	if _, err := os.Stat(string(p)); err != nil {
		return err
	}
	return nil
}

// ContainsString Check if a string is in a slice
func ContainsString(slice []string, target string) bool {
	for _, element := range slice {
		if strings.Contains(element, target) {
			return true
		}
	}
	return false
}

// Validate Check that the source control provider is valid
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

// Validate that the log level is valid
func (l LogLevels) Validate() error {
	switch strings.ToLower(string(l)) {
	case "debug", "info", "warn", "fatal":
		return nil
	default:
		return fmt.Errorf("invalid log level: %s", l)
	}
}

// Validate that the platform is valid
func (p Platform) Validate() error {
	switch strings.ToLower(string(p)) {
	case "android", "ios", "vega":
		return nil
	default:
		return fmt.Errorf("invalid platform: %s", p)
	}
}
