package utils_testing

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/bugsnag/bugsnag-cli/pkg/ios"
)

// Tests expected common use case behaviour for processing <path> value
func TestProcessPathValue(t *testing.T) {
	currentDir, _ := os.Getwd()

	tt := map[string]struct {
		pathValue      string
		dsymValue      string
		projectRoot    string
		expectedResult *ios.DsymUploadInfo
	}{
		"if <path> is set, is a normal directory and --project-root is not set, value of <path> is returned as-is": {
			pathValue:   "../testdata/ios/parent_root",
			dsymValue:   "",
			projectRoot: "",
			expectedResult: &ios.DsymUploadInfo{
				ProjectRoot: "../testdata/ios/parent_root",
			},
		},
		"if <path> is set, is a .xcodeproj directory and --project-root is not set, one directory up from <path> is returned": {
			pathValue:   "../testdata/ios/parent_root/MyTestApp.xcodeproj",
			dsymValue:   "",
			projectRoot: "",
			expectedResult: &ios.DsymUploadInfo{
				ProjectRoot: "../testdata/ios/parent_root",
			},
		},
		"if <path> is set, is a .xcworkspace directory and --project-root is not set, one directory up from <path> is returned": {
			pathValue:   "../testdata/ios/parent_root/MyTestApp.xcworkspace",
			dsymValue:   "",
			projectRoot: "",
			expectedResult: &ios.DsymUploadInfo{
				ProjectRoot: "../testdata/ios/parent_root",
			},
		},
		"if <path> and --project-root are both unset, current working directory is returned as project root": {
			pathValue:   "",
			dsymValue:   "",
			projectRoot: "",
			expectedResult: &ios.DsymUploadInfo{
				ProjectRoot: currentDir,
			},
		},
		"if <path> is a file then set the dsym path to it's value": {
			pathValue:   "../testdata/ios/MyTestApp",
			dsymValue:   "",
			projectRoot: "",
			expectedResult: &ios.DsymUploadInfo{
				DsymPaths: &[]string{"../testdata/ios/MyTestApp"},
			},
		},
		"if <path> is a .zip then set the dsym path to it's value": {
			pathValue:   "../testdata/ios/MyTestApp.zip",
			dsymValue:   "",
			projectRoot: "",
			expectedResult: &ios.DsymUploadInfo{
				DsymPaths: &[]string{"../testdata/ios/MyTestApp.zip"},
			},
		},
		"if <path> is a file then set the dsym path to it's value and if --project-root is set, use it's value for project root": {
			pathValue:   "../testdata/ios/MyTestApp",
			dsymValue:   "",
			projectRoot: "../testdata/ios/parent_root",
			expectedResult: &ios.DsymUploadInfo{
				ProjectRoot: "../testdata/ios/parent_root",
				DsymPaths:   &[]string{"../testdata/ios/MyTestApp"},
			},
		},
	}

	for name, tc := range tt {
		t.Run(name, func(t *testing.T) {
			actualResult, err := ios.ProcessPathValues(tc.pathValue, tc.dsymValue, tc.projectRoot)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedResult, actualResult)
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
