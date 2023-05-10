package android_testing

import (
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/stretchr/testify/assert"
)

func TestBuildVariantsList(t *testing.T) {
	t.Log("Testing building variants list")
	results, err := android.BuildVariantsList("../testdata/android/variants/")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, []string{"debug", "release"}, results, "The variants should match")
}
