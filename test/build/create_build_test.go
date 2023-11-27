package build

import (
	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetBuilderName(t *testing.T) {
	t.Log("Test getting the builders name from the system")
	results := utils.GetSystemUser()
	assert.Equal(t, os.Getenv("USER"), results, "They should be the same")
}

func TestGettingRepoInfo(t *testing.T) {
	t.Log("Test getting repo info map, only setting the commit hash")
	results := utils.GetRepoInfo("", "", "git@github.com:bugsnag/bugsnag-cli", "0123456789")

	assert.Equal(t, map[string]string{
		"repository": "git@github.com:bugsnag/bugsnag-cli",
		"revision":   "0123456789",
	}, results, "They should be the same")

	t.Log("Test getting repo info map, passing all three variables")
	results = utils.GetRepoInfo("", "github", "https://notgithub.com/bugsnag/bugsnag-cli", "0123456789")
	assert.Equal(t, map[string]string{
		"repository": "https://notgithub.com/bugsnag/bugsnag-cli",
		"revision":   "0123456789",
		"provider":   "github",
	}, results, "They should be the same")
}
