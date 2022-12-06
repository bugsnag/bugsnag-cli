package utils_testing

import (
	"os"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGetSystemUser(t *testing.T) {
	t.Log("Test getting the system user")
	results, err := utils.GetSystemUser()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, os.Getenv("USER"), results, "They should be the same")
}
