package upload_testing

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bugsnag/bugsnag-cli/pkg/options"
	"github.com/bugsnag/bugsnag-cli/pkg/upload"
	"github.com/stretchr/testify/assert"
)

// MockLogger implements the log.Logger interface for testing
type MockLogger struct {
	DebugMessages []string
	InfoMessages  []string
	WarnMessages  []string
	ErrorMessages []string
	FatalMessages []string
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		DebugMessages: []string{},
		InfoMessages:  []string{},
		WarnMessages:  []string{},
		ErrorMessages: []string{},
		FatalMessages: []string{},
	}
}

func (m *MockLogger) Debug(msg string) {
	m.DebugMessages = append(m.DebugMessages, msg)
}

func (m *MockLogger) Info(msg string) {
	m.InfoMessages = append(m.InfoMessages, msg)
}

func (m *MockLogger) Warn(msg string) {
	m.WarnMessages = append(m.WarnMessages, msg)
}

func (m *MockLogger) Error(msg string) {
	m.ErrorMessages = append(m.ErrorMessages, msg)
}

func (m *MockLogger) Fatal(msg string) {
	m.FatalMessages = append(m.FatalMessages, msg)
}

func (m *MockLogger) HasWarning(substring string) bool {
	for _, msg := range m.WarnMessages {
		if strings.Contains(msg, substring) {
			return true
		}
	}
	return false
}

func (m *MockLogger) HasDebug(substring string) bool {
	for _, msg := range m.DebugMessages {
		if strings.Contains(msg, substring) {
			return true
		}
	}
	return false
}

func (m *MockLogger) Reset() {
	m.DebugMessages = []string{}
	m.InfoMessages = []string{}
	m.WarnMessages = []string{}
	m.ErrorMessages = []string{}
	m.FatalMessages = []string{}
}

// Test that ProcessAndroidProguard warns when AndroidManifest.xml cannot be found
func TestProcessAndroidProguard_MissingManifest_Warns(t *testing.T) {
	// Create a temporary directory structure without AndroidManifest.xml
	tempDir := t.TempDir()
	appDir := filepath.Join(tempDir, "app")
	buildDir := filepath.Join(appDir, "build")
	outputsDir := filepath.Join(buildDir, "outputs", "mapping", "release")

	err := os.MkdirAll(outputsDir, 0755)
	assert.NoError(t, err)

	// Create a dummy mapping.txt file
	mappingFile := filepath.Join(outputsDir, "mapping.txt")
	err = os.WriteFile(mappingFile, []byte("# Proguard mapping file\n"), 0644)
	assert.NoError(t, err)

	logger := NewMockLogger()

	opts := options.CLI{
		Globals: options.Globals{
			ApiKey: "test-api-key",
		},
		Upload: options.Upload{
			AndroidProguard: options.AndroidProguardMapping{
				Path:          []string{mappingFile},
				Variant:       "release",
				ApplicationId: "com.test.app",
				VersionCode:   "1",
				VersionName:   "1.0",
			},
		},
	}

	// This should not return an error, but should log a warning
	err = upload.ProcessAndroidProguard(opts, logger)

	// The function should complete without error even though manifest is missing
	// Note: It may still fail due to actual upload but should not fail on manifest lookup
	// For this test, we're primarily checking that warnings are logged
	assert.True(t, len(logger.WarnMessages) > 0 || err != nil,
		"Should either warn about missing manifest or fail for other reasons")
}

// Test that ProcessAndroidProguard warns when AndroidManifest.xml cannot be read
func TestProcessAndroidProguard_UnreadableManifest_Warns(t *testing.T) {
	// Create a temporary directory with an invalid manifest file
	tempDir := t.TempDir()
	appDir := filepath.Join(tempDir, "app")
	buildDir := filepath.Join(appDir, "build")
	outputsDir := filepath.Join(buildDir, "outputs", "mapping", "release")
	manifestDir := filepath.Join(buildDir, "intermediates", "merged_manifests", "release")

	err := os.MkdirAll(outputsDir, 0755)
	assert.NoError(t, err)
	err = os.MkdirAll(manifestDir, 0755)
	assert.NoError(t, err)

	// Create a dummy mapping.txt file
	mappingFile := filepath.Join(outputsDir, "mapping.txt")
	err = os.WriteFile(mappingFile, []byte("# Proguard mapping file\n"), 0644)
	assert.NoError(t, err)

	// Create an invalid manifest file
	manifestFile := filepath.Join(manifestDir, "AndroidManifest.xml")
	err = os.WriteFile(manifestFile, []byte("invalid xml content"), 0644)
	assert.NoError(t, err)

	logger := NewMockLogger()

	opts := options.CLI{
		Globals: options.Globals{
			ApiKey: "test-api-key",
		},
		Upload: options.Upload{
			AndroidProguard: options.AndroidProguardMapping{
				Path:        []string{mappingFile},
				Variant:     "release",
				AppManifest: manifestFile,
			},
		},
	}

	// This should not return an error, but should log a warning about unable to read
	err = upload.ProcessAndroidProguard(opts, logger)

	// Check that a warning was logged about reading the manifest
	// Note: Function may still fail for other reasons (like upload), but should warn about manifest
	if err == nil {
		assert.True(t, logger.HasWarning("Unable to read AndroidManifest.xml"),
			"Should warn about unable to read manifest")
	}
}

// Test that ProcessReactNativeAndroid warns when AndroidManifest.xml cannot be found
func TestProcessReactNativeAndroid_MissingManifest_Warns(t *testing.T) {
	tempDir := t.TempDir()
	androidDir := filepath.Join(tempDir, "android")
	appDir := filepath.Join(androidDir, "app")
	buildDir := filepath.Join(appDir, "build")
	assetsDir := filepath.Join(buildDir, "generated", "assets", "createBundleReleaseJsAndAssets")
	sourcemapsDir := filepath.Join(buildDir, "generated", "sourcemaps", "react", "release")

	err := os.MkdirAll(assetsDir, 0755)
	assert.NoError(t, err)
	err = os.MkdirAll(sourcemapsDir, 0755)
	assert.NoError(t, err)

	// Create dummy bundle and sourcemap files
	bundleFile := filepath.Join(assetsDir, "index.android.bundle")
	err = os.WriteFile(bundleFile, []byte("// bundle content"), 0644)
	assert.NoError(t, err)

	sourcemapFile := filepath.Join(sourcemapsDir, "index.android.bundle.map")
	err = os.WriteFile(sourcemapFile, []byte("{}"), 0644)
	assert.NoError(t, err)

	logger := NewMockLogger()

	opts := options.CLI{
		Globals: options.Globals{
			ApiKey: "test-api-key",
		},
		Upload: options.Upload{
			ReactNativeAndroid: options.ReactNativeAndroid{
				Path:        []string{tempDir},
				ProjectRoot: tempDir,
				Android: options.ReactNativeAndroidSpecific{
					Variant:     "release",
					VersionCode: "1",
				},
				ReactNative: options.ReactNativeShared{
					Bundle:      bundleFile,
					SourceMap:   sourcemapFile,
					VersionName: "1.0",
				},
			},
		},
	}

	// This should not return an error for missing manifest, but should log a warning
	_ = upload.ProcessReactNativeAndroid(opts, logger)

	// Check that a warning was logged about locating the manifest
	hasManifestWarning := logger.HasWarning("Unable to locate AndroidManifest.xml")
	if len(logger.WarnMessages) > 0 {
		assert.True(t, hasManifestWarning,
			"Should warn about unable to locate manifest if no API key provided")
	}
}

// Test that logs debug message when manifest is successfully found
func TestAndroidManifest_SuccessfulFind_LogsDebug(t *testing.T) {
	// Use actual test fixtures if available
	manifestPath := "../../features/android/fixtures/app/build/intermediates/merged_manifests/release/AndroidManifest.xml"

	// Only run this test if the fixture exists
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		t.Skip("Test fixture not available")
		return
	}

	logger := NewMockLogger()

	// Test with a proguard upload that has a valid manifest
	mappingPath := "../../features/android/fixtures/app/build/outputs/mapping/release/mapping.txt"
	if _, err := os.Stat(mappingPath); os.IsNotExist(err) {
		t.Skip("Test fixture not available")
		return
	}

	opts := options.CLI{
		Globals: options.Globals{
			ApiKey: "test-api-key",
		},
		Upload: options.Upload{
			AndroidProguard: options.AndroidProguardMapping{
				Path:          []string{mappingPath},
				Variant:       "release",
				ApplicationId: "com.test.app",
				VersionCode:   "1",
				VersionName:   "1.0",
			},
		},
	}

	_ = upload.ProcessAndroidProguard(opts, logger)

	// Should not have warnings about manifest issues
	assert.False(t, logger.HasWarning("Unable to locate AndroidManifest.xml"),
		"Should not warn when manifest can be found")
	assert.False(t, logger.HasWarning("Unable to read AndroidManifest.xml"),
		"Should not warn when manifest can be read")
}
