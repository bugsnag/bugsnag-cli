package android_testing

import (
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/stretchr/testify/assert"
)

func TestGetAndroidNDKRoot(t *testing.T) {
	t.Log("Testing getting Android NDK Root")
	results, err := android.GetAndroidNDKRoot("../testdata/android/sdk/ndk/24.0.8215888")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "../testdata/android/sdk/ndk/24.0.8215888", results, "The paths should match")
}
