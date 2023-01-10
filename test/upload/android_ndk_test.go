package upload_testing

import (
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/stretchr/testify/assert"
)

func TestGetAndroidNDKRoot(t *testing.T) {
	t.Log("Testing getting Android NDK Root")
	results, err := upload.GetAndroidNDKRoot("/opt/homebrew/share/android-commandlinetools/ndk/24.0.8215888")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "/opt/homebrew/share/android-commandlinetools/ndk/24.0.8215888", results, "The paths should match")
}

func TestBuildObjCopyPath(t *testing.T) {
	t.Log("Testing building ObjCopy path")
	results, err := upload.BuildObjCopyPath("/opt/homebrew/share/android-commandlinetools/ndk/24.0.8215888")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "/opt/homebrew/share/android-commandlinetools/ndk/24.0.8215888/toolchains/llvm/prebuilt/darwin-x86_64/bin/llvm-objcopy", results, "The versions should match")
}

func TestGetNDKVersion(t *testing.T) {
	t.Log("Testing getting Android NDK major version")
	results, err := upload.GetNdkVersion("/opt/homebrew/share/android-commandlinetools/ndk/24.0.8215888")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 24, results, "The versions should match")
}

func TestBuildVariantsList(t *testing.T) {
	t.Log("Testing building variants list")
	results, err := upload.BuildVariantsList("../testdata/android/variants/")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, []string{"debug", "release"}, results, "The variants should match")
}
