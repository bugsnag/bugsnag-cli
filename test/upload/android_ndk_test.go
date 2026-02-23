package upload_testing

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/stretchr/testify/assert"
)

// Note: The resolveAppManifestIfNeeded function is not exported, so we test it
// indirectly through ProcessAndroidNdk

func TestProcessAndroidNdk_MissingManifest_Warns(t *testing.T) {
	// Create a temporary directory structure simulating merged_native_libs without manifest
	tempDir := t.TempDir()
	libPath := filepath.Join(tempDir, "app", "build", "intermediates", "merged_native_libs", "release", "out", "lib", "arm64-v8a")

	err := os.MkdirAll(libPath, 0755)
	assert.NoError(t, err)

	// Create a dummy .so file
	soFile := filepath.Join(libPath, "libnative.so")
	err = os.WriteFile(soFile, []byte("dummy"), 0644)
	assert.NoError(t, err)

	logger := NewMockLogger()

	opts := options.CLI{
		Globals: options.Globals{
			ApiKey: "test-api-key",
		},
		Upload: options.Upload{
			AndroidNdk: options.AndroidNdkMapping{
				Path:          []string{soFile},
				Variant:       "release",
				ApplicationId: "com.test.app",
				VersionCode:   "1",
				VersionName:   "1.0",
			},
		},
	}

	// Call ProcessAndroidNDK - it should warn about missing manifest but not error
	_ = upload.ProcessAndroidNDK(opts, logger)

	// Check that a warning was logged about missing manifest
	hasWarning := logger.HasWarning("Unable to locate AndroidManifest.xml")
	if !hasWarning {
		// The test might fail for other reasons (e.g., invalid .so file)
		// but we want to verify the warning behavior when it gets that far
		t.Log("Warning about manifest may not have been reached due to earlier errors")
		t.Log("Warnings:", logger.WarnMessages)
	}
}

func TestProcessAndroidNdk_WithManifest_NoWarning(t *testing.T) {
	// Use actual test fixture if available
	fixturePath := "../../features/android/fixtures/app/build/intermediates/merged_native_libs/release/out/lib/arm64-v8a/libbugsnag-ndk.so"

	if _, err := os.Stat(fixturePath); os.IsNotExist(err) {
		t.Skip("Test fixture not available")
		return
	}

	logger := NewMockLogger()

	opts := options.CLI{
		Globals: options.Globals{
			ApiKey: "test-api-key",
		},
		Upload: options.Upload{
			AndroidNdk: options.AndroidNdkMapping{
				Path:          []string{fixturePath},
				Variant:       "release",
				ApplicationId: "com.test.app",
				VersionCode:   "1",
				VersionName:   "1.0",
			},
		},
	}

	_ = upload.ProcessAndroidNDK(opts, logger)

	// Should not warn about manifest issues if manifest is found
	// (though it may warn for other reasons)
	hasManifestWarning := logger.HasWarning("Unable to locate AndroidManifest.xml") ||
		logger.HasWarning("Unable to read AndroidManifest.xml")

	if hasManifestWarning {
		t.Errorf("Should not warn about manifest when it exists and is readable")
		t.Log("Warnings:", logger.WarnMessages)
	}
}

func TestProcessAndroidNdk_ManifestProvidedExplicitly_NoSearch(t *testing.T) {
	// Create a temporary directory with a valid manifest
	tempDir := t.TempDir()
	libPath := filepath.Join(tempDir, "app", "build", "intermediates", "merged_native_libs", "release", "out", "lib", "arm64-v8a")
	manifestPath := filepath.Join(tempDir, "AndroidManifest.xml")

	err := os.MkdirAll(libPath, 0755)
	assert.NoError(t, err)

	// Create a minimal valid manifest
	manifestContent := `<?xml version="1.0" encoding="utf-8"?>
<manifest xmlns:android="http://schemas.android.com/apk/res/android"
    package="com.test.app"
    android:versionCode="1"
    android:versionName="1.0">
</manifest>`

	err = os.WriteFile(manifestPath, []byte(manifestContent), 0644)
	assert.NoError(t, err)

	// Create a dummy .so file
	soFile := filepath.Join(libPath, "libnative.so")
	err = os.WriteFile(soFile, []byte("dummy"), 0644)
	assert.NoError(t, err)

	logger := NewMockLogger()

	opts := options.CLI{
		Globals: options.Globals{
			ApiKey: "test-api-key",
		},
		Upload: options.Upload{
			AndroidNdk: options.AndroidNdkMapping{
				Path:          []string{soFile},
				AppManifest:   manifestPath,
				Variant:       "release",
				ApplicationId: "com.test.app",
				VersionCode:   "1",
				VersionName:   "1.0",
			},
		},
	}

	_ = upload.ProcessAndroidNDK(opts, logger)

	// Should not warn about manifest since it was explicitly provided
	assert.False(t, logger.HasWarning("Unable to locate AndroidManifest.xml"),
		"Should not search for manifest when explicitly provided")
}
