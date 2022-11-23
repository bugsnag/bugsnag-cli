package utils_testing

import (
	"fmt"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestEndpointBuilding(t *testing.T) {
	defaultUrl := "https://upload.bugsnag.com"
	defaultPort := 443

	t.Log("Testing setting an endpoint with CLI defaults")
	results, err := utils.BuildEndpointUrl(defaultUrl, defaultPort)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, fmt.Sprintf("%s:%d", defaultUrl, defaultPort), results, "They should be the same")

	t.Log("Testing setting a port with the CLI default URL")
	results, err = utils.BuildEndpointUrl(defaultUrl, 8443)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, results, defaultUrl + ":8443", "They should be the same")

	t.Log("Testing setting an endpoint and port")
	results, err = utils.BuildEndpointUrl("https://localhost", 8443)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, results, "https://localhost:8443", "They should be the same")

	t.Log("Testing setting an endpoint with a port and port")
	results, err = utils.BuildEndpointUrl("https://localhost:1234", 8443)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, results, "https://localhost:1234", "They should be the same")
}
