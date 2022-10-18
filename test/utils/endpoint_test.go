package utils_testing

import (
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestEndpointBuilding(t *testing.T) {
	results := utils.BuildEndpointUrl("", 0)
	assert.Equal(t, results, "https://upload.bugsnag.com", "They should be the same")

	results = utils.BuildEndpointUrl("https://localhost", 0)
	assert.Equal(t, results, "https://localhost", "They should be the same")

	results = utils.BuildEndpointUrl("", 8443)
	assert.Equal(t, results, "https://upload.bugsnag.com:8443", "They should be the same")

	results = utils.BuildEndpointUrl("https://localhost", 8443)
	assert.Equal(t, results, "https://localhost:8443", "They should be the same")
}
