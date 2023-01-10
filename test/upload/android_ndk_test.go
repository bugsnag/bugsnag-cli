package upload_testing

import (
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/stretchr/testify/assert"
)

func TestGetAndroidNDKRoot(t *testing.T) {
	t.Log("Testing getting Android NDK Root")
	results, err := upload.GetAndroidNDKRoot("~/Library/Android/sdk/ndk/23.1.7779620")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "~/Library/Android/sdk/ndk/23.1.7779620", results, "The paths should match")
}


func TestBuildObjCopyPath(t *testing.T) {
	t.Log("Testing building ObjCopy path")
	results, err := upload.BuildObjCopyPath("~/Library/Android/sdk/ndk/23.1.7779620")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 23, results, "The versions should match")
}

func TestGetNDKVersion(t *testing.T) {
	t.Log("Testing getting Android NDK major version")
	results, err := upload.GetNdkVersion("~/Library/Android/sdk/ndk/23.1.7779620")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 23, results, "The versions should match")
}

func TestBuildVariantsList(t *testing.T) {
	t.Log("Testing building variants list")
	results, err := upload.BuildVariantsList("../testdata/android/variants/")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, []string{"debug", "release"}, results, "The variants should match")
}
