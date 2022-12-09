package utils_testing_test

import (
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetRepoUrl(t *testing.T) {
	t.Log("Test getting repo url from system")
	results := utils.GetRepoUrl()
	assert.Equal(t, "git@github.com:bugsnag/bugsnag-cli", results, "They should be the same")
}
