package utils_testing

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
)

func requireXcodebuild(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("xcodebuild"); err != nil {
		t.Skip("xcodebuild is required for this test")
	}
}

// Tests expected scenarios where project root is set based on the value of <path> or --project-root
func TestDefaultProjectRoot(t *testing.T) {
	tt := map[string]struct {
		pathValue           string
		projectRootValue    string
		expectedProjectRoot string
	}{
		"<path> contains normal directory and is used as project root": {
			pathValue:           "../testdata/ios/SingleSchemeExample",
			projectRootValue:    "",
			expectedProjectRoot: "../testdata/ios/SingleSchemeExample",
		},
		"<path> contains xcodeproj directory and one directory up is used as project root": {
			pathValue:           "../testdata/ios/SingleSchemeExample/SingleSchemeExample.xcodeproj",
			projectRootValue:    "",
			expectedProjectRoot: "../testdata/ios/SingleSchemeExample",
		},
		"--project-root is set and is used as project root": {
			pathValue:           "../testdata/ios/SingleSchemeExample/SingleSchemeExample.xcodeproj",
			projectRootValue:    "/Path/To/ProjectRoot",
			expectedProjectRoot: "/Path/To/ProjectRoot",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			actualProjectRoot := ios.GetDefaultProjectRoot(tc.pathValue, tc.projectRootValue)
			assert.Equal(t, tc.expectedProjectRoot, actualProjectRoot)
		})
	}
}

// Tests expected common use cases when determining the default scheme
func TestGetDefaultScheme(t *testing.T) {
	requireXcodebuild(t)

	tt := map[string]struct {
		pathValue      string
		expectedScheme string
	}{
		"<path> contains a normal directory and is used to fetch the scheme": {
			pathValue:      "../testdata/ios/SingleSchemeExample/",
			expectedScheme: "SingleSchemeExample",
		},
		"<path> contains a .xcodeproj directory and is used to fetch the scheme": {
			pathValue:      "../testdata/ios/SingleSchemeExample/SingleSchemeExample.xcodeproj",
			expectedScheme: "SingleSchemeExample",
		},
		"<path> contains a .xcworkspace directory and is used to fetch the scheme": {
			pathValue:      "../testdata/ios/WorkspaceExample.xcworkspace",
			expectedScheme: "WorkspaceScheme",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			actualScheme, err := ios.GetDefaultScheme(tc.pathValue)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedScheme, actualScheme)
		})
	}
}

// Tests expected common error scenarios when determining the default scheme
func TestGetDefaultSchemeErrorScenarios(t *testing.T) {
	requireXcodebuild(t)

	tt := map[string]struct {
		pathValue            string
		expectedExceptionMsg string
	}{
		"multiple schemes found results in exception": {
			pathValue:            "../testdata/ios/MultipleSchemeExample/MultipleSchemeExample.xcodeproj",
			expectedExceptionMsg: "multiple schemes found",
		},
		"no schemes found results in exception": {
			pathValue:            "../testdata/ios/parent_root",
			expectedExceptionMsg: "no schemes found",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			_, err := ios.GetDefaultScheme(tc.pathValue)

			require.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedExceptionMsg)
		})
	}
}

// Tests expected use cases when fetching build settings
func TestGetXcodeBuildSettings(t *testing.T) {
	requireXcodebuild(t)

	tt := map[string]struct {
		pathValue      string
		scheme         string
		expectedResult *ios.XcodeBuildSettings
	}{
		"successfully retrieve build settings for xcodeproj and scheme": {
			pathValue: "../testdata/ios/SingleSchemeExample/SingleSchemeExample.xcodeproj",
			scheme:    "SingleSchemeExample",
			expectedResult: &ios.XcodeBuildSettings{
				ConfigurationBuildDir: "Build/Products/Debug-iphoneos",
				InfoPlistPath:         "SingleSchemeExample.app/Info.plist",
				BuiltProductsDir:      "Build/Products/Debug-iphoneos",
				DsymName:              "SingleSchemeExample.app.dSYM",
			},
		},
		"successfully retrieve build settings for xcworkspace and scheme": {
			pathValue: "../testdata/ios/BuildSettingsExample.xcworkspace",
			scheme:    "BuildSettingsScheme",
			expectedResult: &ios.XcodeBuildSettings{
				ConfigurationBuildDir: "Build/Products/Debug-iphoneos",
				InfoPlistPath:         "SingleSchemeExample.app/Info.plist",
				BuiltProductsDir:      "Build/Products/Debug-iphoneos",
				DsymName:              "SingleSchemeExample.app.dSYM",
			},
		},
		"successfully retrieve build settings for path to project root and scheme": {
			pathValue: "../testdata/ios/SingleSchemeExample",
			scheme:    "SingleSchemeExample",
			expectedResult: &ios.XcodeBuildSettings{
				ConfigurationBuildDir: "Build/Products/Debug-iphoneos",
				InfoPlistPath:         "SingleSchemeExample.app/Info.plist",
				BuiltProductsDir:      "Build/Products/Debug-iphoneos",
				DsymName:              "SingleSchemeExample.app.dSYM",
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			actualResult, err := ios.GetXcodeBuildSettings(tc.pathValue, tc.scheme, "")
			require.NoError(t, err)
			assert.NotNil(t, actualResult)

			assert.Contains(t, actualResult.ConfigurationBuildDir, tc.expectedResult.ConfigurationBuildDir)
			assert.Contains(t, actualResult.InfoPlistPath, tc.expectedResult.InfoPlistPath)
			assert.Contains(t, actualResult.BuiltProductsDir, tc.expectedResult.BuiltProductsDir)
			assert.Equal(t, tc.expectedResult.DsymName, actualResult.DsymName)
		})
	}
}
