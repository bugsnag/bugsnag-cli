package utils_testing_test

import (
	"os"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestParseGitUrl(t *testing.T) {
	t.Log("Testing parsing git ssh URL")
	results := utils.ParseGitUrl("git@github.com:bugsnag/bugsnag-cli")
	assert.Equal(t, "https://github.com/bugsnag/bugsnag-cli", results, "They should be the same")

	t.Log("Testing parsing git HTTPS URL")
	results = utils.ParseGitUrl("https://github.com/bugsnag/bugsnag-cli")
	assert.Equal(t, "https://github.com/bugsnag/bugsnag-cli", results, "They should be the same")
}

func TestGetRepoUrl(t *testing.T) {
	t.Log("Test getting repo url from system")
	results, err := utils.GetRepoUrl()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "https://github.com/bugsnag/bugsnag-cli", results, "They should be the same")
}
