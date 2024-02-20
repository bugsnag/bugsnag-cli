package utils_testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
)

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
	tt := map[string]struct {
		pathValue            string
		expectedExceptionMsg string
	}{
		"multiple schemes found results in exception": {
			pathValue:            "../testdata/ios/MultipleSchemeExample/MultipleSchemeExample.xcodeproj",
			expectedExceptionMsg: "Multiple schemes found",
		},
		"no schemes found results in exception": {
			pathValue:            "../testdata/ios/parent_root",
			expectedExceptionMsg: "No schemes found",
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			_, err := ios.GetDefaultScheme(tc.pathValue)

			assert.Contains(t, err.Error(), tc.expectedExceptionMsg)
		})
	}
}

// Tests expected use cases when fetching build settings
func TestGetXcodeBuildSettings(t *testing.T) {
	tt := map[string]struct {
		pathValue      string
		scheme         string
		expectedResult *ios.XcodeBuildSettings
	}{
		"successfully retrieve build settings for xcodeproj and scheme": {
			pathValue: "../../features/base-fixtures/rn0_72/ios/rn0_72.xcodeproj",
			scheme:    "rn0_72",
			expectedResult: &ios.XcodeBuildSettings{
				ConfigurationBuildDir: "Build/Products/Release-iphoneos",
				InfoPlistPath:         "Info.plist",
				BuiltProductsDir:      "Build/Products/Release-iphoneos",
				DsymName:              "rn0_72.app.dSYM",
			},
		},
		"successfully retrieve build settings for xcworkspace and scheme": {
			pathValue: "../../features/base-fixtures/rn0_69/ios/rn0_69.xcworkspace",
			scheme:    "rn0_69",
			expectedResult: &ios.XcodeBuildSettings{
				ConfigurationBuildDir: "Build/Products/Release-iphoneos",
				InfoPlistPath:         "Info.plist",
				BuiltProductsDir:      "Build/Products/Release-iphoneos",
				DsymName:              "rn0_69.app.dSYM",
			},
		},
		"successfully retrieve build settings for path to project root and scheme": {
			pathValue: "../../features/base-fixtures/rn0_70/ios/",
			scheme:    "rn0_70",
			expectedResult: &ios.XcodeBuildSettings{
				ConfigurationBuildDir: "Build/Products/Release-iphoneos",
				InfoPlistPath:         "Info.plist",
				BuiltProductsDir:      "Build/Products/Release-iphoneos",
				DsymName:              "rn0_70.app.dSYM",
			},
		},
		"successfully retrieve build settings for projectRoot and scheme": {
			pathValue: "../../features/base-fixtures/rn0_69/ios/",
			scheme:    "rn0_69",
			expectedResult: &ios.XcodeBuildSettings{
				ConfigurationBuildDir: "Build/Products/Release-iphoneos",
				InfoPlistPath:         "Info.plist",
				BuiltProductsDir:      "Build/Products/Release-iphoneos",
				DsymName:              "rn0_69.app.dSYM",
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			actualResult, err := ios.GetXcodeBuildSettings(tc.pathValue, tc.scheme)
			require.NoError(t, err)
			assert.NotNil(t, actualResult)

			assert.Contains(t, actualResult.ConfigurationBuildDir, tc.expectedResult.ConfigurationBuildDir)
			assert.Contains(t, actualResult.InfoPlistPath, tc.expectedResult.InfoPlistPath)
			assert.Contains(t, actualResult.BuiltProductsDir, tc.expectedResult.BuiltProductsDir)
		})
	}
}
