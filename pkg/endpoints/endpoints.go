package endpoints

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/bugsnag/bugsnag-cli/pkg/options"
)

// Constants defining upload and build endpoint instances.
const (
	SECONDARY_API_PREFIX      = "00000" // API keys starting with this indicate usage of the secondary instance.
	SECONDARY_UPLOAD_ENDPOINT = "https://upload.bugsnag.smartbear.com"
	SECONDARY_BUILD_ENDPOINT  = "https://build.bugsnag.smartbear.com"
	PRIMARY_UPLOAD_ENDPOINT   = "https://upload.bugsnag.com"
	PRIMARY_BUILD_ENDPOINT    = "https://build.bugsnag.com"
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
// The server passed in as the endpoint option is used, if provided. Otherwise the API key is used to determine the appropriate instance.
// If the endpoint URL cannot be built, it returns an error.
//
// Parameters:
//   - apiKey: the project API key.
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
	} else if strings.HasPrefix(apiKey, SECONDARY_API_PREFIX) {
		endpoint = SECONDARY_UPLOAD_ENDPOINT
	} else {
		endpoint = PRIMARY_UPLOAD_ENDPOINT
	}

	endpoint, err := BuildEndpointURL(endpoint+endpointPath, options.Port)

	if err != nil {
		return endpoint, fmt.Errorf("error building upload endpoint URL: %w", err)
	}

	return endpoint, nil
}

// GetDefaultBuildEndpoint selects the appropriate build endpoint based on the API key.
//
// The server passed in as the endpoint option is used, if provided. Otherwise the API key is used to determine the appropriate instance.
// If the endpoint URL cannot be built, it returns an error.
//
// Parameters:
//   - apiKey: the project API key.
//   - options: CLI options that may contain a custom upload API root URL and port.
//
// Returns:
//   - A string containing the resolved build endpoint.
//   - An error if the endpoint URL cannot be built.
func GetDefaultBuildEndpoint(apiKey string, options options.CLI) (string, error) {
	var endpoint string

	if options.CreateBuild.BuildApiRootUrl != "" {
		endpoint = options.CreateBuild.BuildApiRootUrl
	} else if strings.HasPrefix(apiKey, SECONDARY_API_PREFIX) {
		endpoint = SECONDARY_BUILD_ENDPOINT
	} else {
		endpoint = PRIMARY_BUILD_ENDPOINT
	}

	endpoint, err := BuildEndpointURL(endpoint, options.Port)

	if err != nil {
		return endpoint, fmt.Errorf("error building upload endpoint URL: %w", err)
	}

	return endpoint, nil
}
