package android_testing

import (
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/stretchr/testify/assert"
)

func TestGetNDKVersion(t *testing.T) {
	t.Log("Testing getting Android NDK major version")
	results, err := android.GetNdkVersion("/opt/homebrew/share/android-commandlinetools/ndk/24.0.8215888")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 24, results, "The versions should match")
}
