package endpoints

import (
	"fmt"
	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"net/url"
	"strings"
)

// Constants defining upload and build endpoints for Bugsnag and InsightHub.
const (
	HUB_PREFIX     = "00000" // API keys starting with this indicate usage of InsightHub instead of Bugsnag.
	HUB_UPLOAD     = "https://upload.insighthub.smartbear.com"
	HUB_BUILD      = "https://build.insighthub.smartbear.com"
	BUGSNAG_UPLOAD = "https://upload.bugsnag.com"
	BUGSNAG_BUILD  = "https://build.bugsnag.com"
)

// BuildEndpointURL constructs a complete URL from a base URI and optional port.
//
// If the URI already includes a port, it returns the URI as-is.
// Otherwise, if a non-zero port is provided, it appends the port to the base URI.
//
// Parameters:
//   - uri: a string representing the base URI (e.g., "http://localhost").
//   - port: the port number to append if not already specified.
//
// Returns:
//   - A string containing the full URI with port, if applicable.
//   - An error if the URI cannot be parsed.
func BuildEndpointURL(uri string, port int) (string, error) {
	baseURL, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	if baseURL.Port() != "" {
		return baseURL.String(), nil
	}

	if port != 0 {
		baseURL.Host = fmt.Sprintf("%s:%d", baseURL.Hostname(), port)
	}

	return baseURL.String(), nil
}

// GetDefaultUploadEndpoint selects the appropriate upload endpoint based on the API key.
//
// If the API key starts with HUB_PREFIX, it uses the InsightHub upload URL.
// Otherwise, it defaults to the Bugsnag upload URL. The endpointPath is appended to the base URL.
// If a custom upload API root URL is provided in the CLI options, it uses that instead.
// If the endpoint URL cannot be built, it returns an error.
//
// Parameters:
//   - apiKey: the API key used to determine which backend to target.
//   - endpointPath: the specific path to append to the base upload endpoint.
//   - options: CLI options that may contain a custom upload API root URL and port.
//
// Returns:
//   - A string containing the resolved upload endpoint.
//   - An error if the endpoint URL cannot be built.
func GetDefaultUploadEndpoint(apiKey string, endpointPath string, options options.CLI) (string, error) {
	var endpoint string

	if options.Upload.UploadAPIRootUrl != "" {
		endpoint = options.Upload.UploadAPIRootUrl
	} else if strings.HasPrefix(apiKey, HUB_PREFIX) {
		endpoint = HUB_UPLOAD
	} else {
		endpoint = BUGSNAG_UPLOAD
	}

	endpoint, err := BuildEndpointURL(endpoint+endpointPath, options.Port)

	if err != nil {
		return endpoint, fmt.Errorf("error building upload endpoint URL: %w", err)
	}

	return endpoint, nil
}

// GetDefaultBuildEndpoint selects the appropriate build endpoint based on the API key.
//
// If the API key starts with HUB_PREFIX, it uses the InsightHub build URL.
// Otherwise, it defaults to the Bugsnag build URL. If a custom build API root URL is provided
// in the CLI options, it uses that instead. The endpoint URL is built with the specified port.
// If the endpoint URL cannot be built, it returns an error.
//
// Parameters:
//   - apiKey: the API key used to determine which backend to target.
//   - options: CLI options that may contain a custom upload API root URL and port.
//
// Returns:
//   - A string containing the resolved build endpoint.
//   - An error if the endpoint URL cannot be built.
func GetDefaultBuildEndpoint(apiKey string, options options.CLI) (string, error) {
	var endpoint string

	if options.CreateBuild.BuildApiRootUrl != "" {
		endpoint = options.CreateBuild.BuildApiRootUrl
	} else if strings.HasPrefix(apiKey, HUB_PREFIX) {
		endpoint = HUB_BUILD
	} else {
		endpoint = BUGSNAG_BUILD
	}

	endpoint, err := BuildEndpointURL(endpoint, options.Port)

	if err != nil {
		return endpoint, fmt.Errorf("error building upload endpoint URL: %w", err)
	}

	return endpoint, nil
}
