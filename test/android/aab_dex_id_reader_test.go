package android_testing

import (
	"fmt"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/android"
	"github.com/stretchr/testify/assert"
)

func TestAabDexFileReader(t *testing.T) {
	t.Log("Testing reading the BuildID from an AAB dex files")
	results, err := android.GetAppSignature("../testdata/android/aab/")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "f3112c3dbdd73ae5dee677e407af196f101e97f5", fmt.Sprintf("%x", results), "The signatures should match")
}
