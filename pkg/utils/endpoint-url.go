package utils

import (
	"fmt"
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
	if uri == "" {
		uri = BUGSNAG_UPLOAD
	}

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
// If the endpoint matches the default Bugsnag upload URL and the API key starts with
// HUB_PREFIX, it switches to the InsightHub upload URL.
//
// Parameters:
//   - endpoint: the current upload endpoint (may be Bugsnag or Hub).
//   - apiKey: the API key used to determine which backend to target.
//
// Returns:
//   - A string containing the resolved upload endpoint.
func GetDefaultUploadEndpoint(endpoint string, apiKey string) string {
	if strings.Contains(endpoint, BUGSNAG_UPLOAD) {
		if strings.HasPrefix(apiKey, HUB_PREFIX) {
			endpoint = HUB_UPLOAD
		} else {
			endpoint = BUGSNAG_UPLOAD
		}
	}

	return endpoint
}

// GetDefaultBuildEndpoint selects the appropriate build endpoint based on the API key.
//
// If the endpoint matches the default Bugsnag build URL and the API key starts with
// HUB_PREFIX, it switches to the InsightHub build URL.
//
// Parameters:
//   - endpoint: the current build endpoint (may be Bugsnag or Hub).
//   - apiKey: the API key used to determine which backend to target.
//
// Returns:
//   - A string containing the resolved build endpoint.
func GetDefaultBuildEndpoint(endpoint string, apiKey string) string {
	if strings.Contains(endpoint, BUGSNAG_BUILD) {
		if strings.HasPrefix(apiKey, HUB_PREFIX) {
			endpoint = HUB_BUILD
		} else {
			endpoint = BUGSNAG_BUILD
		}
	}

	return endpoint
}
