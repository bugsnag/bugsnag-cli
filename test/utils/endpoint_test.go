package utils_testing

import (
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestEndpointBuilding(t *testing.T) {
	t.Log("Testing the default endpoint")
	results := utils.BuildEndpointUrl("", 0)
	assert.Equal(t, results, "https://upload.bugsnag.com", "They should be the same")

	t.Log("Testing setting an endpoint")
	results = utils.BuildEndpointUrl("https://localhost", 0)
	assert.Equal(t, results, "https://localhost", "They should be the same")

	t.Log("Testing setting a port")
	results = utils.BuildEndpointUrl("", 8443)
	assert.Equal(t, results, "https://upload.bugsnag.com:8443", "They should be the same")

	t.Log("Testing setting an endpoint and port")
	results = utils.BuildEndpointUrl("https://localhost", 8443)
	assert.Equal(t, results, "https://localhost:8443", "They should be the same")
}
