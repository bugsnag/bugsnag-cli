package build

import (
	"os"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/build"
	"github.com/stretchr/testify/assert"
)

func TestSetBuilderName(t *testing.T) {
	t.Log("Test setting builders name")
	results, err := build.SetBuilderName("foobar")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "foobar", results, "They should be the same")

	t.Log("Test not setting the builders name")
	results, err = build.SetBuilderName("")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, os.Getenv("USER"), results, "They should be the same")
}

func TestGettingRepoInfo(t *testing.T) {
	t.Log("Test getting repo info map, only setting the commit hash")
	results := build.GetRepoInfo("", "", "git@github.com:bugsnag/bugsnag-cli", "0123456789")

	assert.Equal(t, map[string]string{
		"repository": "git@github.com:bugsnag/bugsnag-cli",
		"revision":   "0123456789",
	}, results, "They should be the same")

	t.Log("Test getting repo info map, passing all three variables")
	results = build.GetRepoInfo("", "github", "https://notgithub.com/bugsnag/bugsnag-cli", "0123456789")
	assert.Equal(t, map[string]string{
		"repository": "https://notgithub.com/bugsnag/bugsnag-cli",
		"revision":   "0123456789",
		"provider":   "github",
	}, results, "They should be the same")
}
