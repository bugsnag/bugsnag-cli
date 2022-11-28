package utils

import (
	"fmt"
	"net/url"
)

func BuildEndpointUrl(uri string, port int) (string, error) {
	baseUrl, err := url.Parse(uri)

	if err != nil {
		return "", err
	}

	if baseUrl.Port() != "" {
		return baseUrl.String(), nil
	}

	if port != 0 {
		return fmt.Sprintf("%s:%d", baseUrl, port), nil
	}

	return baseUrl.String(), nil
}
