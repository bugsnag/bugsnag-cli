package utils_testing

import (
	"fmt"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestEndpointBuilding(t *testing.T) {

	t.Log("Testing setting an endpoint with CLI defaults")
	results, err := utils.BuildEndpointUrl("https://upload.bugsnag.com", 443)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, fmt.Sprintf("%s:%d", "https://upload.bugsnag.com", 443), results, "They should be the same")

	t.Log("Testing setting a port with the CLI default URL")
	results, err = utils.BuildEndpointUrl("https://upload.bugsnag.com", 8443)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, results, "https://upload.bugsnag.com:8443", "They should be the same")

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

	t.Log("Testing setting an endpoint with no port")
	results, err = utils.BuildEndpointUrl("https://localhost", 0)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, results, "https://localhost", "They should be the same")
}
